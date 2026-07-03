// HTTP resolver
// Detect IP address by HTTP / HTTPS response
package httpresolver

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/kitsuyui/myip/resolvers/base"
)

type HTTPDetector struct {
	URL string `json:"url"`
}

func scoreOfTLS(t *tls.ConnectionState) float64 {
	if t == nil { // HTTP
		return 0.1
	}
	if t.Version >= tls.VersionTLS13 {
		return 1.0
	}
	switch t.Version {
	case tls.VersionTLS12:
		return 0.8
	case tls.VersionTLS11:
		return 0.6
	case tls.VersionTLS10:
		return 0.4
	default:
		return 0.2
	}
}

func (p HTTPDetector) RetrieveIP() (*base.ScoredIP, error) {
	return p.RetrieveIPWithContext(context.Background())
}

// RetrieveIPWithContext retrieves an IP address using an HTTP request bound to ctx.
func (p HTTPDetector) RetrieveIPWithContext(ctx context.Context) (*base.ScoredIP, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
	if err != nil {
		return nil, err
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxResponseHeaderBytes = 4096
	client := http.Client{Transport: transport}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 64))
	if err != nil {
		return nil, err
	}
	ipStr := strings.TrimSpace(string(body))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, &base.NotRetrievedError{}
	}
	return &base.ScoredIP{IP: ip, Score: scoreOfTLS(resp.TLS)}, nil
}

func (p HTTPDetector) String() string {
	return p.URL
}

func (p HTTPDetector) TypeName() string {
	return "http"
}
