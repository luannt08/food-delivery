package api

import (
	"net/http"
	"time"
)

type AuthServer struct {
	serverMux *http.ServeMux
}

func NewAuthServer() *AuthServer {
	mux := http.NewServeMux()
	return &AuthServer{
		serverMux: mux,
	}
}

func (s *AuthServer) RegisterHandler(path string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	s.serverMux.HandleFunc(path, handlerFunc)
}

func (s *AuthServer) Start(addr string) error {
	server := &http.Server{
		Addr:         addr,
		Handler:      s.serverMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server.ListenAndServe()
}
