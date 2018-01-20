package dns_resolver

import (
	"testing"
)

func TestDNSARecordSuccess(t *testing.T) {
	d := DNSDetector{LookupDomainName: "myip.opendns.com.", Resolver: "resolver1.opendns.com:53"}
	ip, err := d.RetrieveIP()
	if err != nil {
		t.Errorf("DNS A Record Query should be succeed")
	}
	if ip == nil {
		t.Errorf("IP must not nil")
	}
}

func TestDNSARecordFailWhenDNSServerIsDownOrNotExists(t *testing.T) {
	d := DNSDetector{LookupDomainName: "example.com.", Resolver: "127.0.0.1:53"}
	ip, err := d.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestDNSARecordFailWhenInvalidLookup(t *testing.T) {
	d := DNSDetector{LookupDomainName: "dummy.opendns.com.", Resolver: "resolver1.opendns.com:53"}
	ip, err := d.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestDNSTXTRecordSuccess(t *testing.T) {
	d := DNSDetector{LookupDomainName: "o-o.myaddr.l.google.com.", Resolver: "ns1.google.com:53", QueryType: "TXT"}
	ip, err := d.RetrieveIP()
	if err != nil {
		t.Errorf("DNS TXT Record Query should be succeed.")
	}
	if ip == nil {
		t.Errorf("IP must not nil")
	}
}

func TestDNSTXTRecordFailWhenDNSServerIsDownOrNotExists(t *testing.T) {
	d := DNSDetector{LookupDomainName: "o-o.myaddr.l.google.com.", Resolver: "127.0.0.1:53", QueryType: "TXT"}
	ip, err := d.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestDNSTXTRecordFailWhenInvalidLookup(t *testing.T) {
	d := DNSDetector{LookupDomainName: "dummy.myaddr.l.google.com.", Resolver: "ns1.google.com:53", QueryType: "TXT"}
	ip, err := d.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}
