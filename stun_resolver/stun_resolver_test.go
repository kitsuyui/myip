package stun_resolver

import (
	"testing"
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
	h := STUNDetector{Host: "127.0.0.1:1000", Protocol: "udp"}
	ip, err := h.RetrieveIP()
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
