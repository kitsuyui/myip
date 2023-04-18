package base

import (
	"context"
	"net"
)

type NotRetrievedError struct {
	Message string
}

func (n NotRetrievedError) Error() string {
	if n.Message != "" {
		return n.Message
	}
	return "No Answer"
}

type TimeoutError struct {
	Message string
}

func (n TimeoutError) Error() string {
	if n.Message != "" {
		return n.Message
	}
	return "Timeout"
}

type ScoredIP struct {
	IP    net.IP
	Score float64
}

type ScoredIPWithMaxScore struct {
	IP       net.IP
	Score    float64
	MaxScore float64
}

type IPRetrievable interface {
	RetrieveIP() (*ScoredIP, error)
	String() string
}

type ScoredIPRetrievable struct {
	IPRetrievable
	Weight float64
	IPv4   bool
	IPv6   bool
}

func (s ScoredIP) addWeight(weight float64) ScoredIP {
	return ScoredIP{s.IP, s.Score * weight}
}

func (p ScoredIPRetrievable) RetrieveIPWithScoring(ctx context.Context) (*ScoredIPWithMaxScore, error) {
	type Result struct {
		ScoredIP *ScoredIP
		Err      error
	}
	c := make(chan Result)
	go func() {
		scoredIP, err := p.RetrieveIP()
		c <- Result{scoredIP, err}
	}()
	select {
	case <-ctx.Done():
		return nil, &TimeoutError{}
	case r := <-c:
		if r.Err == nil {
			return &ScoredIPWithMaxScore{
				IP:       r.ScoredIP.IP,
				Score:    p.Weight * r.ScoredIP.Score,
				MaxScore: p.Weight,
			}, nil
		}
		return nil, r.Err
	}
}
