package stun_resolver

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/kitsuyui/myip/base"
)

func TestSTUNSuccess(t *testing.T) {
	h := STUNDetector{Host: "stun.l.google.com:19302", Protocol: "udp"}
	ip, err := h.RetrieveIP()
	if err != nil {
		t.Errorf("Should be succeed")
	}
	if ip == nil {
		t.Errorf("IP must not nil")
	}
}

func TestSTUNFail(t *testing.T) {
	ctx := context.Background()
	timeout := 500 * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	h := STUNDetector{Host: "127.0.0.1:1000", Protocol: "udp"}
	var err error
	var ip net.IP
	type Result struct {
		IP  net.IP
		Err error
	}
	c := make(chan Result)
	go func() {
		ip, err := h.RetrieveIP()
		c <- Result{ip, err}
	}()
	select {
	case <-ctx.Done():
		err = &base.TimeoutError{}
	case r := <-c:
		ip = r.IP
		err = r.Err
	}
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestSTUNInvalidAddress(t *testing.T) {
	h := STUNDetector{Host: "<>", Protocol: "udp"}
	ip, err := h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestGetString(t *testing.T) {
	STUNDetector{Host: "stun.l.google.com:19302", Protocol: "udp"}.String()
}
