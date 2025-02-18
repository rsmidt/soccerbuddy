package grpc

import (
	"connectrpc.com/connect"
	"net"
	"net/http"
	"testing"
)

type testRequest struct {
	peerData   connect.Peer
	headerData http.Header
}

func (r *testRequest) Peer() connect.Peer  { return r.peerData }
func (r *testRequest) Header() http.Header { return r.headerData }

// newTestRequest creates a testRequest (which implements connect.AnyRequest)
// with the given X-Forwarded-For header values and peer address (i.e. remote addr).
func newTestRequest(headerValues []string, remoteAddr string) AnyRequest {
	h := http.Header{}
	for _, v := range headerValues {
		h.Add("X-Forwarded-For", v)
	}
	return &testRequest{
		peerData:   connect.Peer{Addr: remoteAddr},
		headerData: h,
	}
}

// --- Helper for tests ---

// mustParseIPNet is a helper for tests to parse CIDR and panic on error.
func mustParseIPNet(cidr string) *net.IPNet {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		panic("invalid CIDR: " + cidr)
	}
	return ipnet
}

// --- Tests ---

func TestGetClientIP_Untrusted(t *testing.T) {
	tests := []struct {
		name         string
		headerValues []string
		remoteAddr   string // used for fallback (peer.Addr)
		expected     string // expected IP as string; empty means nil
	}{
		{
			name:         "Single public IP",
			headerValues: []string{"203.0.113.195"},
			remoteAddr:   "203.0.113.195:1234",
			expected:     "203.0.113.195",
		},
		{
			name:         "Multiple IPs, first private then public",
			headerValues: []string{"192.168.1.1, 203.0.113.195"},
			remoteAddr:   "203.0.113.195:1234",
			expected:     "203.0.113.195",
		},
		{
			name:         "Multiple headers",
			headerValues: []string{"192.168.1.1", "203.0.113.195"},
			remoteAddr:   "203.0.113.195:1234",
			expected:     "203.0.113.195",
		},
		{
			name:         "No public IP in header, fallback to peer.Addr",
			headerValues: []string{"192.168.1.1, 10.0.0.5"},
			remoteAddr:   "198.51.100.23:5678",
			expected:     "198.51.100.23",
		},
		{
			name:         "Invalid IP in header is skipped",
			headerValues: []string{"not_an_ip, 203.0.113.195"},
			remoteAddr:   "203.0.113.195:1234",
			expected:     "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestRequest(tt.headerValues, tt.remoteAddr)
			// nil config => untrusted mode
			ip := GetClientIP(req, nil)
			if tt.expected == "" {
				if ip != nil {
					t.Errorf("expected nil, got %v", ip)
				}
			} else {
				if ip == nil || ip.String() != tt.expected {
					t.Errorf("expected %s, got %v", tt.expected, ip)
				}
			}
		})
	}
}

func TestGetClientIP_TrustedProxyCount(t *testing.T) {
	// In trusted proxy count mode, the client IP is determined by
	// skipping (TrustedProxyCount - 1) addresses from the right.
	tests := []struct {
		name              string
		headerValues      []string
		trustedProxyCount int
		expected          string
	}{
		{
			name:              "Single proxy: client IP is rightmost",
			headerValues:      []string{"198.51.100.17, 203.0.113.195"},
			trustedProxyCount: 1,
			expected:          "203.0.113.195",
		},
		{
			name:              "Multiple proxies: returns correct IP",
			headerValues:      []string{"198.51.100.17, 192.0.2.1, 203.0.113.195"},
			trustedProxyCount: 2, // skip last trusted proxy, return the one before it
			expected:          "192.0.2.1",
		},
		{
			name:              "Not enough IPs for proxy count",
			headerValues:      []string{"203.0.113.195"},
			trustedProxyCount: 2,
			expected:          "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestRequest(tt.headerValues, "")
			cfg := &Config{
				TrustedProxyCount: tt.trustedProxyCount,
			}
			ip := GetClientIP(req, cfg)
			if tt.expected == "" {
				if ip != nil {
					t.Errorf("expected nil, got %v", ip)
				}
			} else {
				if ip == nil || ip.String() != tt.expected {
					t.Errorf("expected %s, got %v", tt.expected, ip)
				}
			}
		})
	}
}

func TestGetClientIP_TrustedProxyList(t *testing.T) {
	// Define a trusted proxy list.
	trustedList := []*net.IPNet{
		mustParseIPNet("192.0.2.0/24"),
		mustParseIPNet("198.51.100.0/24"),
	}

	tests := []struct {
		name         string
		headerValues []string
		trustedList  []*net.IPNet
		expected     string
	}{
		{
			name:         "Trusted proxies at right, returns first non-trusted IP",
			headerValues: []string{"203.0.113.195, 192.0.2.17, 198.51.100.23"},
			trustedList:  trustedList,
			expected:     "203.0.113.195",
		},
		{
			name:         "All IPs trusted, returns nil",
			headerValues: []string{"192.0.2.17, 198.51.100.23"},
			trustedList:  trustedList,
			expected:     "",
		},
		{
			name:         "Multiple headers with mix",
			headerValues: []string{"192.0.2.17", "203.0.113.195, 198.51.100.23"},
			trustedList:  trustedList,
			expected:     "203.0.113.195",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := newTestRequest(tt.headerValues, "")
			cfg := &Config{
				TrustedProxies: tt.trustedList,
			}
			ip := GetClientIP(req, cfg)
			if tt.expected == "" {
				if ip != nil {
					t.Errorf("expected nil, got %v", ip)
				}
			} else {
				if ip == nil || ip.String() != tt.expected {
					t.Errorf("expected %s, got %v", tt.expected, ip)
				}
			}
		})
	}
}
