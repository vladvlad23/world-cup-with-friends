package seed

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAtoi(t *testing.T) {
	cases := map[string]int{
		"1":    1,
		"104":  104,
		"0":    0,
		"":     0,
		"12a3": 123, // non-digits skipped
	}
	for in, want := range cases {
		if got := atoi(in); got != want {
			t.Errorf("atoi(%q) = %d, want %d", in, got, want)
		}
	}
}

func TestTeamName(t *testing.T) {
	byID := map[string]rawTeam{
		"1": {ID: "1", NameEn: "Mexico"},
		"2": {ID: "2", NameEn: "South Africa"},
	}
	cases := map[string]string{
		"1":   "Mexico",       // known team
		"2":   "South Africa", // known team
		"0":   "TBD",          // knockout placeholder
		"":    "TBD",          // empty
		"999": "TBD",          // unknown id
	}
	for id, want := range cases {
		if got := teamName(byID, id); got != want {
			t.Errorf("teamName(%q) = %q, want %q", id, got, want)
		}
	}
}

func TestDateParsing(t *testing.T) {
	// The seed parses the upstream "MM/DD/YYYY HH:MM" local_date format.
	got, err := time.Parse("01/02/2006 15:04", "06/11/2026 13:00")
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	want := time.Date(2026, time.June, 11, 13, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Fatalf("parsed %v, want %v", got, want)
	}
}

func TestLoadSeedUsersEmbeddedFallback(t *testing.T) {
	t.Setenv("USERS_FILE", "") // ensure unset -> embedded demo users
	users, err := loadSeedUsers()
	if err != nil {
		t.Fatalf("loadSeedUsers: %v", err)
	}
	if len(users) == 0 {
		t.Fatal("expected embedded demo users")
	}
}

func TestLoadSeedUsersFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "users.json")
	if err := os.WriteFile(path, []byte(`[{"username":"realvlad","display_name":"Real Vlad","password":"s3cret","color":"#111"}]`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("USERS_FILE", path)

	users, err := loadSeedUsers()
	if err != nil {
		t.Fatalf("loadSeedUsers: %v", err)
	}
	if len(users) != 1 || users[0].Username != "realvlad" {
		t.Fatalf("got %+v, want single realvlad user", users)
	}
}

func TestLoadSeedUsersFileErrors(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		t.Setenv("USERS_FILE", filepath.Join(t.TempDir(), "nope.json"))
		if _, err := loadSeedUsers(); err == nil {
			t.Fatal("expected error for missing USERS_FILE")
		}
	})
	t.Run("empty list", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "empty.json")
		_ = os.WriteFile(path, []byte(`[]`), 0o600)
		t.Setenv("USERS_FILE", path)
		if _, err := loadSeedUsers(); err == nil {
			t.Fatal("expected error for empty USERS_FILE")
		}
	})
}

func TestEmbeddedDataPresent(t *testing.T) {
	// The fixture data the seed relies on must be embedded in the binary.
	for _, f := range []string{
		"data/users.json",
		"data/football.matches.json",
		"data/football.teams.json",
		"data/football.stadiums.json",
	} {
		if _, err := dataFS.ReadFile(f); err != nil {
			t.Errorf("missing embedded file %s: %v", f, err)
		}
	}
}
