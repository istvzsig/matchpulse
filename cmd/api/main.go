package main

import (
	"log"
	"net/http"

	"github.com/istvzsig/matchpulse/internal/adapters/httpapi"
	"github.com/istvzsig/matchpulse/internal/adapters/memory"
	"github.com/istvzsig/matchpulse/internal/application"
)

func main() {
	matchRepo := memory.NewMemoryMatchAdapter()
	fixtureRepo := memory.NewMemoryFixtureAdapter()

	viewBuilder := application.NewViewBuilder()
	matchService := application.NewMatchService(matchRepo, viewBuilder)
	eventProcessor := application.NewEventProcessor(matchRepo)
	fixtureService := application.NewFixtureService(fixtureRepo)

	handler := httpapi.NewHandler(matchService, eventProcessor, fixtureService)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	log.Println("matchpulse listening on :8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
