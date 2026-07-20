package memory

import (
	"context"
	"sync"

	"github.com/istvzsig/matchpulse/internal/domain"
	"github.com/istvzsig/matchpulse/internal/ports"
)

// MemoryFixtureAdapter is an in-memory implementation of ports.FixtureRepository.
// Like MemoryMatchAdapter, it lets the application layer be developed and
// tested without a real database — a persistent adapter can implement the
// same port later without changes above this layer.
type MemoryFixtureAdapter struct {
	fixtures []domain.Fixture
	mu       sync.RWMutex // guards concurrent access to fixtures
}

// compile-time interface verification
var _ ports.FixtureRepository = (*MemoryFixtureAdapter)(nil)

func NewMemoryFixtureAdapter() *MemoryFixtureAdapter {
	return &MemoryFixtureAdapter{
		fixtures: make([]domain.Fixture, 0),
	}
}

// Seed preloads fixtures into the adapter. This exists because
// FixtureRepository only defines a read (GetFixtures) — there's no write
// method on the port itself, so tests and bootstrap code use this instead.
// If fixtures ever need to be created dynamically (e.g. via an API), that's
// a sign the port itself should grow a Save/Create method.
func (m *MemoryFixtureAdapter) Seed(fixtures []domain.Fixture) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.fixtures = append(m.fixtures, fixtures...)
}

// GetFixtures returns all known fixtures.
// Uses RLock since reads can happen concurrently with other reads.
func (m *MemoryFixtureAdapter) GetFixtures(
	ctx context.Context,
) ([]domain.Fixture, error) {

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy so the caller can't mutate our internal slice
	// through the returned value.
	out := make([]domain.Fixture, len(m.fixtures))
	copy(out, m.fixtures)

	return out, nil
}
