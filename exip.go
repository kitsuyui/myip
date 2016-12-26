package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"./defaults"
	"github.com/bitly/go-simplejson"
	"github.com/miekg/dns"
)

//go:generate go-bindata -prefix "data/" -pkg defaults -o defaults/defaults.go data/...

const defaultTimeout time.Duration = time.Millisecond * 2000

// ScoredIP ...
type ScoredIP struct {
	IP    net.IP
	Score float64
}

// HTTPDetector ...
type HTTPDetector struct {
	URL     string
	Timeout time.Duration
}

// DNSDetector ...
type DNSDetector struct {
	LookupDomainName string
	Resolver         string
	Timeout          time.Duration
}

// IPRetrievable ...
type IPRetrievable interface {
	RetrieveIP() (net.IP, error)
}

// ScoredIPRetrievable ...
type ScoredIPRetrievable struct {
	IPRetrievable
	Weight float64
}

// NotRetrievedError ...
type NotRetrievedError struct {
	Message string
}

func (n NotRetrievedError) Error() string {
	if n.Message != "" {
		return n.Message
	}
	return "No Answer"
}

// ConfigError ...
type ConfigError struct {
	Message string
}

func (c ConfigError) Error() string {
	if c.Message != "" {
		return c.Message
	}
	return "ConfigError"
}

// RetriveIPWithScoring ...
func (p ScoredIPRetrievable) RetriveIPWithScoring() (*ScoredIP, error) {
	ip, err := p.RetrieveIP()
	return &ScoredIP{ip, p.Weight}, err
}

// RetrieveIP ...
func (p HTTPDetector) RetrieveIP() (net.IP, error) {
	client := http.Client{Timeout: p.Timeout}
	resp, err := client.Get(p.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ipStr := strings.TrimSpace(string(body))
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, &NotRetrievedError{}
	}
	return ip, nil
}

// RetrieveIP ...
func (p DNSDetector) RetrieveIP() (net.IP, error) {
	c := dns.Client{}
	m := dns.Msg{}
	c.Timeout = p.Timeout
	m.SetQuestion(p.LookupDomainName, dns.TypeA)
	result, _, err := c.Exchange(&m, p.Resolver)
	if err != nil {
		return nil, err
	}
	if len(result.Answer) == 0 {
		return nil, &NotRetrievedError{}
	}
	arecord := result.Answer[0].(*dns.A)
	if arecord.A == nil {
		return nil, &NotRetrievedError{}
	}
	return arecord.A, nil
}

// NewScoredIPRetrievableFromJSON ...
func NewScoredIPRetrievableFromJSON(s *simplejson.Json) (*ScoredIPRetrievable, error) {
	weight, err := s.Get("weight").Float64()
	if err != nil {
		weight = 1.0
	}
	if ipr, err := NewIPRetrievableFromJSON(s); err == nil {
		return &ScoredIPRetrievable{ipr, weight}, nil
	}
	return nil, err
}

// NewIPRetrievableFromJSON ...
func NewIPRetrievableFromJSON(s *simplejson.Json) (IPRetrievable, error) {
	retrievingType := s.Get("type").MustString()
	if retrievingType == "http" || retrievingType == "https" {
		url := s.Get("url").MustString()
		var timeout time.Duration
		if timeoutP, err := s.Get("timeout").Float64(); err == nil {
			timeout = time.Duration(timeoutP) * time.Second
		} else {
			timeout = defaultTimeout
		}
		return &HTTPDetector{url, timeout}, nil
	} else if retrievingType == "dns" {
		name := s.Get("name").MustString()
		server := s.Get("server").MustString()
		var timeout time.Duration
		if timeoutP, err := s.Get("timeout").Float64(); err == nil {
			timeout = time.Duration(timeoutP) * time.Second
		} else {
			timeout = defaultTimeout
		}
		return &DNSDetector{name, server, timeout}, nil
	}
	return nil, &ConfigError{"type error"}
}

func settings(data []byte) ([]*ScoredIPRetrievable, error) {
	json, err := simplejson.NewJson(data)
	if err != nil {
		return nil, &ConfigError{}
	}
	sourcesJSON := json.Get("sources").MustArray()
	sources := []*ScoredIPRetrievable{}
	for i := range sourcesJSON {
		j := json.Get("sources").GetIndex(i)
		if ipr, err := NewScoredIPRetrievableFromJSON(j); err == nil {
			sources = append(sources, ipr)
		} else {
			return nil, &ConfigError{}
		}
	}
	return sources, err
}

func defaultSettings() ([]*ScoredIPRetrievable, error) {
	if data, err := defaults.Asset("defaults.json"); err == nil {
		return settings(data)
	}
	return nil, &ConfigError{"assets are not contained"}
}

func pickupMaxScore(siprs []*ScoredIPRetrievable) (*ScoredIP, error) {
	wg := new(sync.WaitGroup)
	maxWeight := 0.0
	m := map[string]float64{}
	for _, sipr := range siprs {
		wg.Add(1)
		maxWeight += sipr.Weight
		go func(sipr *ScoredIPRetrievable) {
			if sip, err := sipr.RetriveIPWithScoring(); err == nil {
				m[sip.IP.String()] += sip.Score
			}
			wg.Done()
		}(sipr)
	}
	wg.Wait()
	maxKey, maxVal := pickMapMaxItem(m)
	return &ScoredIP{net.ParseIP(maxKey), maxVal / maxWeight}, nil
}

func pickMapMaxItem(m map[string]float64) (string, float64) {
	maxVal := 0.0
	maxKey := ""
	for k2, v2 := range m {
		if v2 > maxVal {
			maxKey = k2
			maxVal = v2
		}
	}
	return maxKey, maxVal
}

func main() {
	retrievables, err := defaultSettings()
	if err != nil {
		log.Fatal("config error")
	}
	sip, err := pickupMaxScore(retrievables)
	if err == nil && sip.Score >= 0.5 {
		fmt.Print(sip.IP.String())
	} else {
		os.Exit(1)
	}
}
