package stun_resolver

import (
	"fmt"
	"net"

	"github.com/gortc/stun"
	"github.com/kitsuyui/myip/base"
)

type STUNDetector struct {
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
}

func (p STUNDetector) RetrieveIP() (*base.ScoredIP, error) {
	var ip net.IP
	c, err := stun.Dial(p.Protocol, p.Host)
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
	return &base.ScoredIP{ip, 1.0}, nil
}

func (p STUNDetector) String() string {
	return fmt.Sprintf("%s", p.Host)
}
