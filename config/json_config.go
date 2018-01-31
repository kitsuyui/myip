package config

import (
	"encoding/json"
	"time"

	"github.com/kitsuyui/myip/base"
	"github.com/kitsuyui/myip/dns_resolver"
	"github.com/kitsuyui/myip/http_resolver"
	"github.com/kitsuyui/myip/stun_resolver"
)

type ConfigJSON struct {
	Version string            `json:"version"`
	Sources []json.RawMessage `json:"sources"`
}

type ConfigJSONUnknownType struct {
	Type    string  `json:"type"`
	Weight  float64 `json:"weight"`
	Timeout float64 `json:"timeout"`
}

const defaultTimeout time.Duration = time.Millisecond * 2000

func (c ConfigJSONUnknownType) timeout(base time.Duration) time.Duration {
	if base == 0 {
		return defaultTimeout
	}
	return time.Duration(c.Timeout) * time.Second
}

func parseScoredIPRetrievableFromJSON(j json.RawMessage) (*base.ScoredIPRetrievable, error) {
	var ut ConfigJSONUnknownType
	type scoredipr = base.ScoredIPRetrievable
	if err := json.Unmarshal(j, &ut); err != nil {
		return nil, err
	}
	switch ut.Type {
	case "http":
		var r http_resolver.HTTPDetector
		if err := json.Unmarshal(j, &r); err != nil {
			return nil, err
		}
		r.Timeout = ut.timeout(r.Timeout)
		return &scoredipr{r, ut.Type, ut.Weight}, nil
	case "https":
		var r http_resolver.HTTPDetector
		if err := json.Unmarshal(j, &r); err != nil {
			return nil, err
		}
		r.Timeout = ut.timeout(r.Timeout)
		return &scoredipr{r, ut.Type, ut.Weight}, nil
	case "dns":
		var r dns_resolver.DNSDetector
		if err := json.Unmarshal(j, &r); err != nil {
			return nil, err
		}
		r.Timeout = ut.timeout(r.Timeout)
		return &scoredipr{r, ut.Type, ut.Weight}, nil
	case "stun":
		var s stun_resolver.STUNDetector
		if err := json.Unmarshal(j, &s); err != nil {
			return nil, err
		}
		s.Timeout = ut.timeout(s.Timeout)
		return &scoredipr{s, ut.Type, ut.Weight}, nil
	default:
		return nil, base.ConfigError{}
	}
}

func ScoredIPRetrievablesFromJSON(data []byte) ([]base.ScoredIPRetrievable, error) {
	var j ConfigJSON
	err := json.Unmarshal(data, &j)
	if err != nil {
		return nil, err
	}
	var sir []base.ScoredIPRetrievable
	for _, s := range j.Sources {
		ut, err := parseScoredIPRetrievableFromJSON(s)
		if err != nil {
			return nil, err
		}
		sir = append(sir, *ut)
	}
	return sir, nil
}
