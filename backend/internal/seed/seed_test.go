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

func TestKickoffInstant(t *testing.T) {
	cases := []struct {
		name      string
		localDate string
		stadiumID string
		wantUTC   time.Time
	}{
		{
			// Houston (NRG, stadium 5) is US Central; June -> CDT (UTC-5).
			"houston CDT", "06/23/2026 12:00", "5",
			time.Date(2026, time.June, 23, 17, 0, 0, 0, time.UTC),
		},
		{
			// Mexico City (Estadio Azteca, stadium 1) is UTC-6, no DST.
			"mexico city", "06/11/2026 13:00", "1",
			time.Date(2026, time.June, 11, 19, 0, 0, 0, time.UTC),
		},
		{
			// Los Angeles (SoFi, stadium 16) is Pacific; June -> PDT (UTC-7).
			"los angeles PDT", "06/12/2026 12:00", "16",
			time.Date(2026, time.June, 12, 19, 0, 0, 0, time.UTC),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := kickoffInstant(tc.localDate, tc.stadiumID)
			if err != nil {
				t.Fatalf("kickoffInstant: %v", err)
			}
			if !got.UTC().Equal(tc.wantUTC) {
				t.Fatalf("got %v, want %v", got.UTC(), tc.wantUTC)
			}
		})
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
