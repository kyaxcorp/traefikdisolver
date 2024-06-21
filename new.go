package traefik_client_real_ip

import (
	"context"
	"errors"
	"net"
	"net/http"
)

// New created a new plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.Provider == "" {
		return nil, errors.New("no provider has been defined")
	}

	


	realIPUpdater := &RealIPUpdater{
		next:     next,
		name:     name,
		provider: config.Provider,
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
		switch config.Provider{
			case 
		}

		for _, v := range ips.CFIPs() {
			_, trustip, err := net.ParseCIDR(v)
			if err != nil {
				return nil, err
			}

			realIPUpdater.TrustIP = append(realIPUpdater.TrustIP, trustip)
		}
	}

	return realIPUpdater, nil
}
