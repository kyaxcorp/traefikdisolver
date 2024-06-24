package traefikdisolver

import (
	"net"
	"net/http"

	"github.com/kyaxcorp/traefikdisolver/providers"
	"github.com/kyaxcorp/traefikdisolver/providers/cloudflare"
	"github.com/kyaxcorp/traefikdisolver/providers/cloudfront"
)

// Disolver is a plugin that overwrite true IP.
type Disolver struct {
	next               http.Handler
	name               string
	provider           providers.Provider
	TrustIP            map[providers.Provider][]*net.IPNet
	clientIPHeaderName string
}

// CFVisitorHeader definition for the header value.
type CFVisitorHeader struct {
	Scheme string `json:"scheme"`
}

func (r *Disolver) trust(s string, req *http.Request) *TrustResult {
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

	var ips []*net.IPNet

	switch r.provider {
	case providers.Cloudflare:
		ips = r.TrustIP[providers.Cloudflare]
	case providers.Cloudfront:
		ips = r.TrustIP[providers.Cloudfront]
	case providers.Auto:
		provider := detectProvider(req)
		switch provider {
		case providers.Cloudflare:
			ips = r.TrustIP[providers.Cloudflare]
		case providers.Cloudfront:
			ips = r.TrustIP[providers.Cloudfront]
		}
	}

	for _, network := range ips {
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

func detectProvider(req *http.Request) providers.Provider {
	if req.Header.Get(cloudflare.ClientIPHeaderName) != "" {
		return providers.Cloudflare
	} else if req.Header.Get(cloudfront.ClientIPHeaderName) != "" {
		return providers.Cloudfront
	}
	return providers.Unknown
}
