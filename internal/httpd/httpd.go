package httpd

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"

	"github.com/buglloc/vanityd/internal/projects"
)

//go:embed static
var staticFS embed.FS

type HttpD struct {
	prj  projects.Provider
	zone string
	http http.Server
}

func NewHttpD(prj projects.Provider, opts ...Option) *HttpD {
	h := &HttpD{
		prj:  prj,
		zone: "." + DefaultZone,
		http: http.Server{
			Addr: DefaultAddr,
		},
	}

	for _, opt := range opts {
		opt(h)
	}

	h.http.Handler = h.buildRouter()
	return h
}

func (h *HttpD) buildRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.RequestID,
		Logger(log.Logger),
		middleware.Recoverer,
	)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		host := requestHost(r)
		if strings.HasSuffix(host, h.zone) {
			h.replyProject(w, r, strings.TrimSuffix(host, h.zone))
			return
		}

		h.replyIndex(w, r)
	})

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))
	return r
}

func (h *HttpD) ListenAndServe() error {
	log.Info().Str("addr", h.http.Addr).Msg("listening")

	return h.http.ListenAndServe()
}

func (h *HttpD) Shutdown(ctx context.Context) error {
	return h.http.Shutdown(ctx)
}

func (h *HttpD) replyIndex(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := indexTmpl.Execute(&buf, h.prj.List()); err != nil {
		log.Error().
			Err(err).
			Msg("unable to render index page")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *HttpD) replyProject(w http.ResponseWriter, r *http.Request, prjName string) {
	prj, err := h.prj.Get(prjName)
	if err != nil {
		if errors.Is(err, projects.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Error().
			Err(err).
			Str("prj", prjName).
			Msg("unable to get project info")
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.Redirect(w, r, prj.URL, http.StatusTemporaryRedirect)
}

func requestHost(r *http.Request) (host string) {
	// not standard, but most popular
	host = r.Header.Get("X-Forwarded-Host")
	if host != "" {
		return
	}

	// RFC 7239
	host = r.Header.Get("Forwarded")
	_, _, host = parseForwarded(host)
	if host != "" {
		return
	}

	// if all else fails fall back to request host
	host = r.Host
	return
}

func parseForwarded(forwarded string) (addr, proto, host string) {
	if forwarded == "" {
		return
	}
	for _, forwardedPair := range strings.Split(forwarded, ";") {
		if tv := strings.SplitN(forwardedPair, "=", 2); len(tv) == 2 {
			token, value := tv[0], tv[1]
			token = strings.TrimSpace(token)
			value = strings.TrimSpace(strings.Trim(value, `"`))
			switch strings.ToLower(token) {
			case "for":
				addr = value
			case "proto":
				proto = value
			case "host":
				host = value
			}

		}
	}
	return
}
