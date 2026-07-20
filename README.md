# MatchPulse

## Real-time esports match aggregation platform built with Go

MatchPulse is a backend service that simulates the architecture of a live esports data platform.

It receives match events, processes them, aggregates the current match state, and exposes a frontend-ready view of the game.

The project focuses on understanding the engineering challenges behind real-time systems:

- event processing
- concurrency
- state management
- scalability
- reliability
- observability

---

# 🎮 Project Story

Imagine an esports tournament with thousands of matches happening simultaneously.

Every second, the system receives events:

- Player killed
- Objective captured
- Round started
- Round finished
- Match completed

The backend must:

1. Receive events
2. Validate them
3. Process them safely
4. Update match state
5. Publish the latest view for viewers

This sounds simple.

The difficult parts appear when:

- events arrive out of order
- events are duplicated
- many updates happen simultaneously
- processing becomes slower than incoming traffic
- servers restart during live matches

MatchPulse explores how real backend systems solve these problems.

---

# 🎯 Goals

This project is built to practice production backend engineering concepts using Go.

Main learning goals:

- Clean architecture
- Event-driven design
- Concurrent programming
- Thread-safe state management
- Worker pools
- Context cancellation
- Graceful shutdown
- Caching strategies
- Testing strategies
- Observability
- Performance considerations

---

# 🏗️ Architecture

Current architecture:

```text
			HTTP API

                |
                v

         Event Processor

                |
                v

          Match Service

                |
                v

      In-memory Match State

                |
                v

         View Model API
```

The architecture will evolve as new requirements are introduced.

---

# ✅ Current Status

**Implemented:**

- Domain layer: `MatchState`, `MatchEvent`, `Fixture`, `Team`, `Player`
- Event application logic (`MatchState.Apply`) with versioning/out-of-order rejection
- In-memory `MatchRepository` adapter (`MemoryMatchAdapter`)
- In-memory `FixtureRepository` adapter (`MemoryFixtureAdapter`)
- `EventProcessor` and `FixtureService` application services
- `ViewBuilder` for read-model projection
- Unit tests for score update event application and fixture adapter behavior

**Not yet implemented:**

- HTTP API (Level 1) — `cmd/api/main.go` is a placeholder
- Concurrency safety across multi-step operations (see known issue below)
- Worker pools, graceful shutdown, observability (Levels 4–6)

**Known issue:**

- `EventProcessor.Process` performs Get → Apply → Save as separate steps. Each step is individually thread-safe, but the sequence as a whole is not — concurrent events for the same fixture can race. This is the Level 3 concurrency problem the roadmap below calls out.

---

# 🚀 Development Roadmap

The project grows incrementally.

Each stage introduces a real engineering problem.

## Level 1 - Match Server

Goal:

Create a simple HTTP service.

Features:

- Create matches
- Receive events
- Query match state

Concepts:

- HTTP handlers
- services
- models
- testing

---

## Level 2 - Event Processing

Goal:

Process live match events.

Example:

- Team A scored
- Team B scored
- Round ended

Concepts:

- event modeling
- aggregation
- state transitions

---

## Level 3 - Concurrency

Problem:

Multiple events arrive at the same time.

Concepts:

- goroutines
- mutexes
- race conditions
- Go race detector

---

## Level 4 - Worker Processing

Problem:

Events arrive faster than they can be processed.

Concepts:

- channels
- worker pools
- backpressure

---

## Level 5 - Reliability

Problem:

The server restarts during a live game.

Concepts:

- graceful shutdown
- context cancellation
- recovery strategies

---

## Level 6 - Production Features

Concepts:

- logging
- metrics
- tracing
- caching
- rate limiting

---

# 🧠 Engineering Decisions

Important design choices are documented using ADRs.

Examples:

```text
docs/adr/

ADR-0001-use-in-memory-state.md

ADR-0002-event-processing-model.md

ADR-0003-worker-pool-design.md
```

---

# 🧪 Testing

The project includes:

- unit tests
- integration tests
- concurrency tests

Concurrency checks:

```bash
go test --race ./...
```

---

# 🛠️ Technology

Current:

- Go
- Standard library
- net/http
- sync package
- testing package

Additional technologies will only be introduced when they solve a real problem.

---

# 📚 Learning Notes

This repository is also a learning resource.

Each feature explains:

- the problem
- possible solutions
- tradeoffs
- implementation
- testing approach

The goal is not just to build a system.

The goal is to understand why systems are designed the way they are.

---

# 🎮 Connection to Real Systems

Although MatchPulse is a game-inspired project, the same engineering patterns appear in:

- esports platforms
- payment systems
- IoT platforms
- analytics pipelines
- messaging systems
- real-time applications

---

# Future Ideas

Possible future improvements:

- persistent storage
- distributed event processing
- message brokers
- multiple service instances
- authentication
- Kubernetes deployment

---

# Author

Built while learning and practicing backend engineering with Go.
