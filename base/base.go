package base

import "net"
import "context"

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

type IPRetrievable interface {
	RetrieveIP() (net.IP, error)
	String() string
}

type ScoredIPRetrievable struct {
	IPRetrievable
	Weight float64
}

func (p ScoredIPRetrievable) RetriveIPWithScoring(ctx context.Context) (*ScoredIP, error) {
	type Result struct {
		IP  net.IP
		Err error
	}
	c := make(chan Result)
	go func() {
		ip, err := p.RetrieveIP()
		c <- Result{ip, err}
	}()
	select {
	case <-ctx.Done():
		return nil, &TimeoutError{}
	case r := <-c:
		return &ScoredIP{r.IP, p.Weight}, r.Err
	}
}
