package server

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"sre.qlik.com/palindrome/logger"
)

const (
	tracingID = "Request-Tracing-ID"
)

// RequestTracing is type defined to be used for a context with value
type RequestTracing string

// RegisterRoutes register the endpoints for the service to receive requests on
func (s *server) RegisterRoutes() {
	rootRouter := s.router
	rootRouter.HandleFunc("/health", s.healthHandler()).Methods(http.MethodGet)

	// register routes for the api subrouter (i.e., endpoints prefixed with /api/v1)
	apiRouter := rootRouter.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/messages", s.handleGetMessages()).Methods(http.MethodGet)
	apiRouter.HandleFunc("/messages", s.handlePostMessage()).Methods(http.MethodPost)
	apiRouter.HandleFunc("/messages/{id}", s.handleGetSingleMessage()).Methods(http.MethodGet)
	apiRouter.HandleFunc("/messages/{id}", s.handleDeleteMessage()).Methods(http.MethodDelete)
}

// Logging middleware logs all the incoming requests
func Logging(l logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t := time.Now().UTC()
			defer func() {
				requestID, ok := r.Context().Value(RequestTracing(tracingID)).(string)
				if !ok {
					requestID = "unknown"
				}
				l.Info("%s: %s  Method: %s URL: %s RemoteAddr: %s UserAgent: %s Latency: %v ", tracingID, requestID, r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent(), time.Since(t))
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// Tracing middleware adds a TracingID to each request
func Tracing() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-Id")
			if requestID == "" {
				requestID = uuid.New().String()
			}
			ctx := context.WithValue(r.Context(), RequestTracing(tracingID), requestID)
			w.Header().Set("X-Request-Id", requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
