package application

import (
	"context"

	"github.com/istvzsig/matchpulse/internal/domain"
	"github.com/istvzsig/matchpulse/internal/ports"
)

// EventProcessor coordinates incoming match events.
// Business rules are handled by domain.MatchState.
type EventProcessor struct {
	repository ports.MatchRepository
}

func NewEventProcessor(
	repository ports.MatchRepository,
) *EventProcessor {

	return &EventProcessor{
		repository: repository,
	}
}

// Process applies a single event to the fixture's match state.
// Uses Update rather than separate Get/Save calls so the whole
// read-apply-write sequence is atomic per fixture, even under
// concurrent event submission.
func (processor *EventProcessor) Process(
	ctx context.Context,
	event domain.MatchEvent,
) error {

	return processor.repository.Update(
		ctx,
		event.FixtureID,
		func(state domain.MatchState) (domain.MatchState, error) {
			err := state.Apply(event)
			return state, err
		},
	)
}
