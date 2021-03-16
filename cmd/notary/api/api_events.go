package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	"vax/pkg/httputils"
)

type createEventRequest struct {
	Hash string `json:"hash"`
}

func (s *Server) GetEvent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	eventId := strings.TrimSpace(ps.ByName("event_id"))
	if eventId == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event := s.eventsDao.Get(eventId); event != nil {
		httputils.WriteResponseObject(w, r, event)
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func (s *Server) PostEvent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req createEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Hash == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if event := s.eventsDao.Create(req.Hash); event != nil {
		httputils.WriteResponseObject(w, r, event)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
