package pcsm

import (
	"net/http"

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

// Match represents a match in any state, possibly across multiple
// fields.
type Match struct {
	Phase  string
	Number int `json:"matchNumber"`
	State  string
	Fields []Field
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
