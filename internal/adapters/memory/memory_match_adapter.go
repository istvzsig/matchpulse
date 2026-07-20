package memory

import (
	"context"
	"sync"

	"github.com/istvzsig/matchpulse/internal/domain"
	"github.com/istvzsig/matchpulse/internal/ports"
)

type MemoryMatchAdapter struct {
	matches map[string]domain.MatchState
	mu      sync.RWMutex
}

// compile-time interface verification
var _ ports.MatchRepository = (*MemoryMatchAdapter)(nil)

func NewMemoryMatchAdapter() *MemoryMatchAdapter {
	return &MemoryMatchAdapter{
		matches: make(map[string]domain.MatchState),
	}
}

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

	return state, nil
}

func (m *MemoryMatchAdapter) Save(
	ctx context.Context,
	state domain.MatchState,
) error {

	m.mu.Lock()
	defer m.mu.Unlock()

	m.matches[state.FixtureID] = state

	return nil
}
