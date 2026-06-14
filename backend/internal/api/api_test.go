package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"worldcup/backend/internal/auth"
	"worldcup/backend/internal/db"
	"worldcup/backend/internal/models"
	"worldcup/backend/internal/seed"
)

// testServer spins up a Server backed by the database in DATABASE_URL. The
// schema is applied and seeded, and meetup tables are truncated so each run is
// deterministic. If no database is reachable the test is skipped (CI provides
// one as a service).
func testServer(t *testing.T) (http.Handler, *pgxpool.Pool) {
	t.Helper()
	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = "postgres://worldcup:worldcup@localhost:5432/worldcup?sslmode=disable"
	}
	ctx := context.Background()
	pool, err := db.Connect(ctx, url)
	if err != nil {
		t.Skipf("no database reachable (%v); set DATABASE_URL to run integration tests", err)
	}
	if err := seed.Run(ctx, pool); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if _, err := pool.Exec(ctx, `TRUNCATE meetups RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	t.Cleanup(pool.Close)
	return NewServer(pool, auth.New("test-secret")).Router(), pool
}

func do(t *testing.T, h http.Handler, method, path, token string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func login(t *testing.T, h http.Handler, username, password string) (string, models.User) {
	t.Helper()
	rec := do(t, h, http.MethodPost, "/api/auth/login", "", map[string]string{
		"username": username, "password": password,
	})
	if rec.Code != http.StatusOK {
		t.Fatalf("login %s: status %d (%s)", username, rec.Code, rec.Body.String())
	}
	var out struct {
		Token string      `json:"token"`
		User  models.User `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	return out.Token, out.User
}

func TestLogin(t *testing.T) {
	h, _ := testServer(t)

	t.Run("success", func(t *testing.T) {
		token, user := login(t, h, "alice", "password")
		if token == "" {
			t.Fatal("expected a token")
		}
		if user.Username != "alice" {
			t.Fatalf("user = %q, want alice", user.Username)
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		rec := do(t, h, http.MethodPost, "/api/auth/login", "", map[string]string{
			"username": "alice", "password": "nope",
		})
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})

	t.Run("unknown user", func(t *testing.T) {
		rec := do(t, h, http.MethodPost, "/api/auth/login", "", map[string]string{
			"username": "ghost", "password": "password",
		})
		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want 401", rec.Code)
		}
	})
}

func TestProtectedRoutesRequireAuth(t *testing.T) {
	h, _ := testServer(t)
	for _, path := range []string{"/api/me", "/api/users", "/api/matches", "/api/matches/1"} {
		rec := do(t, h, http.MethodGet, path, "", nil)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("GET %s without token = %d, want 401", path, rec.Code)
		}
	}
}

func TestMatchesSeeded(t *testing.T) {
	h, _ := testServer(t)
	token, _ := login(t, h, "alice", "password")

	rec := do(t, h, http.MethodGet, "/api/matches", token, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var matches []models.Match
	if err := json.Unmarshal(rec.Body.Bytes(), &matches); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(matches) != 104 {
		t.Fatalf("got %d matches, want 104", len(matches))
	}
	if matches[0].HomeTeam == "" || matches[0].Kickoff.IsZero() {
		t.Fatalf("first match not populated: %+v", matches[0])
	}
}

func TestMeetupLifecycle(t *testing.T) {
	h, _ := testServer(t)
	aliceTok, alice := login(t, h, "alice", "password")
	bobTok, bob := login(t, h, "bob", "password")
	_, carol := login(t, h, "carol", "password")

	// Alice creates a bar meetup on match 1, inviting bob and mentioning carol.
	rec := do(t, h, http.MethodPost, "/api/matches/1/meetups", aliceTok, map[string]any{
		"location_name":   "The Anchor",
		"location_url":    "https://maps.example/anchor",
		"location_is_bar": true,
		"note":            "come thirsty @carol",
		"invite_user_ids": []int{bob.ID},
		"mention_user_ids": []int{carol.ID},
	})
	if rec.Code != http.StatusCreated {
		t.Fatalf("create meetup: status %d (%s)", rec.Code, rec.Body.String())
	}
	var created models.Meetup
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create: %v", err)
	}
	if !created.LocationIsBar {
		t.Error("expected location_is_bar true")
	}
	if created.CreatedBy.ID != alice.ID {
		t.Errorf("creator = %d, want %d", created.CreatedBy.ID, alice.ID)
	}
	if len(created.Members) != 1 || created.Members[0].ID != alice.ID {
		t.Errorf("creator should be auto-joined, members = %+v", created.Members)
	}
	if len(created.Invites) != 1 || created.Invites[0].ID != bob.ID {
		t.Errorf("invites = %+v, want [bob]", created.Invites)
	}
	if len(created.Mentions) != 1 || created.Mentions[0].ID != carol.ID {
		t.Errorf("mentions = %+v, want [carol]", created.Mentions)
	}

	// The match now reports a meetup count.
	rec = do(t, h, http.MethodGet, "/api/matches", aliceTok, nil)
	var matches []models.Match
	_ = json.Unmarshal(rec.Body.Bytes(), &matches)
	var found bool
	for _, m := range matches {
		if m.ID == 1 {
			found = true
			if m.MeetupCount != 1 {
				t.Errorf("match 1 meetup_count = %d, want 1", m.MeetupCount)
			}
		}
	}
	if !found {
		t.Fatal("match 1 not in list")
	}

	// Bob joins.
	rec = do(t, h, http.MethodPost, "/api/meetups/"+itoa(created.ID)+"/join", bobTok, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("join: status %d (%s)", rec.Code, rec.Body.String())
	}
	var joined models.Meetup
	_ = json.Unmarshal(rec.Body.Bytes(), &joined)
	if len(joined.Members) != 2 {
		t.Fatalf("members after join = %d, want 2", len(joined.Members))
	}

	// Joining again is idempotent.
	rec = do(t, h, http.MethodPost, "/api/meetups/"+itoa(created.ID)+"/join", bobTok, nil)
	_ = json.Unmarshal(rec.Body.Bytes(), &joined)
	if len(joined.Members) != 2 {
		t.Fatalf("re-join changed member count to %d, want 2", len(joined.Members))
	}
}

func TestCreateMeetupValidation(t *testing.T) {
	h, _ := testServer(t)
	token, _ := login(t, h, "alice", "password")
	rec := do(t, h, http.MethodPost, "/api/matches/1/meetups", token, map[string]any{
		"location_name": "",
	})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("empty location_name status = %d, want 400", rec.Code)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	return string(b[i:])
}
