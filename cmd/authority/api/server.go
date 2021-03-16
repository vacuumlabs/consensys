package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"vax/pkg/hash"
	"vax/pkg/memory"
	"vax/pkg/model/dao"
	"vax/pkg/model/dto"
	"vax/pkg/notary"
	"vax/pkg/queue"
	"vax/pkg/redis"
)

type Server struct {
	actorId      string
	apiBasePath  string
	shipmentsDao dao.Shipment
	router       *httprouter.Router
}

func NewServer(ctx context.Context, actorId string, apiBasePath string) *Server {
	pubSubQueue, err := redis.NewQueue(ctx, actorId, "")
	if err != nil {
		panic(err)
	}

	s := &Server{
		actorId:      actorId,
		apiBasePath:  apiBasePath,
		shipmentsDao: memory.NewShipmentsDao(),
	}

	restApi := httprouter.New()

	s.router = restApi

	go func() {
		msgChannel, eventChannel := pubSubQueue.Listen(ctx)

		for {
			select {
			case msg := <-msgChannel:
				log.Printf("message received: %+v", msg)
				// TODO

			case event := <-eventChannel:
				log.Println(event)
				go s.processEvent(context.Background(), event, pubSubQueue)
			}
		}
	}()

	return s
}

// ServeHTTP is the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) processEvent(ctx context.Context, event dto.Event, pubSubQueue queue.Queue) {
	eventPayload := dto.EventPayload{
		ActorId:  s.actorId,
		Shipment: event.Payload.Shipment,
	}

	if notary.VerifyEvent(event, s.shipmentsDao, false) {
		eventPayload.Type = dto.EventTypeShipmentApproved
	} else {
		eventPayload.Type = dto.EventTypeShipmentRejected
	}

	payloadHash, err := hash.EventPayload(eventPayload)
	if err != nil {
		panic(err)
	}

	responseEvent, err := notary.CreateEvent(payloadHash)
	if err != nil {
		panic(err)
	}

	responseEvent.Payload = &eventPayload
	// TODO: event.Signature

	event.Payload.Shipment.Events = append(event.Payload.Shipment.Events, *responseEvent)
	s.shipmentsDao.Update(event.Payload.Shipment)

	if err := pubSubQueue.SendEvent(ctx, *responseEvent, event.Payload.ActorId, event.Payload.Shipment.Id); err != nil {
		log.Printf("problem sending ack message: %+v", err)
	}
}
