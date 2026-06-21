// STUN Resolver
// Resolve IP address by STUN protocol
// c.f. https://en.wikipedia.org/wiki/STUN
package stunresolver

import (
	"crypto/tls"
	"net"
	"strconv"
	"time"

	"github.com/kitsuyui/myip/resolvers/base"
	"github.com/pion/stun"
)

// dialTimeout limits each STUN dial so goroutines in base.RetrieveIPWithScoring
// cannot outlive the caller's context deadline by more than this duration.
const dialTimeout = 5 * time.Second

type STUNDetector struct {
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
}

func (p STUNDetector) RetrieveIP() (*base.ScoredIP, error) {
	var ip net.IP
	var conn net.Conn
	var err error

	uri, err := stun.ParseURI(p.Host)
	if err != nil {
		return nil, &base.NotRetrievedError{}
	}
	scheme := uri.Scheme
	address := uri.Host + ":" + strconv.Itoa(uri.Port)
	dialer := &net.Dialer{Timeout: dialTimeout}
	if scheme == stun.SchemeTypeSTUNS {
		cfg := &tls.Config{
			ServerName: uri.Host,
		}
		conn, err = tls.DialWithDialer(dialer, p.Protocol, address, cfg)
		if err != nil {
			return nil, &base.NotRetrievedError{}
		}
	} else {
		conn, err = dialer.Dial(p.Protocol, address)
		if err != nil {
			return nil, &base.NotRetrievedError{}
		}
	}
	c, err := stun.NewClient(conn)
	if err != nil {
		conn.Close()
		return nil, &base.NotRetrievedError{}
	}
	defer c.Close()
	var err2 error
	if err := c.Do(stun.MustBuild(stun.TransactionID, stun.BindingRequest), func(res stun.Event) {
		if res.Error != nil {
			err2 = &base.NotRetrievedError{}
			return
		}
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err != nil {
			err2 = &base.NotRetrievedError{}
			return
		}
		ip = xorAddr.IP
	}); err != nil || err2 != nil {
		return nil, &base.NotRetrievedError{}
	}
	if scheme == stun.SchemeTypeSTUNS {
		return &base.ScoredIP{IP: ip, Score: 1.0}, nil
	} else {
		return &base.ScoredIP{IP: ip, Score: 0.1}, nil
	}
}

func (p STUNDetector) String() string {
	return p.Host
}

func (p STUNDetector) TypeName() string {
	return "stun"
}
