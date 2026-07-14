package dnsresolver

import (
	"net"
	"testing"

	"github.com/miekg/dns"
)

func TestDNSARecordSuccess(t *testing.T) {
	server := startTestDNSServer(t, func(req *dns.Msg) []dns.RR {
		return []dns.RR{
			&dns.A{
				Hdr: dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP("203.0.113.10").To4(),
			},
		}
	})
	d := DNSDetector{LookupDomainName: "example.test.", Resolver: server}
	ip, err := d.RetrieveIP()
	if err != nil {
		t.Errorf("DNS A Record Query should be succeed")
	}
	if ip == nil || !ip.IP.Equal(net.ParseIP("203.0.113.10")) {
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
	server := startTestDNSServer(t, func(req *dns.Msg) []dns.RR {
		return nil
	})
	d := DNSDetector{LookupDomainName: "missing.example.test.", Resolver: server}
	ip, err := d.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestDNSTXTRecordSuccess(t *testing.T) {
	server := startTestDNSServer(t, func(req *dns.Msg) []dns.RR {
		return []dns.RR{
			&dns.TXT{
				Hdr: dns.RR_Header{Name: req.Question[0].Name, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
				Txt: []string{"2001:db8::5"},
			},
		}
	})
	d := DNSDetector{LookupDomainName: "txt.example.test.", Resolver: server, QueryType: "TXT"}
	ip, err := d.RetrieveIP()
	if err != nil {
		t.Errorf("DNS TXT Record Query should be succeed.")
	}
	if ip == nil || !ip.IP.Equal(net.ParseIP("2001:db8::5")) {
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
	server := startTestDNSServer(t, func(req *dns.Msg) []dns.RR {
		return nil
	})
	d := DNSDetector{LookupDomainName: "missing-txt.example.test.", Resolver: server, QueryType: "TXT"}
	ip, err := d.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}
func TestUnsupportedQueryType(t *testing.T) {
	d := DNSDetector{LookupDomainName: "example.com.", Resolver: "8.8.8.8:53", QueryType: "MX"}
	ip, err := d.RetrieveIP()
	if err == nil {
		t.Errorf("Unsupported QueryType should return an error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestGetString(t *testing.T) {
	result := DNSDetector{LookupDomainName: "dummy.myaddr.l.google.com.", Resolver: "ns1.google.com:53", QueryType: "TXT"}.String()
	tobe := "TXT,dummy.myaddr.l.google.com.,ns1.google.com:53"
	if result != tobe {
		t.Errorf("The result must be %s", tobe)
	}
	result = DNSDetector{LookupDomainName: "dummy.opendns.com.", Resolver: "resolver1.opendns.com:53"}.String()
	tobe = "A,dummy.opendns.com.,resolver1.opendns.com:53"
	if result != tobe {
		t.Errorf("The result must be %s", tobe)
	}
}

func startTestDNSServer(t *testing.T, answer func(req *dns.Msg) []dns.RR) string {
	t.Helper()

	packetConn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen udp: %v", err)
	}

	server := &dns.Server{
		PacketConn: packetConn,
		Handler: dns.HandlerFunc(func(w dns.ResponseWriter, req *dns.Msg) {
			msg := new(dns.Msg)
			msg.SetReply(req)
			msg.Answer = answer(req)
			if err := w.WriteMsg(msg); err != nil {
				t.Errorf("write dns response: %v", err)
			}
		}),
	}

	go func() {
		if err := server.ActivateAndServe(); err != nil {
			t.Logf("dns server stopped: %v", err)
		}
	}()

	t.Cleanup(func() {
		if err := server.Shutdown(); err != nil {
			t.Errorf("shutdown dns server: %v", err)
		}
	})

	return packetConn.LocalAddr().String()
}
