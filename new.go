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
	}

	if config.TrustIP != nil {
		for _, v := range config.TrustIP {
			_, trustip, err := net.ParseCIDR(v)
			if err != nil {
				return nil, err
			}

			realIPUpdater.TrustIP = append(realIPUpdater.TrustIP, trustip)
		}
	}

	if !config.DisableDefaultCFIPs {
		var ips []string
		switch provider {
		case providers.Cloudflare:
			ips = cloudflare.TrustedIPS()
			realIPUpdater.clientIPHeaderName = cloudflare.ClientIPHeaderName
		case providers.Cloudfront:
			ips = cloudfront.TrustedIPS()
			realIPUpdater.clientIPHeaderName = cloudfront.ClientIPHeaderName
		}

		for _, v := range ips {
			_, trustip, err := net.ParseCIDR(v)
			if err != nil {
				return nil, err
			}

			realIPUpdater.TrustIP = append(realIPUpdater.TrustIP, trustip)
		}
	}

	return realIPUpdater, nil
}
