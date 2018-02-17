package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/kitsuyui/myip/base"
	"github.com/kitsuyui/myip/dns_resolver"
	"github.com/kitsuyui/myip/http_resolver"
	"github.com/kitsuyui/myip/stun_resolver"
	"github.com/kitsuyui/myip/targets"
)

var version string

func typeName(ipr interface{}) string {
	switch ipr.(type) {
	case http_resolver.HTTPDetector:
		return "http"
	case dns_resolver.DNSDetector:
		return "dns"
	case stun_resolver.STUNDetector:
		return "stun"
	}
	return ""
}

func pickupMaxScore(siprs []base.ScoredIPRetrievable, timeout time.Duration) (*base.ScoredIP, error) {
	wg := new(sync.WaitGroup)
	sumOfWeight := 0.0
	m := map[string]float64{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	logger := log.Logger{}
	if *verboseMode {
		logger.SetOutput(os.Stderr)
	} else {
		logger.SetOutput(ioutil.Discard)
	}
	c := make(chan base.ScoredIP)
	defer close(c)
	for _, sipr := range siprs {
		wg.Add(1)
		sumOfWeight += sipr.Weight
		go func(sipr base.ScoredIPRetrievable) {
			sip, err := sipr.RetriveIPWithScoring(ctx)
			if err != nil {
				logger.Printf("Error:%s\ttype:%s\tweight:%1.1f\t%s", err, typeName(sipr.IPRetrievable), sipr.Weight, sipr.String())
				wg.Done()
				return
			}
			logger.Printf("IP:%s\ttype:%s\tweight:%1.1f\t%s", sip.IP.String(), typeName(sipr.IPRetrievable), sipr.Weight, sipr.String())
			c <- *sip
			wg.Done()
		}(sipr)
	}
	go func() {
		for sip := range c {
			m[sip.IP.String()] += sip.Score
		}
	}()
	wg.Wait()
	maxKey, maxVal := pickMapMaxItem(m)
	return &base.ScoredIP{net.ParseIP(maxKey), maxVal / sumOfWeight}, nil
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

var useNewline = flag.Bool("newline", false, "Show IP with newline.")
var cmdVersion = flag.Bool("version", false, "Show version.")
var verboseMode = flag.Bool("verbose", false, "Verbose mode.")
var timeout = flag.Duration("timeout", 1*time.Second, "Timeout duration.")
var scoreThreshold = flag.Float64("threshold", 0.5, "Threshold that should be exceeded by top weighted votes.")

func init() {
	flag.BoolVar(useNewline, "n", false, "Show IP with newline.")
	flag.BoolVar(cmdVersion, "V", false, "Show version.")
	flag.BoolVar(verboseMode, "v", false, "Verbose mode.")
	flag.DurationVar(timeout, "t", 1*time.Second, "Timeout duration.")
	flag.Float64Var(scoreThreshold, "T", 0.5, "Threshold that should be exceeded by top weighted votes.")
}

func main() {
	flag.Parse()
	if *cmdVersion {
		println(version)
		return
	}
	sir := targets.IPRetrievables()
	sip, err := pickupMaxScore(sir, *timeout)
	if err == nil && sip.Score >= *scoreThreshold {
		if *useNewline {
			println(sip.IP.String())
		} else {
			print(sip.IP.String())
		}
	} else {
		os.Exit(1)
	}
}
