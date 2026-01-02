package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dormitory-life/core/internal/config"
	core "github.com/dormitory-life/core/internal/service"
)

type ServerConfig struct {
	Config      config.ServerConfig
	CoreService core.CoreServiceClient
	Logger      *slog.Logger
}

type Server struct {
	server      http.Server
	coreService core.CoreServiceClient
	logger      *slog.Logger
}

func New(cfg ServerConfig) *Server {
	s := new(Server)
	s.server.Addr = fmt.Sprintf(":%d", cfg.Config.Port)
	s.server.Handler = s.setupRouter()
	s.coreService = cfg.CoreService
	s.logger = cfg.Logger

	return s
}

func (s *Server) setupRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /core/ping", s.pingHandler)
	s.server.Handler = s.loggingMiddleware(mux)
	return mux
}

func (s *Server) Start() error {
	s.logger.Debug("server started", slog.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}
