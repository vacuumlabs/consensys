package queue

import (
	"context"
	"vax/pkg/model/dto"
)

type Queue interface {
	SendMessage(ctx context.Context, message dto.Message, actorId string, shipmentId string) error
	SendEvent(ctx context.Context, event dto.Event, actorId string, shipmentId string) error
	Listen(ctx context.Context) (<-chan dto.Message, <-chan dto.Event)
	Close() error
}
