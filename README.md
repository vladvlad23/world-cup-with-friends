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
These embedded demo users are for **local dev / CI only** — see
[Production deployment](#production-deployment) for real users.

## Production deployment

Real credentials must **never** live in git. The seeder reads its user list from
an uncommitted file pointed to by `USERS_FILE`; the committed demo users are only
a dev/CI fallback. Deploy with the production overlay
(`docker-compose.prod.yml`), which mounts `./users.json` into the backend and
sets `USERS_FILE=/config/users.json`.

### Step by step

1. **Get the code onto the server** (clone the repo, or pull on the box).

2. **Configure the environment.** Copy the template and fill in the values:
   ```bash
   cp .env.production.example .env
   ```
   - `JWT_SECRET` — generate with `openssl rand -hex 32`.
   - `PUBLIC_ORIGIN` — the **exact** URL users reach the app on, including scheme
     and any non-standard port (e.g. `https://meetups.example.com`, or
     `http://1.2.3.4:3000`). If this is wrong the app loads but **every login
     fails with HTTP 403** (SvelteKit CSRF check).
   - `POSTGRES_PASSWORD` — a strong DB password.

3. **Create the user list.** Copy the template and edit with real users and
   strong passwords:
   ```bash
   cp users.example.json users.json
   ```
   `users.json` is gitignored — it stays on the server and is never pushed. Each
   object needs `username`, `display_name`, `password`, `color`.

4. **Start it:**
   ```bash
   docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
   ```

5. **Verify the right accounts loaded.** The backend log should say
   `seeding N users from USERS_FILE=/config/users.json` (not the embedded-demo
   line):
   ```bash
   docker compose logs backend | grep seeding
   ```

Then open `PUBLIC_ORIGIN` in a browser and log in.

### Operating notes

- **Passwords are self-service.** Each user logs in with the password you set and
  can change it from the header (**Change password**). Once changed, that
  password is preserved — the seeder does **not** overwrite it on restart.
- **Adding users later:** edit `users.json` and restart the backend
  (`docker compose ... up -d`). Users are upserted by username; existing users'
  changed passwords are kept, new users are seeded.
- **Start prod clean.** On a fresh database only the `users.json` accounts are
  seeded. Don't first-boot the production DB with the demo users and then switch;
  start it with `users.json` in place from the beginning.
- **Lost `users.json`?** Passwords are bcrypt-hashed in the DB and can't be
  recovered from it. Keep a copy somewhere safe, or regenerate the file and
  restart (users can also reset via the app). Back the file up in a password
  manager.
- **TLS:** put a reverse proxy (nginx / Caddy / Traefik) in front terminating
  HTTPS and forwarding to the frontend port, and set `PUBLIC_ORIGIN` to the
  `https://` URL. You can drop the backend `ports:` block entirely — the frontend
  reaches it over the internal Docker network.

## What you can do

1. **Log in** with a hardcoded user.
2. **Calendar** — every match across June/July 2026, each with a badge showing
   how many meetups exist.
3. Click a match → **Propose a meetup** / **See existing meetups** / **Cancel**.
4. **Propose a meetup**: location name (with a bar-reservation warning if it's a
   bar), optional Google Maps link, invite other users, and an `@`-mention note.
5. **See existing meetups** → open one → **Join** or go back.
6. **Change your password** from the header at any time.

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
| POST | `/api/me/password` | Change own password |
| GET | `/api/users` | All app users (for invites/mentions) |
| GET | `/api/matches` | All matches + meetup counts |
| GET | `/api/matches/{id}` | Match + its meetups |
| POST | `/api/matches/{id}/meetups` | Create a meetup |
| GET | `/api/meetups/{id}` | Meetup detail |
| POST | `/api/meetups/{id}/join` | Join a meetup |

## Tests & CI

```bash
# Backend (needs a Postgres; integration tests skip without DATABASE_URL)
cd backend && DATABASE_URL=postgres://worldcup:worldcup@localhost:5432/worldcup?sslmode=disable go test ./...

# Frontend
cd frontend && npm test
```

GitHub Actions (`.github/workflows/ci.yml`) runs on every push/PR to `main`:
a Go job (vet, build, `go test -race` against a Postgres service) and a
frontend job (`npm ci`, test, build).

## Roadmap

- **Betting** (planned). Would add live scores via the upstream API and a
  wagering model on top of the existing meetup/match schema.
