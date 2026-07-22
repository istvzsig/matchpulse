package httpapi

// createMatchRequest is the request body for POST /matches.
type createMatchRequest struct {
	FixtureID string `json:"fixture_id"`
}

// submitEventRequest is the request body for POST /matches/{fixtureID}/events.
// FixtureID itself comes from the URL path, not the body.
type submitEventRequest struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	TeamID  string `json:"team_id,omitempty"`
	Value   int    `json:"value,omitempty"`
	Status  string `json:"status,omitempty"`
	Winner  string `json:"winner,omitempty"`
	Version int64  `json:"version"`
}

type errorResponse struct {
	Error string `json:"error"`
}
