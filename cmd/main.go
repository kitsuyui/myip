package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"

	"log"
	"os"
	"strconv"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/kitsuyui/myip/resolvers/base"
	"github.com/kitsuyui/myip/resolvers/dns_resolver"
	"github.com/kitsuyui/myip/resolvers/http_resolver"
	"github.com/kitsuyui/myip/resolvers/stun_resolver"
	"github.com/kitsuyui/myip/resolvers/targets"
)

var version string
var verboseMode bool

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

func pickUpFirstItemThatExceededThreshold(siprs []base.ScoredIPRetrievable, timeout time.Duration, threshold float64) (*base.ScoredIP, error) {
	sumOfWeight := 0.0
	m := map[string]float64{}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	logger := log.Logger{}
	if verboseMode {
		logger.SetOutput(os.Stderr)
	} else {
		logger.SetOutput(ioutil.Discard)
	}
	c := make(chan base.ScoredIPWithMaxScore)
	defer close(c)
	for _, sipr := range siprs {
		sumOfWeight += sipr.Weight
		go func(sipr base.ScoredIPRetrievable) {
			sip, err := sipr.RetriveIPWithScoring(ctx)
			if err != nil {
				logger.Printf("Error:%s\ttype:%s\tweight:%1.1f\t%s", err, typeName(sipr.IPRetrievable), sipr.Weight, sipr.String())
				return
			}
			logger.Printf("IP:%s\ttype:%s\tweight:%1.1f\t%s", sip.IP.String(), typeName(sipr.IPRetrievable), sip.Score, sipr.String())
			c <- *sip
		}(sipr)
	}
	result := make(chan base.ScoredIP)
	defer close(result)
	go func() {
		for sip := range c {
			key := sip.IP.String()
			m[key] += sip.Score
			sumOfWeight -= (sip.MaxScore - sip.Score)
			currentScore := m[key] / sumOfWeight
			if currentScore > threshold {
				result <- base.ScoredIP{IP: sip.IP, Score: currentScore}
			}
		}
	}()
	sip := <-result
	return &sip, nil
}

var timeout = flag.Duration("timeout", 3*time.Second, "Timeout duration.")

func init() {
	flag.DurationVar(timeout, "t", 3*time.Second, "Timeout duration.")
}

func main() {

	usage := `myip

Usage:
 myip [-v | --verbose] [-4 | -6] [-T=<rate>] [-t=<duration>]
 myip (--help | --version)

Options:
 -h --help               						 Show this screen.
 -V --version            						 Show version.
 -v --verbose            						 Verbose mode.
 -4 --ipv4               						 Prefer IPv4.
 -6 --ipv6               						 Prefer IPv6.
 -n --newline            						 Show IP with newline.
 -N --no-newline         						 Show IP without newline.
 -T=<rate> --threshold=<rate>  			 Threshold that must be exceeded by weighted votes [default: 0.6].
 -t=<duration> --timeout=<duration>  Timeout [default: 3s].
`
	opts, err := docopt.ParseDoc(usage)
	if err != nil {
		log.Fatal(err)
	}
	if showVersion, _ := opts.Bool("--version"); showVersion {
		println(version)
		return
	}
	var sir []base.ScoredIPRetrievable
	if ipv4, _ := opts.Bool("--ipv4"); ipv4 {
		sir = targets.IPv4Retrievables()
	} else if ipv6, _ := opts.Bool("--ipv6"); ipv6 {
		sir = targets.IPv6Retrievables()
	} else {
		sir = targets.IPRetrievables()
	}
	if verbose, _ := opts.Bool("--verbose"); verbose {
		verboseMode = true
	}
	threshold, err := opts.Float64("--threshold")
	if err != nil {
		log.Fatal(err)
	}

	timeoutStr, err := opts.String("--timeout")
	if err != nil {
		log.Fatal(err)
	}

	var duration time.Duration
	f, err := strconv.ParseFloat(timeoutStr, 64)
	if err != nil {
		duration, err = time.ParseDuration(timeoutStr)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		duration = time.Duration(f) * time.Second
	}
	sip, err := pickUpFirstItemThatExceededThreshold(sir, duration, threshold)
	if err == nil && sip.Score >= threshold {
		fmt.Print(sip.IP.String())
		if newline, _ := opts.Bool("--newline"); newline {
			fmt.Println("")
		}
	} else {
		os.Exit(1)
	}
}
