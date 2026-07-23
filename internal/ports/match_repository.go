package ports

import (
	"context"

	"github.com/istvzsig/matchpulse/internal/domain"
)

type MatchRepository interface {
	Get(
		ctx context.Context,
		fixtureID string,
	) (domain.MatchState, error)

	Save(
		ctx context.Context,
		state domain.MatchState,
	) error

	// Update performs an atomic read-modify-write: it fetches the current
	// state for fixtureID, passes it to mutate, and persists whatever
	// mutate returns - all under a single lock, so no other Update/Save
	// for the same fixture can interleave. mutate should return an error
	// (without a valid state) if the mutation itself is invalid; in that
	// case nothing is persisted.
	Update(
		ctx context.Context,
		fixtureID string,
		mutate func(state domain.MatchState) (domain.MatchState, error),
	) error
}
