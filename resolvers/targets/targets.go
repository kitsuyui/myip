package targets

import (
	"github.com/kitsuyui/myip/resolvers/base"
	"github.com/kitsuyui/myip/resolvers/dns_resolver"
	"github.com/kitsuyui/myip/resolvers/http_resolver"
	"github.com/kitsuyui/myip/resolvers/stun_resolver"
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

		// Umbrella DNS resolvers
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "208.67.222.222:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "208.67.220.220:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},

		scored{IPRetrievable: dns{LookupDomainName: "whoami.akamai.net.", Resolver: "ns1-1.akamaitech.net:53", QueryType: "A"}, Weight: 1.0, IPv4: true, IPv6: false},
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
		scored{IPRetrievable: stun{Host: "stun:stun.antisip.com:3478", Protocol: "tcp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.linphone.org:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.voipgate.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.cope.es:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.solcon.nl:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun:stun.uls.co.za:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},

		scored{IPRetrievable: stun{Host: "stun:stunserver2025.stunprotocol.org:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},

		// Cloudflare STUN servers
		scored{IPRetrievable: stun{Host: "stun:stun.cloudflare.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		// Facebook / Meta STUN servers
		scored{IPRetrievable: stun{Host: "stun:stun.fbsbx.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
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
