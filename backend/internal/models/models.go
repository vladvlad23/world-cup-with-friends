package models

import "time"

type User struct {
	ID          int    `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Color       string `json:"color"`
}

type Match struct {
	ID          int       `json:"id"`
	HomeTeam    string    `json:"home_team"`
	AwayTeam    string    `json:"away_team"`
	HomeFlag    string    `json:"home_flag"`
	AwayFlag    string    `json:"away_flag"`
	Stage       string    `json:"stage"`
	GroupLabel  string    `json:"group_label"`
	Kickoff     time.Time `json:"kickoff"`
	StadiumName string    `json:"stadium_name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	MeetupCount int       `json:"meetup_count"`
}

// MeetupSummary is the lightweight shape shown in lists.
type MeetupSummary struct {
	ID            int    `json:"id"`
	MatchID       int    `json:"match_id"`
	LocationName  string `json:"location_name"`
	LocationIsBar bool   `json:"location_is_bar"`
	CreatedBy     User   `json:"created_by"`
	MemberCount   int    `json:"member_count"`
}

// Meetup is the full detail shape.
type Meetup struct {
	ID            int       `json:"id"`
	MatchID       int       `json:"match_id"`
	LocationName  string    `json:"location_name"`
	LocationURL   string    `json:"location_url"`
	LocationIsBar bool      `json:"location_is_bar"`
	Note          string    `json:"note"`
	CreatedBy     User      `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	Invites       []User    `json:"invites"`
	Mentions      []User    `json:"mentions"`
	Members       []User    `json:"members"`
}
