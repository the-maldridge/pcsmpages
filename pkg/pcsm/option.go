package pcsm

import (
	"github.com/hashicorp/go-hclog"
)

// WithAddress configures the address of the PCSM instance to interact
// with.
func WithAddress(a string) Option {
	return func(c *Client) {
		if a == "" {
			return
		}
		c.addr = a
	}
}

// WithLogger configures the parent logger for this client.
func WithLogger(l hclog.Logger) Option {
	return func(c *Client) {
		c.l = l.Named("pcsm")
	}
}
