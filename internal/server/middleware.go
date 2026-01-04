package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info("request started",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.Any("date", time.Now().UTC()))

		next.ServeHTTP(w, r)

		s.logger.Info("request completed",
			slog.String("method", r.Method),
			slog.String("url", r.URL.Path),
			slog.Any("date", time.Now().UTC()))
	})
}

func (s *Server) extractIdsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userId := r.Header.Get("X-User-ID")
        dormitoryId := r.Header.Get("X-Dormitory-ID")
        
        ctx := r.Context()
        if userId != "" {
            ctx = context.WithValue(ctx, "userId", userId)
        }
        if dormitoryId != "" {
            ctx = context.WithValue(ctx, "dormitoryId", dormitoryId)
        }
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}