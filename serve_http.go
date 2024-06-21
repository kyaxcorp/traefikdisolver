package traefikdisolver

import (
	"encoding/json"
	"net/http"
)

func (r *Disolver) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	trustResult := r.trust(req.RemoteAddr)
	if trustResult.isFatal {
		http.Error(rw, "Unknown source", http.StatusInternalServerError)
		return
	}
	if trustResult.isError {
		http.Error(rw, "Unknown source", http.StatusBadRequest)
		return
	}
	if trustResult.directIP == "" {
		http.Error(rw, "Unknown source", http.StatusUnprocessableEntity)
		return
	}
	if trustResult.trusted {
		if req.Header.Get(cfVisitor) != "" {
			var cfVisitorValue CFVisitorHeader
			if err := json.Unmarshal([]byte(req.Header.Get(cfVisitor)), &cfVisitorValue); err != nil {
				req.Header.Set(xCfTrusted, "danger")
				req.Header.Del(cfVisitor)
				req.Header.Del(cfConnectingIP)
				r.next.ServeHTTP(rw, req)
				return
			}
			req.Header.Set(xForwardProto, cfVisitorValue.Scheme)
		}
		req.Header.Set(xCfTrusted, "yes")
		req.Header.Set(xForwardFor, req.Header.Get(cfConnectingIP))
		req.Header.Set(xRealIP, req.Header.Get(cfConnectingIP))
	} else {
		req.Header.Set(xCfTrusted, "no")
		req.Header.Set(xRealIP, trustResult.directIP)
		req.Header.Del(cfVisitor)
		req.Header.Del(cfConnectingIP)
	}
	r.next.ServeHTTP(rw, req)
}
