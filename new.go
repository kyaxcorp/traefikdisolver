package traefikdisolver

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/kyaxcorp/traefikdisolver/providers"
	"github.com/kyaxcorp/traefikdisolver/providers/cloudflare"
	"github.com/kyaxcorp/traefikdisolver/providers/cloudfront"
)

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Provider == "" {
		return nil, errors.New("no provider has been defined")
	}

	provider := providers.Provider(config.Provider)
	if provider.Validate() != nil {
		return nil, errors.New("failed to validate provider")
	}

	realIPUpdater := &Disolver{
		next:     next,
		name:     name,
		provider: provider,
		TrustIP:  make(map[providers.Provider][]*net.IPNet),
	}

	switch provider {
	case providers.Cloudflare:
		realIPUpdater.clientIPHeaderName = cloudflare.ClientIPHeaderName
	case providers.Cloudfront:
		realIPUpdater.clientIPHeaderName = cloudfront.ClientIPHeaderName
	}

	if config.DisableDefaultCFIPs {
		for prov, _ := range providers.ListExisting {
			var trustIPs []string
			var exist bool
			if trustIPs, exist = config.TrustIP[prov.String()]; !exist {
				// Let's get the default ones!
				switch prov {
				case providers.Cloudflare:
					trustIPs = cloudflare.TrustedIPS()
				case providers.Cloudfront:
					trustIPs = cloudfront.TrustedIPS()
				}
			}
			for _, v := range trustIPs {
				_, trustip, err := net.ParseCIDR(v)
				if err != nil {
					return nil, err
				}
				realIPUpdater.TrustIP[prov] = append(realIPUpdater.TrustIP[prov], trustip)
			}
		}
	} else {
		var ips []string
		switch provider {
		case providers.Cloudflare:
			ips = cloudflare.TrustedIPS()
		case providers.Cloudfront:
			ips = cloudfront.TrustedIPS()
		case providers.Auto:
			// ips = auto.TrustedIPS()
		}

		switch provider {
		case providers.Cloudflare, providers.Cloudfront:
			for _, v := range ips {
				_, trustip, err := net.ParseCIDR(v)
				if err != nil {
					return nil, err
				}

				realIPUpdater.TrustIP[provider] = append(realIPUpdater.TrustIP[provider], trustip)
			}
		case providers.Auto:
			ips = cloudflare.TrustedIPS()
			for _, v := range ips {
				_, trustip, err := net.ParseCIDR(v)
				if err != nil {
					return nil, err
				}

				realIPUpdater.TrustIP[providers.Cloudflare] = append(realIPUpdater.TrustIP[providers.Cloudflare], trustip)
			}
			ips = cloudfront.TrustedIPS()
			for _, v := range ips {
				_, trustip, err := net.ParseCIDR(v)
				if err != nil {
					return nil, err
				}

				realIPUpdater.TrustIP[providers.Cloudfront] = append(realIPUpdater.TrustIP[providers.Cloudfront], trustip)
			}
		}
	}

	return realIPUpdater, nil
}
