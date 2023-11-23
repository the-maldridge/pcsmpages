package http

import (
	"github.com/the-maldridge/pcsmpages/pkg/pcsm"

	"github.com/hashicorp/go-hclog"
)

// Option sets parameters of the server.
type Option func(*Server)

// WithLogger sets the logger on the server.
func WithLogger(l hclog.Logger) Option {
	return func(s *Server) { s.l = l.Named("http") }
}

// WithPCSM connects the pages server to a PCSM data source.
func WithPCSM(p *pcsm.Client) Option {
	return func(s *Server) { s.p = p }
}
