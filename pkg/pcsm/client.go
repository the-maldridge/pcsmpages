package pcsm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

// New returns and configures a client.
func New(opts ...Option) *Client {
	c := &Client{
		l:    hclog.NewNullLogger(),
		c:    http.Client{Timeout: time.Second * 5},
		addr: "localhost:9268",
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// GetCurrentMatch is a convenience function to retrieve the match
// that the system believes is current.
func (c *Client) GetCurrentMatch() (*Match, error) {
	return c.GetMatch("", 0)
}

// GetMatch returns the specified phase and match.
func (c *Client) GetMatch(phase string, number int) (*Match, error) {
	nElement := ""
	if number > 0 {
		nElement = strconv.Itoa(number)
	}

	p := path.Join(path.Clean(path.Join("api/public/match", phase, nElement)))
	url := fmt.Sprintf("http://%s/%s", c.addr, p)
	c.l.Debug("Fetching match", "url", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	m := new(Match)
	if err := json.NewDecoder(resp.Body).Decode(m); err != nil {
		c.l.Error("Error decoding json", "error", err)
		return nil, err
	}

	return m, nil
}

// GetCurrentMatchFromSchedule returns whatever match the schedule
// thinks is either Staged or Started.
func (c *Client) GetCurrentMatchFromSchedule() (*Match, error) {
	s, err := c.GetSchedule("")
	if err != nil {
		return nil, err
	}

	for _, m := range s {
		if m.State == "Staged" || m.State == "Started" {
			return &m, nil
		}
	}
	return new(Match), nil
}

// GetSchedule loads the current schedule from PCSM.
func (c *Client) GetSchedule(phase string) (Schedule, error) {
	p := path.Join(path.Clean(path.Join("api/public/schedule", phase)))
	url := fmt.Sprintf("http://%s/%s", c.addr, p)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	s := Schedule{}
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		c.l.Error("Error decoding json", "error", err)
		return nil, err
	}

	sort.Slice(s, func(i, j int) bool {
		return s[i].Number < s[j].Number
	})

	return s, nil
}
