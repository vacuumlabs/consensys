package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"vax/pkg/memory"
	"vax/pkg/model/dao"
)

const (
	eventsPath = "/event"
)

type Server struct {
	apiBasePath   string
	nonceProvider dao.NonceProvider
	eventsDao     dao.Event
	router        *httprouter.Router
}

func NewServer(ctx context.Context, apiBasePath string) *Server {
	nonceProvider := memory.NewMonotonicNonceProvider()

	s := &Server{
		apiBasePath:   apiBasePath,
		nonceProvider: nonceProvider,
		eventsDao:     memory.NewEventsDao(nonceProvider),
	}

	restApi := httprouter.New()

	restApi.GET(apiBasePath+eventsPath+"/:event_id", s.GetEvent)
	restApi.POST(apiBasePath+eventsPath, s.PostEvent)

	s.router = restApi

	return s
}

// ServeHTTP is the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
