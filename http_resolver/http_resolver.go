package http_resolver

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kitsuyui/myip/base"
)

type HTTPDetector struct {
	URL     string        `json:"url"`
	Timeout time.Duration `json:"timeout"`
}

func (p HTTPDetector) RetrieveIP() (net.IP, error) {
	client := http.Client{Timeout: p.Timeout}
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
	return ip, nil
}

func (p HTTPDetector) String() string {
	return fmt.Sprintf("%s", p.URL)
}
