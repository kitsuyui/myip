package targets

import (
	"github.com/kitsuyui/myip/base"
	"github.com/kitsuyui/myip/dns_resolver"
	"github.com/kitsuyui/myip/http_resolver"
	"github.com/kitsuyui/myip/stun_resolver"
)

func IPRetrievables() []base.ScoredIPRetrievable {
	type scored = base.ScoredIPRetrievable
	type http = http_resolver.HTTPDetector
	type dns = dns_resolver.DNSDetector
	type stun = stun_resolver.STUNDetector
	return []base.ScoredIPRetrievable{
		scored{IPRetrievable: http{URL: "http://inet-ip.info/ip"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://whatismyip.akamai.com/"}, Weight: 1.0, IPv4: true, IPv6: false},

		scored{IPRetrievable: http{URL: "https://ipecho.net/plain"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://eth0.me/"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://ifconfig.me/ip"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://checkip.amazonaws.com/"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://wgetip.com/"}, Weight: 1.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: http{URL: "https://ip.tyk.nu/"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://l2.io/ip"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://api.ipify.org/"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://myexternalip.com/raw"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://icanhazip.com"}, Weight: 1.0, IPv4: true, IPv6: true}, // document: https://major.io/icanhazip-com-faq/
		scored{IPRetrievable: http{URL: "https://ifconfig.io/ip"}, Weight: 1.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: http{URL: "https://ifconfig.co/ip"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://ipinfo.io/ip"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://wtfismyip.com/text"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://secure.internode.on.net/webtools/showmyip?textonly=1"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://secure.internode.on.net/webtools/showmyip?textonly=1"}, Weight: 1.0, IPv4: true, IPv6: false},

		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver1.opendns.com:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver2.opendns.com:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver3.opendns.com:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver4.opendns.com:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "whoami.akamai.net.", Resolver: "ns1-1.akamaitech.net:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "whoami.ultradns.net.", Resolver: "pdns1.ultradns.net:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "o-o.myaddr.l.google.com.", Resolver: "ns1.google.com:53", QueryType: "TXT"}, Weight: 1.0, IPv4: true, IPv6: true},

		scored{IPRetrievable: stun{Host: "stun:stun.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun1.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun2.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun3.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun4.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.aa.net.uk:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.hoiio.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.acrobits.cz:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.voip.blackberry.com:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.sip.us:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.stunprotocol.org:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: stun{Host: "stun:stun.antisip.com:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.avigora.fr:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.linphone.org:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.12connect.com:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.voipgate.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.cope.es:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.bluesip.net:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.sippeer.dk:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.solcon.nl:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.sovtest.ru:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.sipnet.net:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.uls.co.za:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},

		scored{IPRetrievable: stun{Host: "stuns:stun.sipnet.ru:5349", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
	}
}

func IPv4Retrievables() (sir []base.ScoredIPRetrievable) {
	for _, sipr := range IPRetrievables() {
		if sipr.IPv4 && !sipr.IPv6 {
			sir = append(sir, sipr)
		}
	}
	return sir
}

func IPv6Retrievables() (sir []base.ScoredIPRetrievable) {
	for _, sipr := range IPRetrievables() {
		if sipr.IPv6 {
			sir = append(sir, sipr)
		}
	}
	return sir
}
