package application

import "github.com/istvzsig/matchpulse/internal/domain"

type MatchView struct {
	MatchID string         `json:"match_id"`
	Status  string         `json:"status"`
	Scores  map[string]int `json:"scores"`
	Winner  string         `json:"winner"`
}

type ViewBuilder struct {
}

func NewViewBuilder() *ViewBuilder {
	return &ViewBuilder{}
}

func (b *ViewBuilder) Build(
	state domain.MatchState,
) MatchView {

	return MatchView{
		MatchID: state.FixtureID,
		Status:  string(state.Status),
		Scores:  state.Scores,
		Winner:  state.Winner,
	}
}
