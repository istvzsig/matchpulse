package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/istvzsig/matchpulse/internal/adapters/memory"
	"github.com/istvzsig/matchpulse/internal/application"
)

func newTestHandler() *Handler {
	matchRepo := memory.NewMemoryMatchAdapter()
	viewBuilder := application.NewViewBuilder()
	matchService := application.NewMatchService(matchRepo, viewBuilder)
	eventProcessor := application.NewEventProcessor(matchRepo)
	fixtureRepo := memory.NewMemoryFixtureAdapter()
	fixtureService := application.NewFixtureService(fixtureRepo)

	return NewHandler(matchService, eventProcessor, fixtureService)
}

func TestCreateAndGetMatch(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Create
	createReq := httptest.NewRequest(
		http.MethodPost,
		"/matches",
		strings.NewReader(`{"fixture_id":"F1"}`),
	)
	createRec := httptest.NewRecorder()
	mux.ServeHTTP(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", createRec.Code, createRec.Body.String())
	}

	// Get
	getReq := httptest.NewRequest(http.MethodGet, "/matches/F1", nil)
	getRec := httptest.NewRecorder()
	mux.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", getRec.Code, getRec.Body.String())
	}

	var view application.MatchView
	if err := json.NewDecoder(getRec.Body).Decode(&view); err != nil {
		t.Fatal(err)
	}

	if view.MatchID != "F1" {
		t.Fatalf("expected match_id F1, got %s", view.MatchID)
	}

	if view.Status != "scheduled" {
		t.Fatalf("expected status scheduled, got %s", view.Status)
	}
}

func TestGetMatchNotFound(t *testing.T) {
	h := newTestHandler()
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/matches/unknown", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}
