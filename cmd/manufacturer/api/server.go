package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"vax/pkg/memory"
	"vax/pkg/model/dao"
)

const (
	shipmentPath = "/shipment"
)

type Server struct {
	actorId      string
	apiBasePath  string
	shipmentsDao dao.Shipment
	router       *httprouter.Router
}

func NewServer(ctx context.Context, actorId string, apiBasePath string) *Server {
	s := &Server{
		actorId:      actorId,
		apiBasePath:  apiBasePath,
		shipmentsDao: memory.NewShipmentsDao(),
	}

	restApi := httprouter.New()

	restApi.POST(apiBasePath+shipmentPath, s.PostShipment)

	s.router = restApi

	return s
}

// ServeHTTP is the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
