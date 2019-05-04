package http_resolver

import (
	"crypto/tls"
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

func TestGetString(t *testing.T) {
	HTTPDetector{URL: "http://example.com/"}.String()
}

func TestScoreOfTLS(t *testing.T) {
	tls13Score := scoreOfTLS(&tls.ConnectionState{Version: tls.VersionTLS13})
	tls12Score := scoreOfTLS(&tls.ConnectionState{Version: tls.VersionTLS12})
	tls11Score := scoreOfTLS(&tls.ConnectionState{Version: tls.VersionTLS11})
	tls10Score := scoreOfTLS(&tls.ConnectionState{Version: tls.VersionTLS10})
	httpScore := scoreOfTLS(nil)

	if !(tls13Score > tls12Score && tls12Score > tls11Score && tls11Score > tls10Score && tls10Score > httpScore) {
		t.Errorf("Handshake strength order is not good")
	}
}
