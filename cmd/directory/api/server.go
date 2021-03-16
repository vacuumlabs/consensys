package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	actorPath   = "/actor"
	actorIdPath = "/actorId"
)

type Server struct {
	apiBasePath string
	router      *httprouter.Router
}

func NewServer(ctx context.Context, apiBasePath string) *Server {
	s := &Server{
		apiBasePath: apiBasePath,
	}

	restApi := httprouter.New()

	restApi.GET(apiBasePath+actorPath+"/:actor_name", s.GetActorByName)
	restApi.GET(apiBasePath+actorIdPath+"/:actor_id", s.GetActorById)

	s.router = restApi

	return s
}

// ServeHTTP is the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
