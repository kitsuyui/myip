package http_resolver

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("203.0.113.10\n"))
	}))
	defer server.Close()

	h := HTTPDetector{URL: server.URL}
	scoredIP, err := h.RetrieveIP()
	if err != nil {
		t.Errorf("Should be succeed: %v", err)
	}
	if scoredIP == nil {
		t.Fatal("IP must not nil")
	}
	if scoredIP.IP.String() != "203.0.113.10" {
		t.Errorf("unexpected IP: %s", scoredIP.IP.String())
	}
}

func TestHTTPFail(t *testing.T) {
	h := HTTPDetector{URL: "://bad-url"}
	ip, err := h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestHTTPFail2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not an ip"))
	}))
	defer server.Close()

	h := HTTPDetector{URL: server.URL}
	ip, err := h.RetrieveIP()
	if err == nil {
		t.Errorf("This should be error")
	}
	if ip != nil {
		t.Errorf("IP should be nil when error")
	}
}

func TestHTTPContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	h := HTTPDetector{URL: server.URL}
	ip, err := h.RetrieveIPWithContext(ctx)
	if err == nil {
		t.Fatal("expected context error")
	}
	if ip != nil {
		t.Fatalf("IP should be nil when context is canceled: %v", ip)
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}

func TestGetString(t *testing.T) {
	result := HTTPDetector{URL: "http://example.com/"}.String()
	if result != "http://example.com/" {
		t.Errorf("The result must be http://example.com/")
	}
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
