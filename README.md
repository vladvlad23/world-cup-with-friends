# World Cup 2026 Meetup Planner

Plan social meetups around the 2026 FIFA World Cup matches. Log in, browse the
match calendar, and for any match propose a meetup, see existing ones, and join.

- **Frontend:** SvelteKit (Node adapter)
- **Backend:** Go (chi + pgx)
- **Database:** PostgreSQL
- **Match data:** the 104 fixtures (teams, dates, stadiums) are seeded into
  Postgres from static data originally published by
  [rezarahiminia/worldcup2026](https://github.com/rezarahiminia/worldcup2026).
  The app has **no runtime dependency** on that API.

## Quick start (Docker)

```bash
cp .env.example .env      # optionally edit secrets/ports
docker compose up --build
```

Then open **http://localhost:3000**.

Demo users (all password `password`): `alice`, `bob`, `carol`, `dave`, `erin`.
Edit `backend/internal/seed/data/users.json` to change them (re-seeded on every
backend start).

## What you can do

1. **Log in** with a hardcoded user.
2. **Calendar** — every match across June/July 2026, each with a badge showing
   how many meetups exist.
3. Click a match → **Propose a meetup** / **See existing meetups** / **Cancel**.
4. **Propose a meetup**: location name (with a bar-reservation warning if it's a
   bar), optional Google Maps link, invite other users, and an `@`-mention note.
5. **See existing meetups** → open one → **Join** or go back.

## Architecture

```
Browser ──> SvelteKit (Node, :3000) ──> Go API (:8080) ──> Postgres
```

The JWT lives in an httpOnly cookie set by SvelteKit; the browser never calls
the Go API directly (no CORS, token never exposed to client JS). SvelteKit
`load` functions and `/api/*` proxy endpoints forward requests to Go with the
token attached.

## Local development (without Docker)

```bash
# 1. Postgres
docker run -d --name wc-db -e POSTGRES_USER=worldcup -e POSTGRES_PASSWORD=worldcup \
  -e POSTGRES_DB=worldcup -p 5432:5432 postgres:16-alpine

# 2. Backend (seeds on first run)
cd backend && go run ./cmd/server      # listens on :8080

# 3. Frontend
cd frontend && npm install && BACKEND_URL=http://localhost:8080 npm run dev
```

## API endpoints (Go backend)

| Method | Path | Description |
| ------ | ---- | ----------- |
| POST | `/api/auth/login` | Login, returns JWT |
| GET | `/api/me` | Current user |
| GET | `/api/users` | All app users (for invites/mentions) |
| GET | `/api/matches` | All matches + meetup counts |
| GET | `/api/matches/{id}` | Match + its meetups |
| POST | `/api/matches/{id}/meetups` | Create a meetup |
| GET | `/api/meetups/{id}` | Meetup detail |
| POST | `/api/meetups/{id}/join` | Join a meetup |

## Roadmap

- **Betting** (planned). Would add live scores via the upstream API and a
  wagering model on top of the existing meetup/match schema.
