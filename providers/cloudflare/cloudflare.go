// Package ips contains a list of current cloud flare IP ranges
package cloudflare

// CFIPs is the CloudFlare Server IP list (this is checked on build).
func TrustedIPS() []string {
	return []string{
		"0.0.0.0/0",
		"::/0",
	}
}

const ClientIPHeaderName = "CF-Connecting-IP"
const CfVisitor = "CF-Visitor"
const XCfTrusted = "X-Is-Trusted"
