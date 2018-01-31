package base

import "net"

type NotRetrievedError struct {
	Message string
}

func (n NotRetrievedError) Error() string {
	if n.Message != "" {
		return n.Message
	}
	return "No Answer"
}

type ConfigError struct {
	Message string
}

func (c ConfigError) Error() string {
	if c.Message != "" {
		return c.Message
	}
	return "ConfigError"
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
	Type   string
	Weight float64
}

func (p ScoredIPRetrievable) RetriveIPWithScoring() (*ScoredIP, error) {
	ip, err := p.RetrieveIP()
	return &ScoredIP{ip, p.Weight}, err
}
