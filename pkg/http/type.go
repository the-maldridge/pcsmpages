package http

import (
	"net/http"

	"github.com/the-maldridge/pcsmpages/pkg/pcsm"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/hashicorp/go-hclog"
)

// Server wraps up all the request routers and associated components
// that serve various parts of the nbuild stack.
type Server struct {
	l     hclog.Logger
	r     chi.Router
	n     *http.Server
	p     *pcsm.Client
	tmpls *pongo2.TemplateSet
}
