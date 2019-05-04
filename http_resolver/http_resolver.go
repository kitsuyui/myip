package http_resolver

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/kitsuyui/myip/base"
)

type HTTPDetector struct {
	URL string `json:"url"`
}

func scoreOfTLS(t *tls.ConnectionState) float64 {
	if t == nil { // HTTP
		return 0.1
	}
	switch t.Version {
	case tls.VersionTLS13:
		return 1.0
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
	client := http.Client{}
	resp, err := client.Get(p.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ipStr := strings.TrimSpace(string(body))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, &base.NotRetrievedError{}
	}
	return &base.ScoredIP{ip, scoreOfTLS(resp.TLS)}, nil
}

func (p HTTPDetector) String() string {
	return fmt.Sprintf("%s", p.URL)
}
