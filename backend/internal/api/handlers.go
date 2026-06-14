package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"worldcup/backend/internal/auth"
	"worldcup/backend/internal/models"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}

	var (
		id    int
		hash  string
		user  models.User
	)
	err := s.pool.QueryRow(r.Context(),
		`SELECT id, username, display_name, color, password_hash FROM users WHERE username = $1`,
		body.Username).Scan(&id, &user.Username, &user.DisplayName, &user.Color, &hash)
	if errors.Is(err, pgx.ErrNoRows) || (err == nil && bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)) != nil) {
		writeErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	user.ID = id

	token, err := s.auth.Issue(id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "token error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	u, err := s.getUser(r.Context(), uid)
	if err != nil {
		writeErr(w, http.StatusNotFound, "user not found")
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.listUsers(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (s *Server) handleMatches(w http.ResponseWriter, r *http.Request) {
	rows, err := s.pool.Query(r.Context(), `
		SELECT m.id, m.home_team, m.away_team, m.home_flag, m.away_flag, m.stage,
		       m.group_label, m.kickoff, m.stadium_name, m.city, m.country,
		       count(mu.id) AS meetup_count
		FROM matches m
		LEFT JOIN meetups mu ON mu.match_id = m.id
		GROUP BY m.id
		ORDER BY m.kickoff, m.id`)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	defer rows.Close()

	matches := []models.Match{}
	for rows.Next() {
		var m models.Match
		if err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.HomeFlag, &m.AwayFlag,
			&m.Stage, &m.GroupLabel, &m.Kickoff, &m.StadiumName, &m.City, &m.Country,
			&m.MeetupCount); err != nil {
			writeErr(w, http.StatusInternalServerError, "scan error")
			return
		}
		matches = append(matches, m)
	}
	writeJSON(w, http.StatusOK, matches)
}

func (s *Server) handleMatch(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	var m models.Match
	err = s.pool.QueryRow(r.Context(), `
		SELECT id, home_team, away_team, home_flag, away_flag, stage, group_label,
		       kickoff, stadium_name, city, country
		FROM matches WHERE id = $1`, id).Scan(
		&m.ID, &m.HomeTeam, &m.AwayTeam, &m.HomeFlag, &m.AwayFlag, &m.Stage,
		&m.GroupLabel, &m.Kickoff, &m.StadiumName, &m.City, &m.Country)
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, "match not found")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	meetups, err := s.listMeetupSummaries(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"match": m, "meetups": meetups})
}

func (s *Server) handleCreateMeetup(w http.ResponseWriter, r *http.Request) {
	matchID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	var body struct {
		LocationName   string `json:"location_name"`
		LocationURL    string `json:"location_url"`
		LocationIsBar  bool   `json:"location_is_bar"`
		Note           string `json:"note"`
		InviteUserIDs  []int  `json:"invite_user_ids"`
		MentionUserIDs []int  `json:"mention_user_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.LocationName == "" {
		writeErr(w, http.StatusBadRequest, "location_name is required")
		return
	}
	uid := auth.UserID(r.Context())

	tx, err := s.pool.Begin(r.Context())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	defer tx.Rollback(r.Context())

	var meetupID int
	err = tx.QueryRow(r.Context(), `
		INSERT INTO meetups (match_id, location_name, location_url, location_is_bar, note, created_by)
		VALUES ($1,$2,$3,$4,$5,$6) RETURNING id`,
		matchID, body.LocationName, nullable(body.LocationURL), body.LocationIsBar, body.Note, uid).Scan(&meetupID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not create meetup")
		return
	}

	for _, inv := range body.InviteUserIDs {
		if _, err := tx.Exec(r.Context(),
			`INSERT INTO meetup_invites (meetup_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			meetupID, inv); err != nil {
			writeErr(w, http.StatusInternalServerError, "invite error")
			return
		}
	}
	for _, mn := range body.MentionUserIDs {
		if _, err := tx.Exec(r.Context(),
			`INSERT INTO meetup_mentions (meetup_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
			meetupID, mn); err != nil {
			writeErr(w, http.StatusInternalServerError, "mention error")
			return
		}
	}
	// Creator auto-joins.
	if _, err := tx.Exec(r.Context(),
		`INSERT INTO meetup_members (meetup_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		meetupID, uid); err != nil {
		writeErr(w, http.StatusInternalServerError, "member error")
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		writeErr(w, http.StatusInternalServerError, "commit error")
		return
	}

	m, err := s.getMeetup(r.Context(), meetupID)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (s *Server) handleMeetup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	m, err := s.getMeetup(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		writeErr(w, http.StatusNotFound, "meetup not found")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (s *Server) handleJoinMeetup(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}
	uid := auth.UserID(r.Context())
	_, err = s.pool.Exec(r.Context(),
		`INSERT INTO meetup_members (meetup_id, user_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		id, uid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not join")
		return
	}
	m, err := s.getMeetup(r.Context(), id)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}
