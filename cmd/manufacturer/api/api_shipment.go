package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"vax/pkg/hash"
	"vax/pkg/httputils"
	"vax/pkg/model/dto"
	"vax/pkg/notary"
	"vax/pkg/queue"
	"vax/pkg/redis"
)

type createShipmentRequest struct {
	VaccineName    string `json:"vaccineName"`
	Quantity       uint64 `json:"quantity"`
	ExpirationDays int64  `json:"expirationDays"`
	AuthorityId    string `json:"authorityId"`
	CustomerId     string `json:"customerId"`
}

var (
	ErrQueueClosed      = errors.New("queue closed")
	ErrEventNotVerified = errors.New("event not verified")
)

func (s *Server) PostShipment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req createShipmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shipment := s.shipmentsDao.Create(req.VaccineName, req.Quantity, req.ExpirationDays, s.actorId, req.AuthorityId, req.CustomerId)

	ctx := r.Context()
	pubSubQueue, err := redis.NewQueue(ctx, s.actorId, shipment.Id)
	if err != nil {
		panic(err)
	}
	defer pubSubQueue.Close()

	_, eventsChannel := pubSubQueue.Listen(ctx)

	if err = s.sendEventToActor(ctx, &shipment, dto.EventTypeShipmentCreated, req.AuthorityId, pubSubQueue, eventsChannel); err != nil {
		panic(err)
	}

	if err = s.sendEventToActor(ctx, &shipment, dto.EventTypeShipmentShipped, req.CustomerId, pubSubQueue, eventsChannel); err != nil {
		panic(err)
	}

	httputils.WriteResponseObject(w, r, shipment)
}

func (s *Server) sendEventToActor(ctx context.Context, shipment *dto.Shipment, eventType dto.EventType, authorityId string, pubSubQueue queue.Queue, eventsResponsesChannel <-chan dto.Event) error {
	eventPayload := dto.EventPayload{
		Type:     eventType,
		ActorId:  s.actorId,
		Shipment: *shipment,
	}

	payloadHash, err := hash.EventPayload(eventPayload)
	if err != nil {
		return err
	}

	event, err := notary.CreateEvent(payloadHash)
	if err != nil {
		return err
	}

	event.Payload = &eventPayload
	// TODO: event.Signature

	shipment.Events = append(shipment.Events, *event)
	s.shipmentsDao.Update(*shipment)

	// Do not use shipment ID because other actors have global queues and only
	// the manufacturer needs a "shipment queue" so we can receive the response
	// directly back to the HTTP handler.
	if err = pubSubQueue.SendEvent(ctx, *event, authorityId, ""); err != nil {
		return err
	}

	// wait for the response
	eventMsg, open := <-eventsResponsesChannel
	if !open {
		return ErrQueueClosed
	}

	if notary.VerifyEvent(eventMsg, s.shipmentsDao, true) {
		shipment.Events = append(shipment.Events, eventMsg)
		s.shipmentsDao.Update(*shipment)
	} else {
		return ErrEventNotVerified
	}

	return nil
}
