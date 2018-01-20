package http_resolver

import (
	"testing"
)

func TestHTTPSuccess(t *testing.T) {
	h := HTTPDetector{URL: "http://whatismyip.akamai.com/"}
	ip, err := h.RetrieveIP()
	if err != nil {
		t.Errorf("Should be succeed")
	}
	if ip == nil {
		t.Errorf("IP must not nil")
	}
}

func TestHTTPFail(t *testing.T) {
	h := HTTPDetector{URL: "http://127.0.0.1/"}
	ip, err := h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestHTTPFail2(t *testing.T) {
	h := HTTPDetector{URL: "http://example.com/"}
	ip, err := h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}
