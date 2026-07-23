package application

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/istvzsig/matchpulse/internal/adapters/memory"
	"github.com/istvzsig/matchpulse/internal/domain"
)

// TestEventProcessorConcurrentScoreUpdates fires many score-update events
// at the same fixture concurrently. Before the atomic Update fix, this
// was prone to lost updates (two goroutines reading the same state,
// each applying their own increment, then one Save overwriting the
// other). With Update, every increment must be reflected in the total.
func TestEventProcessorConcurrentScoreUpdates(t *testing.T) {
	repository := memory.NewMemoryMatchAdapter()
	processor := NewEventProcessor(repository)

	fixtureID := "F1"
	if err := repository.Save(context.Background(), domain.NewMatchState(fixtureID)); err != nil {
		t.Fatal(err)
	}

	const eventCount = 100

	var wg sync.WaitGroup
	wg.Add(eventCount)

	var mu sync.Mutex
	successCount := 0

	for i := 1; i <= eventCount; i++ {
		go func(version int64) {
			defer wg.Done()

			event := domain.MatchEvent{
				ID:        "E",
				FixtureID: fixtureID,
				Type:      domain.EventScoreUpdate,
				TeamID:    "T1",
				Value:     1,
				Version:   version,
				Timestamp: time.Now(),
			}

			err := processor.Process(context.Background(), event)

			if err == nil {
				mu.Lock()
				successCount++
				mu.Unlock()
				return
			}

			// Out-of-order version conflicts are expected under
			// concurrent submission — the versioning scheme assumes
			// ordered delivery per fixture. Anything else is a real failure.
			if !errors.Is(err, domain.ErrOutdatedEvent) {
				t.Errorf("unexpected error: %v", err)
			}
		}(int64(i))
	}

	wg.Wait()

	state, err := repository.Get(context.Background(), fixtureID)
	if err != nil {
		t.Fatal(err)
	}

	// The key assertion: score must exactly match the number of events
	// that actually succeeded. Before the Update fix, this could fail
	// even with zero errors reported, because a lost update doesn't
	// throw an error — it just silently disappears.
	if state.Scores["T1"] != successCount {
		t.Fatalf(
			"lost update detected: %d events succeeded but score is %d",
			successCount, state.Scores["T1"],
		)
	}

	t.Logf("%d/%d events succeeded (rest rejected as out-of-order, as expected)", successCount, eventCount)
}
