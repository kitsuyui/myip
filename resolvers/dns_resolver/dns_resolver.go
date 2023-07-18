// DNS Resolver
// Resolve IP address by DNS query
package dns_resolver

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"

	"github.com/kitsuyui/myip/resolvers/base"
)

type DNSDetector struct {
	LookupDomainName string `json:"name"`
	Resolver         string `json:"server"`
	QueryType        string `json:"queryType"`
}

func (p DNSDetector) RetrieveIP() (*base.ScoredIP, error) {
	if strings.ToUpper(p.QueryType) == "TXT" {
		return p.RetrieveIPByTXTRecord()
	}
	return p.RetrieveIPByARecord()
}

func (p DNSDetector) RetrieveIPByARecord() (*base.ScoredIP, error) {
	c := dns.Client{}
	m := dns.Msg{}
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
	return &base.ScoredIP{IP: arecord.A, Score: 1.0}, nil
}

func (p DNSDetector) RetrieveIPByTXTRecord() (*base.ScoredIP, error) {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion(p.LookupDomainName, dns.TypeTXT)
	result, _, err := c.Exchange(&m, p.Resolver)
	if err != nil {
		return nil, err
	}
	if len(result.Answer) == 0 {
		return nil, &base.NotRetrievedError{}
	}
	txtRecord := result.Answer[0].(*dns.TXT)
	if txtRecord.Txt == nil {
		return nil, &base.NotRetrievedError{}
	}
	ip := net.ParseIP(txtRecord.Txt[0])
	return &base.ScoredIP{IP: ip, Score: 1.0}, nil
}

func (p DNSDetector) String() string {
	qt := "A"
	if strings.ToUpper(p.QueryType) == "TXT" {
		qt = "TXT"
	}
	return fmt.Sprintf("%s,%s,%s", qt, p.LookupDomainName, p.Resolver)
}
