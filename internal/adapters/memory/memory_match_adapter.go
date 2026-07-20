package memory

import (
	"context"
	"sync"

	"github.com/istvzsig/matchpulse/internal/domain"
	"github.com/istvzsig/matchpulse/internal/ports"
)

// MemoryMatchAdapter is an in-memory implementation of ports.MatchRepository.
// It exists so the application layer can be developed and tested without a
// real database - a persistent adapter (e.g. Postgres/Redis) can implement
// the same port later without any changes above this layer.
type MemoryMatchAdapter struct {
	matches map[string]domain.MatchState
	mu      sync.RWMutex // guards concurrent access to matches
}

// compile-time interface verification
var _ ports.MatchRepository = (*MemoryMatchAdapter)(nil)

func NewMemoryMatchAdapter() *MemoryMatchAdapter {
	return &MemoryMatchAdapter{
		matches: make(map[string]domain.MatchState),
	}
}

// Get returns the current state for a fixture.
// Uses RLock since reads can happen concurrently with other reads.
func (m *MemoryMatchAdapter) Get(
	ctx context.Context,
	fixtureID string,
) (domain.MatchState, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	state, ok := m.matches[fixtureID]

	if !ok {
		return domain.MatchState{}, domain.ErrMatchNotFound
	}

	// state is a value copy of the map entry, so the caller can't
	// mutate our internal state through the returned struct.
	return state, nil
}

// Save writes the full state for a fixture, overwriting any previous value.
// Note: this is NOT atomic with the Get a caller may have done beforehand -
// a Get + Apply + Save sequence (see EventProcessor) can race if two events
// for the same fixture are processed concurrently. Fixing that requires
// either a per-fixture lock held across the whole sequence, or a
// compare-and-swap on Version here.
func (m *MemoryMatchAdapter) Save(
	ctx context.Context,
	state domain.MatchState,
) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.matches[state.FixtureID] = state

	return nil
}
