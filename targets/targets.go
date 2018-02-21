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
		scored{IPRetrievable: http{URL: "http://ipecho.net/plain"}, Weight: 0.5, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://ipecho.net/plain"}, Weight: 0.5, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://inet-ip.info/ip"}, Weight: 0.5, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://eth0.me/"}, Weight: 0.5, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://ipcheck.ieserver.net/"}, Weight: 0.5, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://ifconfig.me/ip"}, Weight: 0.5, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://smart-ip.net/myip"}, Weight: 0.5, IPv4: true, IPv6: true},
		scored{IPRetrievable: http{URL: "http://whatismyip.akamai.com/"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "http://checkip.amazonaws.com/"}, Weight: 1.0, IPv4: true, IPv6: false},

		scored{IPRetrievable: http{URL: "https://bot.whatismyipaddress.com/"}, Weight: 3.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: http{URL: "https://icanhazip.com/"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://wgetip.com/"}, Weight: 3.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: http{URL: "https://ident.me/"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://4.ifcfg.me/ip"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://ip.tyk.nu/"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://tnx.nl/ip"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://l2.io/ip"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://api.ipify.org/"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://myexternalip.com/raw"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://icanhazip.com"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://ifcfg.me/ip"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://ifconfig.io/ip"}, Weight: 3.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: http{URL: "https://ifconfig.co/ip"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://ipinfo.io/ip"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://wtfismyip.com/text"}, Weight: 3.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: http{URL: "https://secure.internode.on.net/webtools/showmyip?textonly=1"}, Weight: 3.0, IPv4: true, IPv6: false},

		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver1.opendns.com:53", QueryType: "A"}, Weight: 2.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver2.opendns.com:53", QueryType: "A"}, Weight: 2.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver3.opendns.com:53", QueryType: "A"}, Weight: 2.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "myip.opendns.com.", Resolver: "resolver4.opendns.com:53", QueryType: "A"}, Weight: 2.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "whoami.akamai.net.", Resolver: "ns1-1.akamaitech.net:53", QueryType: "A"}, Weight: 2.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "whoami.ultradns.net.", Resolver: "pdns1.ultradns.net:53", QueryType: "A"}, Weight: 2.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: dns{LookupDomainName: "o-o.myaddr.l.google.com.", Resolver: "ns1.google.com:53", QueryType: "TXT"}, Weight: 2.0, IPv4: true, IPv6: true},

		scored{IPRetrievable: stun{Host: "stun.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun1.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun2.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun3.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun4.l.google.com:19302", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.services.mozilla.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stunserver.org:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.aa.net.uk:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.a-mm.tv:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.hoiio.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.unseen.is:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.acrobits.cz:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.petcube.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.voxox.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.ipcomms.net:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.voip.blackberry.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.sip.us:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "numb.viagenie.ca:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: stun{Host: "stun.ssl7.net:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.stunprotocol.org:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: stun{Host: "stun.dus.net:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: true},
		scored{IPRetrievable: stun{Host: "stun.antisip.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.avigora.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.3cx.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.avigora.fr:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.advfn.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.linphone.org:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.12connect.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.voipgate.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.cope.es:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.it1.hr:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.bluesip.net:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.sippeer.dk:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.solcon.nl:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.srce.hr:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.miwifi.com:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.sovtest.ru:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.sipnet.net:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.sipnet.ru:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
		scored{IPRetrievable: stun{Host: "stun.uls.co.za:3478", Protocol: "udp"}, Weight: 1.0, IPv4: true, IPv6: false},
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
