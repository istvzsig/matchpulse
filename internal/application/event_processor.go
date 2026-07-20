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

func (processor *EventProcessor) Process(
	ctx context.Context,
	event domain.MatchEvent,
) error {

	state, err := processor.repository.Get(
		ctx,
		event.FixtureID,
	)

	if err != nil {
		return err
	}

	err = state.Apply(event)

	if err != nil {
		return err
	}

	return processor.repository.Save(
		ctx,
		state,
	)
}
