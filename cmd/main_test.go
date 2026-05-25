package main

import (
	"net"
	"testing"
	"time"

	"github.com/kitsuyui/myip/resolvers/base"
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
