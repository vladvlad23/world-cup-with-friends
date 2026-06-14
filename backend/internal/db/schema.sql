CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    username      TEXT NOT NULL UNIQUE,
    display_name  TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    color         TEXT NOT NULL DEFAULT '#888888'
);

-- Set once a user changes their own password, so the seeder stops overwriting it.
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_changed BOOLEAN NOT NULL DEFAULT FALSE;

CREATE TABLE IF NOT EXISTS matches (
    id           INTEGER PRIMARY KEY,
    home_team    TEXT NOT NULL,
    away_team    TEXT NOT NULL,
    home_flag    TEXT,
    away_flag    TEXT,
    stage        TEXT NOT NULL,
    group_label  TEXT,
    kickoff      TIMESTAMPTZ NOT NULL,
    stadium_name TEXT,
    city         TEXT,
    country      TEXT
);

CREATE TABLE IF NOT EXISTS meetups (
    id              SERIAL PRIMARY KEY,
    match_id        INTEGER NOT NULL REFERENCES matches(id),
    location_name   TEXT NOT NULL,
    location_url    TEXT,
    location_is_bar BOOLEAN NOT NULL DEFAULT FALSE,
    note            TEXT NOT NULL DEFAULT '',
    created_by      INTEGER NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_meetups_match ON meetups(match_id);

CREATE TABLE IF NOT EXISTS meetup_invites (
    meetup_id INTEGER NOT NULL REFERENCES meetups(id) ON DELETE CASCADE,
    user_id   INTEGER NOT NULL REFERENCES users(id),
    PRIMARY KEY (meetup_id, user_id)
);

CREATE TABLE IF NOT EXISTS meetup_mentions (
    meetup_id INTEGER NOT NULL REFERENCES meetups(id) ON DELETE CASCADE,
    user_id   INTEGER NOT NULL REFERENCES users(id),
    PRIMARY KEY (meetup_id, user_id)
);

CREATE TABLE IF NOT EXISTS meetup_members (
    meetup_id INTEGER NOT NULL REFERENCES meetups(id) ON DELETE CASCADE,
    user_id   INTEGER NOT NULL REFERENCES users(id),
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (meetup_id, user_id)
);
