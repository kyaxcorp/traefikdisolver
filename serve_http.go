package traefikdisolver

import (
	"encoding/json"
	"net"
	"net/http"

	"github.com/kyaxcorp/traefikdisolver/providers"
	"github.com/kyaxcorp/traefikdisolver/providers/cloudflare"
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
		switch r.provider {
		case providers.Cloudflare:
			if req.Header.Get(cloudflare.CfVisitor) != "" {
				var cfVisitorValue CFVisitorHeader
				if err := json.Unmarshal([]byte(req.Header.Get(cloudflare.CfVisitor)), &cfVisitorValue); err != nil {
					req.Header.Set(cloudflare.XCfTrusted, "danger")
					req.Header.Del(cloudflare.CfVisitor)
					req.Header.Del(cloudflare.ClientIPHeaderName)
					r.next.ServeHTTP(rw, req)
					return
				}
				req.Header.Set(xForwardProto, cfVisitorValue.Scheme)
			}
		case providers.Cloudfront:
		}

		switch r.provider {
		case providers.Cloudflare:
			req.Header.Set(cloudflare.XCfTrusted, "yes")
		case providers.Cloudfront:
		}

		clientIP, _, _ := net.SplitHostPort(req.Header.Get(r.clientIPHeaderName))

		req.Header.Set(xForwardFor, clientIP)
		req.Header.Set(xRealIP, clientIP)
	} else {
		switch r.provider {
		case providers.Cloudflare:
			req.Header.Set(cloudflare.XCfTrusted, "no")
			req.Header.Del(cloudflare.CfVisitor)
			req.Header.Del(cloudflare.ClientIPHeaderName)
		case providers.Cloudfront:
		}
		req.Header.Set(xRealIP, trustResult.directIP)
	}
	r.next.ServeHTTP(rw, req)
}
