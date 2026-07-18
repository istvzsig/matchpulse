package domain

import (
	"strconv"
	"time"
)

type MatchStatus string

const (
	StatusScheduled MatchStatus = "scheduled"
	StatusLive      MatchStatus = "live"
	StatusFinished  MatchStatus = "finished"
)

type MatchState struct {
	FixtureID string
	Status    MatchStatus
	Scores    map[string]int
	Winner    string
	Version   int64
	UpdatedAt time.Time
}

func NewMatchState(fixtureID string) MatchState {
	return MatchState{
		FixtureID: fixtureID,
		Status:    StatusScheduled,
		Scores:    make(map[string]int),
		Version:   0,
		UpdatedAt: time.Now(),
	}
}

func (m *MatchState) Apply(event MatchEvent) error {

	// Reject old events
	if event.Version <= m.Version {
		return ErrOutdatedEvent
	}

	switch event.Type {

	case EventScoreUpdate:
		if m.Status == StatusFinished {
			return ErrMatchFinished
		}

		m.Scores[event.TeamID] += event.Value

	case EventStatusChange:
		m.Status = MatchStatus(strconv.Itoa(event.Value))

	case EventWinnerUpdate:
		if m.Status != StatusFinished {
			return ErrInvalidWinner
		}

		m.Winner = event.Winner

	default:
		return ErrInvalidEvent
	}

	m.Version = event.Version
	m.UpdatedAt = event.Timestamp

	return nil
}
