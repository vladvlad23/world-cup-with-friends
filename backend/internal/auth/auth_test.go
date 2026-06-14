package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIssueParseRoundTrip(t *testing.T) {
	s := New("secret")
	token, err := s.Issue(42)
	if err != nil {
		t.Fatalf("issue: %v", err)
	}
	uid, err := s.parse(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if uid != 42 {
		t.Fatalf("got uid %d, want 42", uid)
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	token, _ := New("secret-a").Issue(1)
	if _, err := New("secret-b").parse(token); err == nil {
		t.Fatal("expected error for token signed with a different secret")
	}
}

func TestParseRejectsGarbage(t *testing.T) {
	if _, err := New("secret").parse("not.a.jwt"); err == nil {
		t.Fatal("expected error for malformed token")
	}
}

func TestMiddleware(t *testing.T) {
	s := New("secret")
	token, _ := s.Issue(7)

	var seenUID int
	handler := s.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		seenUID = UserID(r.Context())
	}))

	cases := []struct {
		name       string
		header     string
		wantStatus int
		wantUID    int
	}{
		{"valid", "Bearer " + token, http.StatusOK, 7},
		{"missing", "", http.StatusUnauthorized, 0},
		{"no bearer prefix", token, http.StatusUnauthorized, 0},
		{"invalid token", "Bearer garbage", http.StatusUnauthorized, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			seenUID = 0
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.header != "" {
				req.Header.Set("Authorization", tc.header)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if seenUID != tc.wantUID {
				t.Fatalf("uid = %d, want %d", seenUID, tc.wantUID)
			}
		})
	}
}
