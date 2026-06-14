package api

import (
	"context"

	"worldcup/backend/internal/models"
)

func (s *Server) getUser(ctx context.Context, id int) (models.User, error) {
	var u models.User
	err := s.pool.QueryRow(ctx,
		`SELECT id, username, display_name, color FROM users WHERE id = $1`, id).
		Scan(&u.ID, &u.Username, &u.DisplayName, &u.Color)
	return u, err
}

func (s *Server) listUsers(ctx context.Context) ([]models.User, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, username, display_name, color FROM users ORDER BY display_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Color); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (s *Server) listMeetupSummaries(ctx context.Context, matchID int) ([]models.MeetupSummary, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT mu.id, mu.match_id, mu.location_name, mu.location_is_bar,
		       u.id, u.username, u.display_name, u.color,
		       (SELECT count(*) FROM meetup_members mm WHERE mm.meetup_id = mu.id) AS member_count
		FROM meetups mu
		JOIN users u ON u.id = mu.created_by
		WHERE mu.match_id = $1
		ORDER BY mu.created_at`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.MeetupSummary{}
	for rows.Next() {
		var m models.MeetupSummary
		if err := rows.Scan(&m.ID, &m.MatchID, &m.LocationName, &m.LocationIsBar,
			&m.CreatedBy.ID, &m.CreatedBy.Username, &m.CreatedBy.DisplayName, &m.CreatedBy.Color,
			&m.MemberCount); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (s *Server) getMeetup(ctx context.Context, id int) (models.Meetup, error) {
	var m models.Meetup
	var locURL *string
	err := s.pool.QueryRow(ctx, `
		SELECT mu.id, mu.match_id, mu.location_name, mu.location_url, mu.location_is_bar,
		       mu.note, mu.created_at,
		       u.id, u.username, u.display_name, u.color
		FROM meetups mu
		JOIN users u ON u.id = mu.created_by
		WHERE mu.id = $1`, id).Scan(
		&m.ID, &m.MatchID, &m.LocationName, &locURL, &m.LocationIsBar,
		&m.Note, &m.CreatedAt,
		&m.CreatedBy.ID, &m.CreatedBy.Username, &m.CreatedBy.DisplayName, &m.CreatedBy.Color)
	if err != nil {
		return m, err
	}
	if locURL != nil {
		m.LocationURL = *locURL
	}

	if m.Invites, err = s.usersVia(ctx, "meetup_invites", id); err != nil {
		return m, err
	}
	if m.Mentions, err = s.usersVia(ctx, "meetup_mentions", id); err != nil {
		return m, err
	}
	if m.Members, err = s.usersVia(ctx, "meetup_members", id); err != nil {
		return m, err
	}
	return m, nil
}

// usersVia returns the users linked to a meetup through one of the join tables.
// table is a fixed internal constant, never user input.
func (s *Server) usersVia(ctx context.Context, table string, meetupID int) ([]models.User, error) {
	q := `SELECT u.id, u.username, u.display_name, u.color
	      FROM ` + table + ` j JOIN users u ON u.id = j.user_id
	      WHERE j.meetup_id = $1 ORDER BY u.display_name`
	rows, err := s.pool.Query(ctx, q, meetupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.Color); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
