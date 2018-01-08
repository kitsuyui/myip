package dns_resolver

import (
	"net"
	"time"

	"github.com/miekg/dns"

	"../base"
)

type DNSDetector struct {
	LookupDomainName string        `json:"name"`
	Resolver         string        `json:"resolver"`
	Timeout          time.Duration `json:"timeout"`
}

func (p DNSDetector) RetrieveIP() (net.IP, error) {
	c := dns.Client{}
	m := dns.Msg{}
	c.Timeout = p.Timeout
	m.SetQuestion(p.LookupDomainName, dns.TypeA)
	result, _, err := c.Exchange(&m, p.Resolver)
	if err != nil {
		return nil, err
	}
	if len(result.Answer) == 0 {
		return nil, &base.NotRetrievedError{}
	}
	arecord := result.Answer[0].(*dns.A)
	if arecord.A == nil {
		return nil, &base.NotRetrievedError{}
	}
	return arecord.A, nil
}
