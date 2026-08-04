package targets

import (
	"slices"
	"testing"

	"github.com/kitsuyui/myip/resolvers/base"
)

func TestIPv4RetrievablesMatchesIPv4CapabilitiesInOrder(t *testing.T) {
	want := filterRetrievables(IPRetrievables(), func(retrievable base.ScoredIPRetrievable) bool {
		return retrievable.IPv4
	})

	got := IPv4Retrievables()

	if !slices.EqualFunc(got, want, equalScoredIPRetrievable) {
		t.Fatalf("IPv4Retrievables() mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestIPv6RetrievablesMatchesIPv6CapabilitiesInOrder(t *testing.T) {
	want := filterRetrievables(IPRetrievables(), func(retrievable base.ScoredIPRetrievable) bool {
		return retrievable.IPv6
	})

	got := IPv6Retrievables()

	if !slices.EqualFunc(got, want, equalScoredIPRetrievable) {
		t.Fatalf("IPv6Retrievables() mismatch\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestVariantSelectorsIncludeDualStackRetrievables(t *testing.T) {
	ipv4Retrievables := IPv4Retrievables()
	ipv6Retrievables := IPv6Retrievables()

	foundDualStack := false
	for _, retrievable := range IPRetrievables() {
		if !retrievable.IPv4 || !retrievable.IPv6 {
			continue
		}
		foundDualStack = true

		if !containsRetrievable(ipv4Retrievables, retrievable) {
			t.Fatalf("dual-stack retrievable missing from IPv4 selection: %s", retrievable.String())
		}
		if !containsRetrievable(ipv6Retrievables, retrievable) {
			t.Fatalf("dual-stack retrievable missing from IPv6 selection: %s", retrievable.String())
		}
	}

	if !foundDualStack {
		t.Fatal("expected at least one dual-stack retrievable")
	}
}

func filterRetrievables(all []base.ScoredIPRetrievable, keep func(base.ScoredIPRetrievable) bool) []base.ScoredIPRetrievable {
	var filtered []base.ScoredIPRetrievable
	for _, retrievable := range all {
		if keep(retrievable) {
			filtered = append(filtered, retrievable)
		}
	}
	return filtered
}

func containsRetrievable(all []base.ScoredIPRetrievable, want base.ScoredIPRetrievable) bool {
	return slices.ContainsFunc(all, func(got base.ScoredIPRetrievable) bool {
		return equalScoredIPRetrievable(got, want)
	})
}

func equalScoredIPRetrievable(a, b base.ScoredIPRetrievable) bool {
	return a.String() == b.String() &&
		a.Weight == b.Weight &&
		a.IPv4 == b.IPv4 &&
		a.IPv6 == b.IPv6
}
