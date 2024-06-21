package traefikdisolver

import (
	"net"
	"net/http"

	"github.com/kyaxcorp/traefikdisolver/providers"
)

// Disolver is a plugin that overwrite true IP.
type Disolver struct {
	next               http.Handler
	name               string
	provider           providers.Provider
	TrustIP            []*net.IPNet
	clientIPHeaderName string
}

// CFVisitorHeader definition for the header value.
type CFVisitorHeader struct {
	Scheme string `json:"scheme"`
}

func (r *Disolver) trust(s string) *TrustResult {
	temp, _, err := net.SplitHostPort(s)
	if err != nil {
		return &TrustResult{
			isFatal:  true,
			isError:  true,
			trusted:  false,
			directIP: "",
		}
	}
	ip := net.ParseIP(temp)
	if ip == nil {
		return &TrustResult{
			isFatal:  false,
			isError:  true,
			trusted:  false,
			directIP: "",
		}
	}
	for _, network := range r.TrustIP {
		if network.Contains(ip) {
			return &TrustResult{
				isFatal:  false,
				isError:  false,
				trusted:  true,
				directIP: ip.String(),
			}
		}
	}
	return &TrustResult{
		isFatal:  false,
		isError:  false,
		trusted:  false,
		directIP: ip.String(),
	}
}

// TrustResult for Trust IP test result.
type TrustResult struct {
	isFatal  bool
	isError  bool
	trusted  bool
	directIP string
}
