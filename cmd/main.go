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
 -T=<rate> --threshold=<rate>  			 Threshold in [0.0, 1.0] that must be exceeded by weighted votes [default: 0.6].
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
			fmt.Print(sip.IP.String())
		} else {
			fmt.Println(sip.IP.String())
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
