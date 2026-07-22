package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/istvzsig/matchpulse/internal/application"
	"github.com/istvzsig/matchpulse/internal/domain"
)

// Handler wires the application services to HTTP endpoints.
type Handler struct {
	matchService   *application.MatchService
	eventProcessor *application.EventProcessor
	fixtureService *application.FixtureService
}

func NewHandler(
	matchService *application.MatchService,
	eventProcessor *application.EventProcessor,
	fixtureService *application.FixtureService,
) *Handler {

	return &Handler{
		matchService:   matchService,
		eventProcessor: eventProcessor,
		fixtureService: fixtureService,
	}
}

// RegisterRoutes attaches all handlers to the given mux.
// Uses Go 1.22+ method+path patterns (e.g. "POST /matches").
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /matches", h.CreateMatch)
	mux.HandleFunc("GET /matches/{fixtureID}", h.GetMatch)
	mux.HandleFunc("POST /matches/{fixtureID}/events", h.SubmitEvent)
	mux.HandleFunc("GET /fixtures", h.GetFixtures)
}

func (h *Handler) CreateMatch(w http.ResponseWriter, r *http.Request) {
	var req createMatchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.FixtureID == "" {
		writeError(w, http.StatusBadRequest, "fixture_id is required")
		return
	}

	view, err := h.matchService.CreateMatch(r.Context(), req.FixtureID)

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, view)
}

func (h *Handler) GetMatch(w http.ResponseWriter, r *http.Request) {
	fixtureID := r.PathValue("fixtureID")

	view, err := h.matchService.GetMatch(r.Context(), fixtureID)

	if err != nil {
		if errors.Is(err, domain.ErrMatchNotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, view)
}

func (h *Handler) SubmitEvent(w http.ResponseWriter, r *http.Request) {
	fixtureID := r.PathValue("fixtureID")

	var req submitEventRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	event := domain.MatchEvent{
		ID:        req.ID,
		FixtureID: fixtureID,
		Type:      domain.EventType(req.Type),
		TeamID:    req.TeamID,
		Value:     req.Value,
		Status:    domain.MatchStatus(req.Status),
		Winner:    req.Winner,
		Timestamp: time.Now(),
		Version:   req.Version,
	}

	if err := h.eventProcessor.Process(r.Context(), event); err != nil {
		switch {
		case errors.Is(err, domain.ErrMatchNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, domain.ErrOutdatedEvent),
			errors.Is(err, domain.ErrMatchFinished),
			errors.Is(err, domain.ErrInvalidWinner),
			errors.Is(err, domain.ErrInvalidEvent):
			// client sent something the current state can't accept
			writeError(w, http.StatusConflict, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	view, err := h.matchService.GetMatch(r.Context(), fixtureID)

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, view)
}

func (h *Handler) GetFixtures(w http.ResponseWriter, r *http.Request) {
	fixtures, err := h.fixtureService.GetFixtures(r.Context())

	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, fixtures)
}
