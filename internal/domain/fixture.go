package domain

import "time"

type Fixture struct {
	ID          string
	Tournament  string
	ScheduledAt time.Time
	Teams       []Team
	Players     []Player
}
