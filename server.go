package probes

import (
	"context"
	"net/http"
)

type ServerOpts struct {
	Addr   string
	Probes *Probes

	// The path for readiness probes (default: /readyz)
	ReadyPath string

	// The path for health probes (default: healthz)
	HealthPath string

	// The path for startup probes (default: startedz)
	StartupPath string
}

type Server struct {
	probes *Probes
	server *http.Server
}

func NewServer(opts ServerOpts) *Server {
	if opts.ReadyPath == "" {
		opts.ReadyPath = "/readyz"
	}
	if opts.HealthPath == "" {
		opts.HealthPath = "/healthz"
	}
	if opts.StartupPath == "" {
		opts.StartupPath = "/startedz"
	}
	s := &Server{
		probes: opts.Probes,
	}
	mux := http.NewServeMux()
	mux.HandleFunc(opts.HealthPath, s.healthy)
	mux.HandleFunc(opts.ReadyPath, s.ready)
	mux.HandleFunc(opts.StartupPath, s.started)
	s.server = &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}

	return s
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) healthy(w http.ResponseWriter, r *http.Request) {
	if s.probes.IsHealthy() {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func (s *Server) ready(w http.ResponseWriter, r *http.Request) {
	if s.probes.IsReady() {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func (s *Server) started(w http.ResponseWriter, r *http.Request) {
	if s.probes.IsStarted() {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}
