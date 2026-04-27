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
