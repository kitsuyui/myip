package stun_resolver

import (
	"fmt"
	"net"
	"time"

	"github.com/gortc/stun"
	"github.com/kitsuyui/myip/base"
)

type STUNDetector struct {
	Host     string `json:"host"`
	Protocol string `json:"protocol"`
}

func (p STUNDetector) RetrieveIP() (net.IP, error) {
	var ip net.IP
	c, err := stun.Dial(p.Protocol, p.Host)
	if err != nil {
		return nil, &base.NotRetrievedError{}
	}
	deadline := time.Now().AddDate(0, 0, 1)
	var err2 error
	if err := c.Do(stun.MustBuild(stun.TransactionID, stun.BindingRequest), deadline, func(res stun.Event) {
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
	return ip, nil
}

func (p STUNDetector) String() string {
	return fmt.Sprintf("%s", p.Host)
}
