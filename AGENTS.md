# AGENTS.md

Guidance for AI coding agents working in this repository. Humans should read `README.md` instead.

## Project

Octagon is a tournament management toolkit for **Super Smash Bros. Ultimate** tournaments, built for *The Octagon*, a weekly event in Seattle. It automates the repetitive parts of running a bracket: seeding players by rating, detecting and resolving conflicts, reporting sets, updating stream overlays, and more.

The project values:

- **UX.** Commands and TUIs save the TO time under pressure. Sensible defaults, clear errors, no surprises.
- **High test coverage.** New features ship with tests. Bug fixes ship with a regression test.
- **Clean, maintainable code.** Small packages with clear boundaries. Idiomatic Go.

## Domain glossary

Use these terms precisely. Agents that conflate them produce broken UX.

- **TO** — Tournament Organizer. The primary user of this tool.
- **SSBU** — Super Smash Bros. Ultimate (the game).
- **Set** — a match between two entrants, typically best-of-3 or best-of-5 games.
- **Game** — a single round inside a set.
- **Event** — a competition within a tournament (e.g. Ultimate Singles, Redemption Bracket).
- **Bracket** — the tournament structure; Octagon uses double-elimination by default.
- **Entrant** — a registered participant (one player in singles, a team in doubles).
- **Seed / seeding** — initial bracket placement based on player strength. Good seeding prevents top players from meeting early.
- **Conflict** — two players who should not be placed near each other in seeding (e.g. same crew, same region, personal history).
- **Pool** — a round-robin or mini-bracket group in the first phase of a multi-phase tournament.
- **DQ** — disqualification, usually for no-show.
- **Redemption bracket** — a secondary bracket for players eliminated early.
- **PR** — Power Ranking. **Not** "pull request" in this codebase. When an agent output could be ambiguous, prefer "power ranking" or "pull request" explicitly.
- **start.gg** — primary tournament hosting backend (GraphQL API, dev endpoint).
- **parry.gg** — secondary tournament hosting backend (gRPC API). Integration is in progress; not all features are at parity with start.gg.
- **SchuStats** — external Elo-style rating service for SSBU players, read via Firebase.

## Repository map

```
AGENTS.md                    this file
README.md                    human-facing docs
main.go                      CLI entry point (binary name: octagon)
cmd/                         CLI subcommands (urfave/cli/v3)
  root.go                    registers all subcommands
  attendees/                 list attendees; `add` is broken and logs a warning
  bracket/                   print, seed, and list seeds
  cache/                     inspect, populate, clear local cache
  conflicts/                 create/list manual conflicts and set requests
  rating/                    player rating lookups + per-player bias multipliers
  report/                    interactive set reporting (bubbletea TUI)
  stream/                    stream overlay automation; `select` is a bubbletea TUI
startgg/                     start.gg GraphQL client (genqlient) — dev API, in use
  prod/                      UNUSED scaffolding for web-prod GraphQL schema
parrygg/                     parry.gg gRPC client
  pb/                        protobuf-generated code (do not hand-edit)
obs/                         OBS Studio WebSocket client (goobs)
ocr/                         image-to-text via Tesseract (gosseract)
brackets/                    bracket data models and round logic (pure, stdlib only)
conflicts/                   conflict-checking algorithm, parallel Monte Carlo scorer
seeding/                     seeding algorithm (ratings + conflicts)
stream/                      OCR loop + OBS/parry.gg stream state updater
tournament/                  set fetching + reporting wrapper around startgg
ratings/                     SchuStats ratings via Firebase Realtime DB
cache/                       Badger v4 key-value cache (singleton global db)
characters/                  SSBU character metadata + fuzzy name lookup
config/                      env loading + aliases.json + bias.json
api/                         HTTP API server (gin) — consumed by web/
web/                         React 19 + TypeScript frontend (see web/AGENTS.md)
tools/                       standalone Python utilities (unrelated to go build)
```

## Setup

Required Go version: **1.24.1** (see `go.mod`).

### Configuration

Config is loaded in this order (first hit wins):

1. `~/.config/octagon/octagonrc`
2. `./.env` (in the current working directory)

If neither exists, the command fatals.

Required and optional environment variables:

```dotenv
# start.gg — required for any tournament operation
STARTGG_API_KEY=

# parry.gg — required for `octagon stream` subcommands
PARRYGG_API_KEY=

# Firebase (SchuStats ratings) — required for rating/seeding flows
FIREBASE_API_KEY=
FIREBASE_DATABASE_URL=

# OBS WebSocket — required for stream overlay flows
OBS_WS_PORT=4455
OBS_WS_PASS=

# Optional
DEBUG=1    # debug-level logging from charmbracelet/log
PORT=8080  # api/ HTTP server port
```

Additional Firebase vars (`FIREBASE_AUTH_DOMAIN`, `FIREBASE_PROJECT_ID`, `FIREBASE_STORAGE_BUCKET`, `FIREBASE_MESSAGING_SENDER_ID`, `FIREBASE_APP_ID`, `FIREBASE_MEASUREMENT_ID`) may appear in `.env` for the web frontend but are not read by Go code.

Never commit `.env` or `octagonrc`. `.env` is gitignored.

## Build, run, test

```bash
# Build the CLI binary
go build

# Install to $GOPATH/bin so `octagon` runs from anywhere
go build && go install

# Run without installing
go run . <subcommand> [flags]

# Run with debug logs
DEBUG=1 go run . <subcommand>

# Run all tests
go test ./...

# Run tests for a single package
go test ./conflicts/...

# Run a single test by name
go test ./conflicts/ -run TestScorer

# Start the dev HTTP API (feeds web/)
./start-dev.sh
```

Always run `go test ./...` before declaring a task complete.

## Code style

- Format with `gofmt`/`goimports`. No exceptions.
- Package comments on every package. Existing pattern: `// Package foo does X` or `/* Package foo does X */`.
- Error wrapping: `fmt.Errorf("context: %w", err)`. Always include actionable context.
- Logging: **`github.com/charmbracelet/log`** exclusively. Use structured fields: `log.Debug("fetching", "url", u)`. Do not introduce `fmt.Println` or the standard `log` package. (`api/main.go` uses `fmt.Printf` in one spot; treat that as a bug to fix, not a pattern to copy.)
- External clients (start.gg, parry.gg, OBS, Firebase, OCR) use a `sync.Once` singleton with a `GetClient(ctx)` or `Client()` constructor, reading credentials from env at init. Follow this pattern for new backends.
- Prefer small files focused on one responsibility. `conflicts/` is the reference for splitting: models, checker, scorer, io, cache, printer.
- There is **no service-layer abstraction** between backends and consumers. Both `cmd/*` and `api/` import `startgg`, `parrygg`, `tournament`, `seeding`, `conflicts` directly. `tournament/` is the closest thing to a service wrapper; follow it when adding new cross-cutting flows.

## Testing conventions

- Table-driven tests are the norm. See `conflicts/scorer_test.go` and `brackets/bracket_test.go` for the pattern.
- Tests live next to the code they test: `foo.go` paired with `foo_test.go`.
- External API calls (start.gg, parry.gg, Firebase, OBS) are not mocked with a framework. Isolate them behind small functions so the surrounding logic can be tested without a network.
- `web/` uses React Testing Library + Jest via `react-scripts test`.
- New feature → add tests. Bug fix → add a regression test that fails without your fix.

## UX conventions for CLI and TUI work

- Subcommand names are verbs or short nouns: `attendees`, `report`, `bracket`, `stream`.
- Every subcommand provides `-h` / `--help` output via urfave/cli. Keep `Usage` strings short and action-oriented.
- Flags have both long and short forms where reasonable; short forms must not collide within a command.
- Defaults target *The Octagon* specifically: `--tournament=octagon`, `--event=<Ultimate Singles>`. Overridable.
- Destructive operations (cache clear, bracket reset) confirm before acting or require an explicit flag.
- TUIs use **bubbletea + bubbles + lipgloss**. Existing TUIs: `octagon report` and `octagon stream select`. Keybindings: `q`/`Ctrl+C` quits, `Esc` goes back, `Enter` submits/selects, `Tab` moves between fields. `cmd/report/model.go` is the reference state-machine implementation.
- Errors shown to the user are phrased for the TO, not the developer. Surface what to do next, not just what broke.

## Caching

- Backed by **Badger v4** at `/tmp/octagon-cache`. Ephemeral — wiped on system reboot.
- Default TTL is 90 days. Use `cache.Set` / `cache.Get` from `cache/`.
- Cache anything slow and deterministic (tournament metadata, player data). Do not cache anything with per-request auth state or PII beyond public tournament data.
- `octagon cache clear` drops the whole cache.

## External backends

### start.gg (primary)

- GraphQL via **Khan/genqlient**. Dev API at `https://api.start.gg/gql/alpha`, Bearer auth via `STARTGG_API_KEY`.
- `startgg/prod/` is **unused scaffolding** generated against the web-prod schema (`https://www.start.gg/api/-/gql`). Zero Go files import it. Do not build features on it without first wiring it up.
- `startgg.ID` is `type ID any` (binds GraphQL `ID` to `interface{}`). Consumers type-assert `.(float64)`. A non-float64 ID silently drops the record. When adding code that handles IDs, match the existing cast pattern in `tournament/tournament.go`.
- Regenerating the client after a schema change is documented in `startgg/AGENTS.md`.

### parry.gg (secondary, in progress)

- gRPC + TLS to `api.parry.gg:443`. Client is a struct exposing all service stubs (`TournamentService`, `EventService`, etc.). Auth via `X-API-KEY` metadata, injected by `client.WithAuth(ctx)`.
- **Security gotcha:** the current `loggingInterceptor` in `parrygg/client.go` logs full gRPC metadata at INFO level on every call, which includes the API key. Do not run with INFO logs in any shared or recorded environment until this is fixed. Fixing it is in scope for any PR touching that file.
- `parrygg/pb/` is generated from https://github.com/parry-gg/protos. Do not hand-edit. Regeneration config is not checked in; it is a manual step.
- See `parrygg/AGENTS.md` for the parity matrix with start.gg.

### OBS Studio

- WebSocket via **goobs**. See `obs/AGENTS.md` for source-naming rules and gotchas.

### Firebase (SchuStats ratings)

- Read-only lookup of player ratings. Requires `FIREBASE_API_KEY` and `FIREBASE_DATABASE_URL`.
- Used by `ratings/` and consumed by `seeding/`.

## Commit and PR conventions

Conventional commits only. Types used in this repo:

```
feat:     new feature
fix:      bug fix
perf:     performance improvement
refactor: code change that neither fixes a bug nor adds a feature
docs:     documentation only
test:     test-only change
chore:    tooling, deps, non-code
style:    formatting, whitespace, no logic change
```

Examples:

```
feat: 'conflict remove' subcommand
fix: increase pageSize in getParticipants to prevent missing players
perf: cached getTournament responses
```

- Keep PRs small and scoped to one concern.
- Update the relevant `AGENTS.md` in the same PR as behavior changes that affect agent context.
- Do not force-push to `main`.

## Gotchas

### start.gg

- Only `startgg.GetSeedsPaginated` actually paginates. `GetParticipants` (perPage 100), `GetReportableSets` (perPage 500), and `GetSets` (perPage 100) are **single-page** and silently miss data beyond the first page. `tournament.FetchReportableSets` warns when exactly 500 sets are returned — treat that warning as a bug that needs fixing, not a normal state.
- `startgg.ID` is `interface{}`. Assume every ID cast through `.(float64)` can drop the record.
- `cmd/attendees add` does not work (start.gg API does not permit it) and is explicitly stubbed with a warning.

### parry.gg

- Parity with start.gg is partial — see `parrygg/AGENTS.md`.
- The logging interceptor currently leaks the API key at INFO level (see External backends above).

### Config and cache

- `config.Load()` fatals if neither `~/.config/octagon/octagonrc` nor `./.env` is present.
- `api/main.go` bypasses `config.Load()` and calls `godotenv.Load()` directly. Changes to config loading must update both call sites.
- Cache is in `/tmp/octagon-cache` with a 90-day TTL. Wiped on system reboot. Do not rely on long-term persistence.
- `cache.Open()` must run before any other package calls `cache.Set`/`Get`, or they nil-deref. `main.go` and `api/main.go` do this. Tests swap the global `db` in `TestMain`.

### Workflow dependencies

- `octagon conflicts create` uses fuzzy-matching against the player cache. Run `octagon cache populate` first, or creation will fail to resolve names.
- `octagon seeding` auto-generates conflicts from the **previous** tournament's sets by parsing the current slug (`octagon-<n>` → `octagon-<n-1>`). Non-standard slugs break this.
- Conflict priority is signed: `1..3` = conflict (avoid pairing), `-1..-3` = set request (prefer pairing). Same field, sign flips semantics.

### OBS and stream

- OBS source names are case-sensitive and must match exactly as shown in the OBS Sources panel. The current code hardcodes `"VOD"` for screenshots and `"Player 1 Name"`, `"Player 2 Name"`, `"Player 1 Score"`, `"Player 2 Score"` for text fields.
- Stream pixel bounds in `stream/bounds.go` assume a 1920×1080 capture. Different resolutions break OCR silently.
- `octagon stream select` hardcodes `TournamentSlug="octagon"`, `EventSlug="ultimate-singles"`, `PhaseSlug="main"`, `BracketSlug="bracket"`. Flagify before reusing for other events.
- `stream.Command` registers with `Aliases: []string{"c", "conflict"}`, which collides with `conflicts`. Treat this as a bug to fix, not documented behavior.

### Firebase and ratings

- Firebase auth is API-key-only (`option.WithAPIKey`), no service-account JSON. `FIREBASE_API_KEY` and `FIREBASE_DATABASE_URL` are both required for any rating or seeding flow.
- Seeding continues silently with a 0 rating when Firebase fetch fails for a player. If seed quality drops mysteriously, check Firebase first.

### Terminology

- **`PR` is ambiguous.** In commit messages and comments it refers to Power Ranking, never pull requests. Use "pull request" in full when unclear.

### Dead code to avoid modeling

- `startgg/prod/` — unused scaffolding.
- `cmd/attendees/add.go` — documented broken.
- `AmazonQ.md` — references a non-existent `cmd/report/sets.go`. Slated for deletion.

## Nested AGENTS.md files

Consult these for subsystem-specific rules:

- `startgg/AGENTS.md` — start.gg GraphQL client, schema regeneration
- `parrygg/AGENTS.md` — parry.gg gRPC client, parity matrix
- `obs/AGENTS.md` — OBS WebSocket integration
- `web/AGENTS.md` — React frontend

Agents should read the nearest `AGENTS.md` walking up from the file being edited; rules nested deeper override the root.
