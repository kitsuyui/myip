// STUN Resolver
// Resolve IP address by STUN protocol
// c.f. https://en.wikipedia.org/wiki/STUN
package stun_resolver

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"

	"github.com/kitsuyui/myip/resolvers/base"
	"gortc.io/stun"
)

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
	if scheme == stun.SchemeSecure {
		cfg := &tls.Config{
			ServerName: uri.Host,
		}
		conn, err = tls.Dial(p.Protocol, address, cfg)
		if err != nil {
			return nil, &base.NotRetrievedError{}
		}
	} else {
		conn, err = net.Dial(p.Protocol, address)
		if err != nil {
			return nil, &base.NotRetrievedError{}
		}
	}
	c, err := stun.NewClient(conn)
	if err != nil {
		return nil, &base.NotRetrievedError{}
	}
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
	defer c.Close()
	if scheme == stun.SchemeSecure {
		return &base.ScoredIP{IP: ip, Score: 1.0}, nil
	} else {
		return &base.ScoredIP{IP: ip, Score: 0.1}, nil
	}
}

func (p STUNDetector) String() string {
	return fmt.Sprintf("%s", p.Host)
}
