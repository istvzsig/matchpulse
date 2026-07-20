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
}
