package domain

import "errors"

var (
	ErrInvalidEvent  = errors.New("invalid event")
	ErrOutdatedEvent = errors.New("outdated event")
	ErrMatchFinished = errors.New("match already finished")
	ErrInvalidWinner = errors.New("winner can only be assigned to finished match")
)
