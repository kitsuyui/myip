package main

import (
	"context"
	"fmt"
	"io"
	"math"

	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/kitsuyui/myip/resolvers/base"
	"github.com/kitsuyui/myip/resolvers/targets"
)

var version string
var verboseMode bool

const usageText = `myip

Usage:
 myip [-v | --verbose] [-4 | -6] [-N | --no-newline] [-T=<rate>] [-t=<duration>]
 myip (--help | --version)

Options:
 -h --help               						 Show this screen.
 -V --version            						 Show version.
 -v --verbose            						 Verbose mode.
 -4 --ipv4               						 Prefer IPv4.
 -6 --ipv6               						 Prefer IPv6.
 -N --no-newline         						 Show IP without newline.
 -T=<rate> --threshold=<rate>  			 Threshold in [0.0, 1.0] that must be exceeded by weighted votes [default: 0.6].
 -t=<duration> --timeout=<duration>  Timeout [default: 3s].
`

type namer interface {
	TypeName() string
}

func typeName(ipr interface{}) string {
	if n, ok := ipr.(namer); ok {
		return n.TypeName()
	}
	return fmt.Sprintf("%T", ipr)
}

func validateThreshold(threshold float64) error {
	if math.IsNaN(threshold) || math.IsInf(threshold, 0) || threshold < 0 || threshold > 1 {
		return fmt.Errorf("threshold must be between 0 and 1")
	}
	return nil
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
		logger.SetOutput(io.Discard)
	}
	type result struct {
		scoredIP base.ScoredIPWithMaxScore
		source   base.ScoredIPRetrievable
		err      error
	}
	c := make(chan result, len(siprs))
	var wg sync.WaitGroup
	for _, sipr := range siprs {
		sumOfWeight += sipr.Weight
		wg.Add(1)
		go func(sipr base.ScoredIPRetrievable) {
			defer wg.Done()
			sip, err := sipr.RetrieveIPWithScoring(ctx)
			if err != nil {
				logger.Printf("Error:%s\ttype:%s\tweight:%1.1f\t%s", err, typeName(sipr.IPRetrievable), sipr.Weight, sipr.String())
				c <- result{source: sipr, err: err}
				return
			}
			logger.Printf("IP:%s\ttype:%s\tweight:%1.1f\t%s", sip.IP.String(), typeName(sipr.IPRetrievable), sip.Score, sipr.String())
			c <- result{scoredIP: *sip, source: sipr}
		}(sipr)
	}
	go func() {
		wg.Wait()
		close(c)
	}()
	winner := func() (*base.ScoredIP, bool) {
		return pickWinner(m, sumOfWeight, threshold)
	}
	for {
		select {
		case <-ctx.Done():
			return nil, &base.TimeoutError{}
		case r, ok := <-c:
			if !ok {
				return nil, &base.NotRetrievedError{}
			}
			if r.err != nil {
				sumOfWeight -= r.source.Weight
			} else {
				key := r.scoredIP.IP.String()
				m[key] += r.scoredIP.Score
				sumOfWeight -= (r.scoredIP.MaxScore - r.scoredIP.Score)
			}
			if sip, ok := winner(); ok {
				cancel()
				return sip, nil
			}
		}
	}
}

func pickWinner(scores map[string]float64, sumOfWeight float64, threshold float64) (*base.ScoredIP, bool) {
	if sumOfWeight <= 0 {
		return nil, false
	}

	var winnerIP string
	var winnerScore float64
	found := false
	for ip, score := range scores {
		currentScore := score / sumOfWeight
		if currentScore <= threshold {
			continue
		}
		if !found || currentScore > winnerScore || currentScore == winnerScore && ip < winnerIP {
			winnerIP = ip
			winnerScore = currentScore
			found = true
		}
	}
	if !found {
		return nil, false
	}
	return &base.ScoredIP{IP: net.ParseIP(winnerIP), Score: winnerScore}, true
}

func writeIP(out io.Writer, ip string, noNewline bool) error {
	if noNewline {
		_, err := fmt.Fprint(out, ip)
		return err
	}
	_, err := fmt.Fprintln(out, ip)
	return err
}

func selectRetrievables(ipv4 bool, ipv6 bool) []base.ScoredIPRetrievable {
	if ipv4 {
		return targets.IPv4Retrievables()
	}
	if ipv6 {
		return targets.IPv6Retrievables()
	}
	return targets.IPRetrievables()
}

func main() {
	opts, err := docopt.ParseDoc(usageText)
	if err != nil {
		log.Fatal(err)
	}
	if showVersion, _ := opts.Bool("--version"); showVersion {
		println(version)
		return
	}
	ipv4, _ := opts.Bool("--ipv4")
	ipv6, _ := opts.Bool("--ipv6")
	sir := selectRetrievables(ipv4, ipv6)
	if verbose, _ := opts.Bool("--verbose"); verbose {
		verboseMode = true
	}
	threshold, err := opts.Float64("--threshold")
	if err != nil {
		log.Fatal(err)
	}
	if err := validateThreshold(threshold); err != nil {
		log.Fatal(err)
	}

	timeoutStr, err := opts.String("--timeout")
	if err != nil {
		log.Fatal(err)
	}

	duration, err := time.ParseDuration(timeoutStr)
	if err != nil {
		log.Fatal(err)
	}
	sip, err := pickUpFirstItemThatExceededThreshold(sir, duration, threshold)
	if err == nil {
		if noNewline, _ := opts.Bool("--no-newline"); noNewline {
			if err := writeIP(os.Stdout, sip.IP.String(), true); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := writeIP(os.Stdout, sip.IP.String(), false); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		} else {
			fmt.Fprintln(os.Stderr, "error: consensus threshold not reached")
		}
		os.Exit(1)
	}
}
