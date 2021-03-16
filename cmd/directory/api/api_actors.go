package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

func (s *Server) GetActorByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	actorName := strings.TrimSpace(ps.ByName("actor_name"))
	if actorName == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO
}

func (s *Server) GetActorById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	actorId := strings.TrimSpace(ps.ByName("actor_id"))
	if actorId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// TODO
}
