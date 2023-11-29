package pcsm

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
)

// Client knows how to connect to PCSM, get data from it, and refresh
// that data.
type Client struct {
	l    hclog.Logger
	c    http.Client
	addr string
}

// Option configures the client.
type Option func(c *Client)

// A Schedule is a list of matches that occur in a serial order.
type Schedule []Match

// Time wraps the normal Time to make deserialization work
// correctly for the C# timestamp.
type Time struct {
	time.Time
}

// Match represents a match in any state, possibly across multiple
// fields.
type Match struct {
	Phase         string
	Number        int `json:"matchNumber"`
	State         string
	Start         Time `json:"matchStart"`
	End           Time `json:"matchEnd"`
	TimeRemaining time.Duration
	Fields        []Field
}

// Field has a number and a set of teams.
type Field struct {
	Number int `json:"fieldNumber"`
	Teams  []Team
}

// Team represents a team on the field or in a schedule.
type Team struct {
	Quadrant string

	Name        string
	DisplayName string
	Number      int `json:"teamNumber"`
	Ticker      string
}

// A Scoreboard represents the data necessary to draw a numerically
// sorted set of teams based on the number of points gained.
type Scoreboard struct {
	Teams     []*RankedTeam
	Divisions []*Division
}

// A RankedTeam represent teams that have scoreboard information.
type RankedTeam struct {
	Rank       int
	Page       int
	Name       string
	Number     int `json:"teamNumber"`
	Score      int
	NextMatch  string `json:"nextPhase"`
	Completed  int    `json:"matchesCompleted"`
	Division   string
	DivisionID string `json:"divisionId"`
}

// Division maps between division display names and IDs.
type Division struct {
	ID   string
	Name string
}
