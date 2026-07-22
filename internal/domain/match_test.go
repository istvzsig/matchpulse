package domain

import (
	"testing"
	"time"
)

func TestMatchStateApplyScoreUpdate(t *testing.T) {

	state := NewMatchState("F123")

	event := MatchEvent{
		ID:        "E1",
		FixtureID: "F123",
		Type:      EventScoreUpdate,
		TeamID:    "T1",
		Value:     2,
		Version:   1,
		Timestamp: time.Now(),
	}

	err := state.Apply(event)

	if err != nil {
		t.Fatal(err)
	}

	if state.Scores["T1"] != 2 {
		t.Fatalf(
			"expected score 2 got %d",
			state.Scores["T1"],
		)
	}

}

func TestMatchStateApplyStatusChange(t *testing.T) {

	state := NewMatchState("F123")

	event := MatchEvent{
		ID:        "E1",
		FixtureID: "F123",
		Type:      EventStatusChange,
		Status:    StatusLive,
		Version:   1,
		Timestamp: time.Now(),
	}

	err := state.Apply(event)

	if err != nil {
		t.Fatal(err)
	}

	if state.Status != StatusLive {
		t.Fatalf(
			"expected status %s got %s",
			StatusLive, state.Status,
		)
	}
}
