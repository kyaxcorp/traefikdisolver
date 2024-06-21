package traefikdisolver

// Config the plugin configuration.
type Config struct {
	Provider            string   `json:"provider,omitempty"`
	TrustIP             []string `json:"trustip,omitempty"`
	DisableDefaultCFIPs bool     `json:"disableDefault,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Provider:            "", // TODO: if no provider has been set...
		TrustIP:             []string{},
		DisableDefaultCFIPs: false,
	}
}
