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

// Update performs an atomic read-modify-write for a single fixture.
// The whole operation - read, mutate, write - happens under one write
// lock, so two concurrent Update calls for the same fixture can never
// interleave and silently lose one of their changes.
//
// Note: this holds the adapter's single global lock for the duration of
// the mutation, which means Update for fixture A briefly blocks Get/Save/
// Update for fixture B too. That's the simplest correct option for an
// in-memory map; if match volume grows large enough for that to matter,
// the next step would be per-fixture locks (e.g. a map[string]*sync.Mutex)
// instead of one lock for the whole adapter.
func (m *MemoryMatchAdapter) Update(
	ctx context.Context,
	fixtureID string,
	mutate func(state domain.MatchState) (domain.MatchState, error),
) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.matches[fixtureID]

	if !ok {
		return domain.ErrMatchNotFound
	}

	newState, err := mutate(state)

	if err != nil {
		return err
	}

	m.matches[fixtureID] = newState

	return nil
}
