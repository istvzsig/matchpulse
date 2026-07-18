package ports

import (
	"context"

	"github.com/istvzsig/matchpulse/internal/domain"
)

type FixtureRepository interface {
	GetFixtures(ctx context.Context) ([]domain.Fixture, error)
}
