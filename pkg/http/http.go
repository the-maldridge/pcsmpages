package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New initializes the server with its default routers.
func New(opts ...Option) (*Server, error) {
	sbl, err := pongo2.NewSandboxedFilesystemLoader("theme/p2")
	if err != nil {
		return nil, err
	}

	s := Server{
		r:     chi.NewRouter(),
		n:     &http.Server{},
		tmpls: pongo2.NewSet("html", sbl),
	}

	// Apply options
	for _, o := range opts {
		o(&s)
	}

	s.tmpls.Debug = true

	s.r.Use(middleware.Logger)
	s.r.Use(middleware.Heartbeat("/healthz"))
	s.fileServer(s.r, "/static", http.Dir("theme/static"))

	s.r.Get("/debug", s.debug)
	s.r.Get("/bar/field/{fnum}", s.barField)
	s.r.Get("/box/match", s.boxMatch)
	s.r.Get("/q/{offset}", s.qDisplay)
	s.r.Get("/clock", s.clock)
	s.r.Get("/clock/end", s.clockEnd)

	s.r.Get("/scoreboard", s.scoreboardAll)
	s.r.Get("/scoreboard/page/{page}/{count}", s.scoreboardPage)

	return &s, nil
}

// Serve binds, initializes the mux, and serves forever.
func (s *Server) Serve(bind string) error {
	s.l.Info("HTTP is starting")
	s.n.Addr = bind
	s.n.Handler = s.r
	return s.n.ListenAndServe()
}

// Shutdown requests the underlying server gracefully cease operation.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.n.Shutdown(ctx)
}

func (s *Server) templateErrorHandler(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, "Error while rendering template: %s\n", err)
}

func (s *Server) fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (s *Server) doTemplate(w http.ResponseWriter, r *http.Request, tmpl string, ctx pongo2.Context) {
	if ctx == nil {
		ctx = pongo2.Context{}
	}
	t, err := s.tmpls.FromCache(tmpl)
	if err != nil {
		s.templateErrorHandler(w, err)
		return
	}
	if err := t.ExecuteWriter(ctx, w); err != nil {
		s.templateErrorHandler(w, err)
	}
}

func (s *Server) debug(w http.ResponseWriter, r *http.Request) {
	m, err := s.p.GetScoreboard()
	if err != nil {
		s.l.Error("Error getting current match", "error", err)
	}

	json.NewEncoder(w).Encode(m)
}

func (s *Server) barField(w http.ResponseWriter, r *http.Request) {
	m, err := s.p.GetCurrentMatch()
	if err != nil {
		s.l.Error("Error getting current match", "error", err)
		s.doTemplate(w, r, "views/fieldBar.p2", pongo2.Context{})
		return
	}

	fnum, err := strconv.Atoi(chi.URLParam(r, "fnum"))
	if err != nil {
		s.l.Error("Bad field number!", "fnum", chi.URLParam(r, "fnum"))
		s.doTemplate(w, r, "views/fieldBar.p2", pongo2.Context{})
		return
	}

	ctx := pongo2.Context{"Phase": m.Phase, "Number": m.Number}
	for _, f := range m.Fields {
		s.l.Debug("Field", "want", fnum, "got", f.Number)
		if f.Number == fnum {
			ctx["Teams"] = f.Teams
			break
		}
	}
	s.l.Debug("Context", "ctx", ctx)
	s.doTemplate(w, r, "views/fieldBar.p2", ctx)
}

func (s *Server) boxMatch(w http.ResponseWriter, r *http.Request) {
	m, err := s.p.GetCurrentMatch()
	if err != nil {
		s.l.Error("Error getting current match", "error", err)
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err})
		return
	}

	s.doTemplate(w, r, "views/boxMatch.p2", pongo2.Context{"Match": m})
}

func (s *Server) qDisplay(w http.ResponseWriter, r *http.Request) {
	m, err := s.p.GetCurrentMatch()
	if err != nil {
		s.l.Error("Error getting current match", "error", err)
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err})
		return
	}

	offset, err := strconv.Atoi(chi.URLParam(r, "offset"))
	if err != nil {
		s.l.Error("Bad offset number!", "offset", chi.URLParam(r, "offset"))
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err})
		return
	}

	sch, err := s.p.GetSchedule(m.Phase)
	if err != nil {
		s.l.Error("Error getting current schedule", "error", err)
		s.doTemplate(w, r, "errors/internal.p2", pongo2.Context{"error": err})
		return
	}

	ctx := pongo2.Context{
		"Phase": m.Phase,
		"Running": m,
		"Refresh": r.URL.Query().Get("dwell"),
	}

	if m.Number+offset < len(sch)+1 {
		for _, sm := range sch {
			if sm.Number == m.Number+offset {
				ctx["QueueHere"] = sm
				break
			}
		}
	}
	if r.URL.Query().Get("hidecurrent") != "" {
		delete(ctx, "Running")
	}

	s.l.Debug("Context", "ctx", ctx)
	s.doTemplate(w, r, "views/qHUD.p2", ctx)
}

func (s *Server) clock(w http.ResponseWriter, r *http.Request) {
	s.doTemplate(w, r, "views/clock.p2", pongo2.Context{})
}

func (s *Server) clockEnd(w http.ResponseWriter, r *http.Request) {
	m, err := s.p.GetCurrentMatch()
	if err != nil {
		s.l.Error("Error getting current match", "error", err)
	}

	d := struct {
		MatchEnd time.Time
	}{
		MatchEnd: m.End.ToTime(),
	}

	json.NewEncoder(w).Encode(d)
}
