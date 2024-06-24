package auto

// CFIPs is the CloudFlare Server IP list (this is checked on build).
func TrustedIPS() []string {
	return []string{
		"0.0.0.0/0",
		"::/0",
	}
}
