package server

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/dormitory-life/core/internal/config"
	core "github.com/dormitory-life/core/internal/service"
	"github.com/gorilla/mux"
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

func (s *Server) setupRouter() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/core/ping", s.pingHandler).Methods("GET")

	router.HandleFunc("/core/dormitories/grades", s.getDormitoriesAvgGradesHandler).Methods("GET")
	router.HandleFunc("/core/dormitories/{dormitory_id}/grades", s.getDormitoryAvgGradesHandler).Methods("GET")
	router.HandleFunc("/core/dormitories/{dormitory_id}/grades", s.createDormitoryGradeHandler).Methods("POST")

	router.HandleFunc("/core/dormitories", s.getDormitoriesHandler).Methods("GET")
	router.HandleFunc("/core/dormitories/{dormitory_id}", s.getDormitoryByIdHandler).Methods("GET")
	router.HandleFunc("/core/dormitories", s.createDormitoryHandler).Methods("POST")
	router.HandleFunc("/core/dormitories/{dormitory_id}", s.updateDormitoryHandler).Methods("PUT")
	router.HandleFunc("/core/dormitories/{dormitory_id}", s.deleteDormitoryHandler).Methods("DELETE")

	router.HandleFunc("/core/dormitories/{dormitory_id}/photos", s.createDormitoryPhotosHandler).Methods("POST")
	router.HandleFunc("/core/dormitories/{dormitory_id}/photos", s.deleteDormitoryPhotosHandler).Methods("DELETE")

	router.HandleFunc("/core/dormitories/support", s.createSupportRequestHandler).Methods("POST")

	router.HandleFunc("/core/dormitories/{dormitory_id}/reviews", s.getReviewsHandler).Methods("GET")
	router.HandleFunc("/core/dormitories/{dormitory_id}/reviews", s.createReviewHandler).Methods("POST")
	router.HandleFunc("/core/dormitories/{dormitory_id}/reviews/{review_id}", s.deleteReviewHandler).Methods("DELETE")

	return s.loggingMiddleware(s.extractIdsMiddleware(router))
}

func (s *Server) Start() error {
	s.logger.Debug("server started", slog.String("address", s.server.Addr))
	return s.server.ListenAndServe()
}
