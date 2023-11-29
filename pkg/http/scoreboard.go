package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2/v4"
	"github.com/go-chi/chi/v5"
)

func (s *Server) scoreboardAll(w http.ResponseWriter, r *http.Request) {
	sb, err := s.p.GetScoreboard()
	if err != nil {
		s.l.Error("Error getting scoreboard data", "error", err)
	}

	ctx := pongo2.Context{
		"Board":        sb,
		"ShowDivision": r.URL.Query().Get("showdivision") != "",
	}

	s.doTemplate(w, r, "views/scoreboard.p2", ctx)
}

func (s *Server) scoreboardPage(w http.ResponseWriter, r *http.Request) {
	sb, err := s.p.GetScoreboard()
	if err != nil {
		s.l.Error("Error getting scoreboard data", "error", err)
	}

	page, err := strconv.Atoi(chi.URLParam(r, "page"))
	if err != nil {
		s.l.Error("Bad scoreboard page number", "error", err)
		return
	}

	count, err := strconv.Atoi(chi.URLParam(r, "count"))
	if err != nil {
		s.l.Error("Bad scoreboard count", "error", err)
		return
	}

	start := (page * count)
	end := start + count
	nextpage := page + 1
	if end >= len(sb.Teams) {
		nextpage = 0
		end = len(sb.Teams)-1
	}
	sb.Teams = sb.Teams[start:end]

	ctx := pongo2.Context{
		"Board":        sb,
		"ShowDivision": r.URL.Query().Get("showdivision") != "",
		"Refresh":      r.URL.Query().Get("dwell"),
		"Next":         fmt.Sprintf("/scoreboard/page/%d/%d?%s", nextpage, count, r.URL.Query().Encode()),
	}

	s.doTemplate(w, r, "views/scoreboard.p2", ctx)

}
