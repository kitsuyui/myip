package http_resolver

import (
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
	return &base.ScoredIP{ip, 1.0}, nil
}

func (p HTTPDetector) String() string {
	return fmt.Sprintf("%s", p.URL)
}
