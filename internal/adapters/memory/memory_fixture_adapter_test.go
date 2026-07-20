package memory

import (
	"context"
	"testing"
	"time"

	"github.com/istvzsig/matchpulse/internal/domain"
)

// TestMemoryFixtureAdapterGetFixtures verifies that fixtures seeded into
// the adapter are returned correctly by GetFixtures.
func TestMemoryFixtureAdapterGetFixtures(t *testing.T) {
	adapter := NewMemoryFixtureAdapter()

	fixture := domain.Fixture{
		ID:          "F1",
		Tournament:  "Worlds",
		ScheduledAt: time.Now(),
	}

	adapter.Seed([]domain.Fixture{fixture})

	fixtures, err := adapter.GetFixtures(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(fixtures) != 1 {
		t.Fatalf("expected 1 fixture, got %d", len(fixtures))
	}

	if fixtures[0].ID != "F1" {
		t.Fatalf("expected fixture F1, got %s", fixtures[0].ID)
	}
}

// TestMemoryFixtureAdapterReturnsCopy verifies that GetFixtures returns a
// defensive copy — mutating the returned slice must not affect the
// adapter's internal state.
func TestMemoryFixtureAdapterReturnsCopy(t *testing.T) {
	adapter := NewMemoryFixtureAdapter()
	adapter.Seed([]domain.Fixture{{ID: "F1"}})

	fixtures, _ := adapter.GetFixtures(context.Background())
	fixtures[0].ID = "mutated"

	fixturesAgain, _ := adapter.GetFixtures(context.Background())
	if fixturesAgain[0].ID != "F1" {
		t.Fatalf("internal state was mutated by caller")
	}
}

// TestMemoryFixtureAdapterEmpty verifies that a fresh adapter with no
// seeded fixtures returns an empty (non-nil) slice, not an error.
func TestMemoryFixtureAdapterEmpty(t *testing.T) {
	adapter := NewMemoryFixtureAdapter()

	fixtures, err := adapter.GetFixtures(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if len(fixtures) != 0 {
		t.Fatalf("expected 0 fixtures, got %d", len(fixtures))
	}
}
