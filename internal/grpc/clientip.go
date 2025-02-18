package grpc

import (
	"connectrpc.com/connect"
	"net"
	"net/http"
	"strings"
)

// AnyRequest re-declares a subset of the `connect.AnyRequest` interface because it cannot be implemented
// outside the connect package.
type AnyRequest interface {
	Peer() connect.Peer
	Header() http.Header
}

// Config holds configuration for determining the client IP.
type Config struct {
	// If TrustedProxyCount > 0, then use the “trusted proxy count” method.
	TrustedProxyCount int

	// Alternatively, if TrustedProxies is non‑empty, then use that list.
	TrustedProxies []*net.IPNet
}

// GetClientIP returns the client IP address from the given request.
// It processes all X‑Forwarded‑For headers (splitting comma‑separated values)
// and then selects the IP based on the provided configuration.
// If no X‑Forwarded‑For header is present (or no valid IP is found),
// it falls back to parsing req.peer.Addr.
func GetClientIP(req AnyRequest, cfg *Config) net.IP {
	ips := parseXFF(req)
	if len(ips) == 0 {
		// Fallback: try req.peer.Addr.
		return ipFromPeerAddr(req.Peer().Addr)
	}

	// If configuration is provided, use one of the trusted methods.
	if cfg != nil {
		// Use TrustedProxies list (if provided) first.
		if len(cfg.TrustedProxies) > 0 {
			// Scan from the rightmost IP.
			for i := len(ips) - 1; i >= 0; i-- {
				ip := ips[i]
				if !isTrusted(ip, cfg.TrustedProxies) {
					return ip
				}
			}
			// All addresses came from a trusted proxy.
			return nil
		}

		// Otherwise, if TrustedProxyCount is set, pick the IP by count.
		if cfg.TrustedProxyCount > 0 {
			if len(ips) >= cfg.TrustedProxyCount {
				// When there is one proxy, the client IP is the rightmost IP.
				// In general, skip TrustedProxyCount - 1 addresses from the right.
				index := len(ips) - cfg.TrustedProxyCount
				return ips[index]
			}
			// Not enough IPs.
			return nil
		}
	}

	// Untrusted mode: choose the first IP (from the left)
	// that is a valid IP and not private/internal.
	for _, ip := range ips {
		if !isPrivateIP(ip) {
			return ip
		}
	}

	// If none qualifies, return the peer.Addr.
	return ipFromPeerAddr(req.Peer().Addr)
}

// parseXFF extracts and parses all X‑Forwarded‑For header values from req.header.
func parseXFF(req AnyRequest) []net.IP {
	var ips []net.IP
	// Note: req.header is case‑insensitive.
	headers := req.Header().Values("X-Forwarded-For")
	for _, header := range headers {
		// Some proxies join multiple IPs with commas.
		parts := strings.Split(header, ",")
		for _, part := range parts {
			ipStr := strings.TrimSpace(part)
			if ip := net.ParseIP(ipStr); ip != nil {
				ips = append(ips, ip)
			}
		}
	}
	return ips
}

// ipFromPeerAddr extracts the IP from a peer address (typically "IP:port").
func ipFromPeerAddr(addr string) net.IP {
	if addr == "" {
		return nil
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		// If splitting fails, try to parse the whole string.
		return net.ParseIP(addr)
	}
	return net.ParseIP(host)
}

// isPrivateIP returns true if ip is in one of the private or loopback networks.
func isPrivateIP(ip net.IP) bool {
	privateBlocks := []*net.IPNet{
		mustParseCIDR("10.0.0.0/8"),
		mustParseCIDR("172.16.0.0/12"),
		mustParseCIDR("192.168.0.0/16"),
		mustParseCIDR("127.0.0.0/8"),
		mustParseCIDR("fc00::/7"),
		mustParseCIDR("fe80::/10"),
		mustParseCIDR("::1/128"),
	}
	for _, block := range privateBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// mustParseCIDR is a helper that panics if the CIDR cannot be parsed.
func mustParseCIDR(cidr string) *net.IPNet {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		panic("invalid CIDR: " + cidr)
	}
	return network
}

// isTrusted returns true if the given ip is contained in any of the trusted proxy ranges.
func isTrusted(ip net.IP, trusted []*net.IPNet) bool {
	for _, block := range trusted {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}
