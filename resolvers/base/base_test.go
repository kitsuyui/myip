package base

import (
	"context"
	"errors"
	"testing"
	"time"
)

type contextAwareRetriever struct {
	canceled chan struct{}
}

func (r contextAwareRetriever) RetrieveIP() (*ScoredIP, error) {
	return nil, errors.New("RetrieveIP fallback should not be called")
}

func (r contextAwareRetriever) RetrieveIPWithContext(ctx context.Context) (*ScoredIP, error) {
	<-ctx.Done()
	close(r.canceled)
	return nil, ctx.Err()
}

func (r contextAwareRetriever) String() string {
	return "context-aware"
}

func TestRetrieveIPWithScoringPropagatesContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	retriever := contextAwareRetriever{canceled: make(chan struct{})}
	_, err := ScoredIPRetrievable{IPRetrievable: retriever, Weight: 1.0}.RetrieveIPWithScoring(ctx)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if _, ok := err.(*TimeoutError); !ok {
		t.Fatalf("expected TimeoutError, got %T", err)
	}

	select {
	case <-retriever.canceled:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("context was not propagated to retriever")
	}
}
