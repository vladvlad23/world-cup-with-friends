package seed

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

//go:embed data/*.json
var dataFS embed.FS

// rawMatch mirrors the upstream football.matches.json shape (all strings).
type rawMatch struct {
	ID         string `json:"id"`
	HomeTeamID string `json:"home_team_id"`
	AwayTeamID string `json:"away_team_id"`
	Group      string `json:"group"`
	LocalDate  string `json:"local_date"`
	StadiumID  string `json:"stadium_id"`
	Type       string `json:"type"`
}

type rawTeam struct {
	ID     string `json:"id"`
	NameEn string `json:"name_en"`
	Flag   string `json:"flag"`
}

type rawStadium struct {
	ID        string `json:"id"`
	NameEn    string `json:"name_en"`
	CityEn    string `json:"city_en"`
	CountryEn string `json:"country_en"`
}

type seedUser struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
	Color       string `json:"color"`
}

// Run seeds users (always upserted) and matches (only if the table is empty).
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	if err := seedUsers(ctx, pool); err != nil {
		return fmt.Errorf("seed users: %w", err)
	}
	if err := seedMatches(ctx, pool); err != nil {
		return fmt.Errorf("seed matches: %w", err)
	}
	return nil
}

func seedUsers(ctx context.Context, pool *pgxpool.Pool) error {
	users, err := loadSeedUsers()
	if err != nil {
		return err
	}
	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		// Display name and colour always track the seed file, but the password
		// is only (re)set from the seed while the user hasn't changed it
		// themselves — otherwise a restart would wipe out their new password.
		_, err = pool.Exec(ctx, `
			INSERT INTO users (username, display_name, password_hash, color)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (username)
			DO UPDATE SET display_name = EXCLUDED.display_name,
			              color = EXCLUDED.color,
			              password_hash = CASE WHEN users.password_changed
			                                   THEN users.password_hash
			                                   ELSE EXCLUDED.password_hash END`,
			u.Username, u.DisplayName, string(hash), u.Color)
		if err != nil {
			return err
		}
	}
	return nil
}

// loadSeedUsers reads the user list to seed. In production set USERS_FILE to a
// path (mounted, uncommitted) so real credentials never live in the repo or
// image. Without it, the embedded demo users are used (local dev / tests).
func loadSeedUsers() ([]seedUser, error) {
	var users []seedUser
	if path := os.Getenv("USERS_FILE"); path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read USERS_FILE %s: %w", path, err)
		}
		if err := json.Unmarshal(b, &users); err != nil {
			return nil, fmt.Errorf("parse USERS_FILE %s: %w", path, err)
		}
		if len(users) == 0 {
			return nil, fmt.Errorf("USERS_FILE %s contains no users", path)
		}
		log.Printf("seeding %d users from USERS_FILE=%s", len(users), path)
		return users, nil
	}
	if err := readJSON("data/users.json", &users); err != nil {
		return nil, err
	}
	log.Printf("seeding %d demo users (embedded — set USERS_FILE to override in production)", len(users))
	return users, nil
}

func seedMatches(ctx context.Context, pool *pgxpool.Pool) error {
	var count int
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM matches`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	var matches []rawMatch
	if err := readJSON("data/football.matches.json", &matches); err != nil {
		return err
	}
	var teams []rawTeam
	if err := readJSON("data/football.teams.json", &teams); err != nil {
		return err
	}
	var stadiums []rawStadium
	if err := readJSON("data/football.stadiums.json", &stadiums); err != nil {
		return err
	}

	teamByID := map[string]rawTeam{}
	for _, t := range teams {
		teamByID[t.ID] = t
	}
	stadiumByID := map[string]rawStadium{}
	for _, s := range stadiums {
		stadiumByID[s.ID] = s
	}

	batch := &pgx.Batch{}
	for _, m := range matches {
		kickoff, err := time.Parse("01/02/2006 15:04", m.LocalDate)
		if err != nil {
			return fmt.Errorf("parse date %q: %w", m.LocalDate, err)
		}
		home := teamName(teamByID, m.HomeTeamID)
		away := teamName(teamByID, m.AwayTeamID)
		st := stadiumByID[m.StadiumID]
		batch.Queue(`
			INSERT INTO matches
			  (id, home_team, away_team, home_flag, away_flag, stage, group_label, kickoff, stadium_name, city, country)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
			ON CONFLICT (id) DO NOTHING`,
			atoi(m.ID), home, away,
			teamByID[m.HomeTeamID].Flag, teamByID[m.AwayTeamID].Flag,
			m.Type, m.Group, kickoff,
			st.NameEn, st.CityEn, st.CountryEn)
	}
	return pool.SendBatch(ctx, batch).Close()
}

// teamName resolves a team id to a display name, mapping the upstream "0"
// placeholder (knockout slots not yet decided) to "TBD".
func teamName(byID map[string]rawTeam, id string) string {
	if id == "0" || id == "" {
		return "TBD"
	}
	if t, ok := byID[id]; ok {
		return t.NameEn
	}
	return "TBD"
}

func readJSON(path string, dst any) error {
	b, err := dataFS.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}

func atoi(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			continue
		}
		n = n*10 + int(c-'0')
	}
	return n
}
