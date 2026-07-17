package main

import (
	"bytes"
	"math"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/kitsuyui/myip/resolvers/base"
	"github.com/kitsuyui/myip/resolvers/targets"
)

type fakeRetriever struct {
	ip    string
	score float64
	delay time.Duration
	err   error
}

func (f fakeRetriever) RetrieveIP() (*base.ScoredIP, error) {
	time.Sleep(f.delay)
	if f.err != nil {
		return nil, f.err
	}
	return &base.ScoredIP{IP: net.ParseIP(f.ip), Score: f.score}, nil
}

func (f fakeRetriever) String() string {
	return f.ip
}

func TestPickUpReturnsErrorWhenNoResultExceedsThreshold(t *testing.T) {
	retrievers := []base.ScoredIPRetrievable{
		{IPRetrievable: fakeRetriever{ip: "192.0.2.1", score: 1.0}, Weight: 1.0},
		{IPRetrievable: fakeRetriever{ip: "192.0.2.2", score: 1.0}, Weight: 1.0},
	}
	done := make(chan error, 1)
	go func() {
		_, err := pickUpFirstItemThatExceededThreshold(retrievers, time.Second, 1.1)
		done <- err
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected error")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("pickUpFirstItemThatExceededThreshold did not return")
	}
}

func TestPickUpRequiresScoreToExceedThreshold(t *testing.T) {
	retrievers := []base.ScoredIPRetrievable{
		{IPRetrievable: fakeRetriever{ip: "192.0.2.1", score: 1.0}, Weight: 1.0},
	}

	_, err := pickUpFirstItemThatExceededThreshold(retrievers, time.Second, 1.0)
	if err == nil {
		t.Fatal("expected error")
	}
	if _, ok := err.(*base.NotRetrievedError); !ok {
		t.Fatalf("expected NotRetrievedError, got %T", err)
	}
}

func TestValidateThreshold(t *testing.T) {
	cases := []struct {
		name      string
		threshold float64
		wantErr   bool
	}{
		{name: "below range", threshold: -0.01, wantErr: true},
		{name: "above range", threshold: 1.01, wantErr: true},
		{name: "minimum boundary", threshold: 0, wantErr: false},
		{name: "maximum boundary", threshold: 1, wantErr: false},
		{name: "nan threshold", threshold: math.NaN(), wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateThreshold(tc.threshold)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error for threshold %v", tc.threshold)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("did not expect error for threshold %v: %v", tc.threshold, err)
			}
		})
	}
}

func TestPickWinnerChoosesHighestScoreDeterministically(t *testing.T) {
	scores := map[string]float64{
		"192.0.2.1": 1.0,
		"192.0.2.2": 2.0,
	}

	for i := 0; i < 100; i++ {
		sip, ok := pickWinner(scores, 3.0, 0.2)
		if !ok {
			t.Fatal("expected winner")
		}
		if sip.IP.String() != "192.0.2.2" {
			t.Fatalf("unexpected IP: %s", sip.IP.String())
		}
	}
}

func TestPickWinnerBreaksScoreTiesByIP(t *testing.T) {
	scores := map[string]float64{
		"192.0.2.2": 1.0,
		"192.0.2.1": 1.0,
	}

	for i := 0; i < 100; i++ {
		sip, ok := pickWinner(scores, 2.0, 0.2)
		if !ok {
			t.Fatal("expected winner")
		}
		if sip.IP.String() != "192.0.2.1" {
			t.Fatalf("unexpected IP: %s", sip.IP.String())
		}
	}
}

func TestPickUpDoesNotPanicOnLateResultAfterWinner(t *testing.T) {
	retrievers := []base.ScoredIPRetrievable{
		{IPRetrievable: fakeRetriever{ip: "192.0.2.1", score: 1.0}, Weight: 1.0},
		{IPRetrievable: fakeRetriever{ip: "192.0.2.2", score: 1.0, delay: 20 * time.Millisecond}, Weight: 1.0},
	}
	sip, err := pickUpFirstItemThatExceededThreshold(retrievers, time.Second, 0.4)
	if err != nil {
		t.Fatalf("expected winner: %v", err)
	}
	if sip.IP.String() != "192.0.2.1" {
		t.Fatalf("unexpected IP: %s", sip.IP.String())
	}

	time.Sleep(50 * time.Millisecond)
}

func TestUsageTextDoesNotExposeRedundantNewlineFlag(t *testing.T) {
	if strings.Contains(usageText, "--newline") {
		t.Fatalf("usage text must not expose redundant newline flag: %q", usageText)
	}
	if !strings.Contains(usageText, "--no-newline") {
		t.Fatalf("usage text must expose no-newline flag: %q", usageText)
	}
}

func TestWriteIP(t *testing.T) {
	t.Run("default output ends with newline", func(t *testing.T) {
		var buf bytes.Buffer
		if err := writeIP(&buf, "192.0.2.1", false); err != nil {
			t.Fatalf("writeIP returned error: %v", err)
		}
		if got := buf.String(); got != "192.0.2.1\n" {
			t.Fatalf("unexpected output: %q", got)
		}
	})

	t.Run("no-newline output omits trailing newline", func(t *testing.T) {
		var buf bytes.Buffer
		if err := writeIP(&buf, "192.0.2.1", true); err != nil {
			t.Fatalf("writeIP returned error: %v", err)
		}
		if got := buf.String(); got != "192.0.2.1" {
			t.Fatalf("unexpected output: %q", got)
		}
	})
}

func TestSelectRetrievablesUsesIPv6OnlyTargetsWithoutIPv4Fallback(t *testing.T) {
	selected := selectRetrievables(false, true)
	ipv6Only := targets.IPv6Retrievables()
	allTargets := targets.IPRetrievables()

	if len(selected) == 0 {
		t.Fatal("expected IPv6-capable targets")
	}
	if !slices.Equal(selected, ipv6Only) {
		t.Fatal("expected --ipv6 selection to match IPv6-capable targets exactly")
	}
	if len(selected) >= len(allTargets) {
		t.Fatal("expected --ipv6 selection to be narrower than the default target set")
	}
	for _, retrievable := range selected {
		if !retrievable.IPv6 {
			t.Fatalf("found non-IPv6 target in --ipv6 selection: %+v", retrievable)
		}
	}
}

func TestReadmeDocumentsIPv6OnlyContractWithoutFallback(t *testing.T) {
	readmePath := filepath.Join("..", "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	text := string(content)

	if !strings.Contains(text, "does not fall back to IPv4") {
		t.Fatal("README must document that -6 does not fall back to IPv4")
	}
	if strings.Contains(text, "fallbacks to IPv4") {
		t.Fatal("README must not claim that -6 falls back to IPv4")
	}
}
