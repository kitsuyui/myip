// DNS Resolver
// Resolve IP address by DNS query
package dnsresolver

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/kitsuyui/myip/resolvers/base"
)

// Supported QueryType values for DNSDetector.
const (
	QueryTypeA   = "A"
	QueryTypeTXT = "TXT"
)

// dnsTimeout limits each DNS exchange so goroutines in base.RetrieveIPWithScoring
// cannot outlive the caller's context deadline by more than this duration.
const dnsTimeout = 5 * time.Second

type DNSDetector struct {
	LookupDomainName string `json:"name"`
	Resolver         string `json:"server"`
	QueryType        string `json:"queryType"`
}

func (p DNSDetector) RetrieveIP() (*base.ScoredIP, error) {
	switch strings.ToUpper(p.QueryType) {
	case QueryTypeTXT:
		return p.RetrieveIPByTXTRecord()
	case QueryTypeA, "":
		return p.RetrieveIPByARecord()
	default:
		return nil, &base.NotRetrievedError{Message: fmt.Sprintf("unsupported QueryType: %q (supported: %q, %q)", p.QueryType, QueryTypeA, QueryTypeTXT)}
	}
}

func (p DNSDetector) RetrieveIPByARecord() (*base.ScoredIP, error) {
	c := dns.Client{Timeout: dnsTimeout}
	m := dns.Msg{}
	m.SetQuestion(p.LookupDomainName, dns.TypeA)
	result, _, err := c.Exchange(&m, p.Resolver)
	if err != nil {
		return nil, err
	}
	if len(result.Answer) == 0 {
		return nil, &base.NotRetrievedError{}
	}
	arecord, ok := result.Answer[0].(*dns.A)
	if !ok || arecord.A == nil {
		return nil, &base.NotRetrievedError{}
	}
	return &base.ScoredIP{IP: arecord.A, Score: 1.0}, nil
}

func (p DNSDetector) RetrieveIPByTXTRecord() (*base.ScoredIP, error) {
	c := dns.Client{Timeout: dnsTimeout}
	m := dns.Msg{}
	m.SetQuestion(p.LookupDomainName, dns.TypeTXT)
	result, _, err := c.Exchange(&m, p.Resolver)
	if err != nil {
		return nil, err
	}
	if len(result.Answer) == 0 {
		return nil, &base.NotRetrievedError{}
	}
	txtRecord, ok := result.Answer[0].(*dns.TXT)
	if !ok || len(txtRecord.Txt) == 0 {
		return nil, &base.NotRetrievedError{}
	}
	ip := net.ParseIP(txtRecord.Txt[0])
	if ip == nil {
		return nil, &base.NotRetrievedError{Message: fmt.Sprintf("TXT record %q is not a valid IP address", txtRecord.Txt[0])}
	}
	return &base.ScoredIP{IP: ip, Score: 1.0}, nil
}

func (p DNSDetector) String() string {
	qt := QueryTypeA
	if strings.ToUpper(p.QueryType) == QueryTypeTXT {
		qt = QueryTypeTXT
	}
	return fmt.Sprintf("%s,%s,%s", qt, p.LookupDomainName, p.Resolver)
}

func (p DNSDetector) TypeName() string {
	return "dns"
}
