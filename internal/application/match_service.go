package application

import (
	"context"

	"github.com/istvzsig/matchpulse/internal/domain"
	"github.com/istvzsig/matchpulse/internal/ports"
)

// MatchService handles match lifecycle operations that sit outside
// event processing: creation and read-model queries.
type MatchService struct {
	repository  ports.MatchRepository
	viewBuilder *ViewBuilder
}

func NewMatchService(
	repository ports.MatchRepository,
	viewBuilder *ViewBuilder,
) *MatchService {

	return &MatchService{
		repository:  repository,
		viewBuilder: viewBuilder,
	}
}

// CreateMatch initializes a new match state for a fixture and persists it.
func (s *MatchService) CreateMatch(
	ctx context.Context,
	fixtureID string,
) (MatchView, error) {

	state := domain.NewMatchState(fixtureID)

	if err := s.repository.Save(ctx, state); err != nil {
		return MatchView{}, err
	}

	return s.viewBuilder.Build(state), nil
}

// GetMatch fetches current state and projects it into a MatchView.
func (s *MatchService) GetMatch(
	ctx context.Context,
	fixtureID string,
) (MatchView, error) {

	state, err := s.repository.Get(ctx, fixtureID)

	if err != nil {
		return MatchView{}, err
	}

	return s.viewBuilder.Build(state), nil
}
