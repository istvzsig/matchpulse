package application

import (
	"context"

	"github.com/istvzsig/matchpulse/internal/domain"
	"github.com/istvzsig/matchpulse/internal/ports"
)

type FixtureService struct {
	repository ports.FixtureRepository
}

func NewFixtureService(
	repository ports.FixtureRepository,
) *FixtureService {

	return &FixtureService{
		repository: repository,
	}
}

func (s *FixtureService) GetFixtures(
	ctx context.Context,
) ([]domain.Fixture, error) {

	return s.repository.GetFixtures(ctx)
}
