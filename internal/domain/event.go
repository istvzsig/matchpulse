package domain

import "time"

type EventType string

const (
	EventScoreUpdate  EventType = "score_update"
	EventStatusChange EventType = "status_change"
	EventWinnerUpdate EventType = "winner_update"
)

type MatchEvent struct {
	ID        string
	FixtureID string
	Type      EventType
	TeamID    string
	Value     int
	Winner    string
	Timestamp time.Time
	Version   int64
}
