package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"

	"github.com/kitsuyui/myip/base"
	"github.com/kitsuyui/myip/config"
	"github.com/kitsuyui/myip/defaults"
)

var version string

//go:generate go-bindata -prefix "data/" -pkg defaults -o defaults/defaults.go data/...

func pickupMaxScore(siprs []base.ScoredIPRetrievable) (*base.ScoredIP, error) {
	wg := new(sync.WaitGroup)
	maxWeight := 0.0
	m := map[string]float64{}
	for _, sipr := range siprs {
		wg.Add(1)
		maxWeight += sipr.Weight
		go func(sipr base.ScoredIPRetrievable) {
			if sip, err := sipr.RetriveIPWithScoring(); err == nil {
				m[sip.IP.String()] += sip.Score
			}
			wg.Done()
		}(sipr)
	}
	wg.Wait()
	maxKey, maxVal := pickMapMaxItem(m)
	return &base.ScoredIP{net.ParseIP(maxKey), maxVal / maxWeight}, nil
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

func init() {
	flag.BoolVar(useNewline, "n", false, "Show IP with newline.")
	flag.BoolVar(cmdVersion, "V", false, "Show version.")
}

func main() {
	flag.Parse()
	if *cmdVersion {
		println(version)
		return
	}
	data, err := defaults.Asset("defaults.json")
	if err != nil {
		log.Fatal(err)
	}
	sir, err := config.ScoredIPRetrievablesFromJSON(data)
	if err != nil {
		log.Fatal(err)
	}
	sip, err := pickupMaxScore(sir)
	if err == nil && sip.Score >= 0.5 {
		if *useNewline {
			println(sip.IP.String())
		} else {
			print(sip.IP.String())
		}
	} else {
		os.Exit(1)
	}
}
