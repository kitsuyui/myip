package dns_resolver

import (
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"

	"../base"
)

type DNSDetector struct {
	LookupDomainName string        `json:"name"`
	Resolver         string        `json:"server"`
	Timeout          time.Duration `json:"timeout"`
	QueryType        string        `json:"queryType"`
}

func (p DNSDetector) RetrieveIP() (net.IP, error) {
	if strings.ToUpper(p.QueryType) == "TXT" {
		return p.RetrieveIPByTXTRecord()
	}
	return p.RetrieveIPByARecord()
}

func (p DNSDetector) RetrieveIPByARecord() (net.IP, error) {
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

func (p DNSDetector) RetrieveIPByTXTRecord() (net.IP, error) {
	c := dns.Client{}
	m := dns.Msg{}
	c.Timeout = p.Timeout
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
	return ip, nil
}
