package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"vax/pkg/model/dto"
	"vax/pkg/queue"
)

type redisQueue struct {
	channelName string
	pubSub      *redis.PubSub
}

func NewQueue(ctx context.Context, actorId string, shipmentId string) (queue.Queue, error) {
	channelName := getChannelNameForActor(actorId, shipmentId)
	pubSub := client.Subscribe(ctx, channelName)

	// Wait for the subscription confirmation.
	// TODO: this first message can something different than ACK, if the channel was already active
	if _, err := pubSub.Receive(ctx); err != nil {
		return nil, err
	}

	return &redisQueue{
		channelName: channelName,
		pubSub:      pubSub,
	}, nil
}

func (q *redisQueue) SendMessage(ctx context.Context, message dto.Message, actorId string, shipmentId string) error {
	return q.rawSend(ctx, actorId, shipmentId, message)
}

func (q *redisQueue) SendEvent(ctx context.Context, event dto.Event, actorId string, shipmentId string) error {
	return q.rawSend(ctx, actorId, shipmentId, event)
}

func (q *redisQueue) Listen(ctx context.Context) (<-chan dto.Message, <-chan dto.Event) {
	redisChannel := q.pubSub.Channel()
	messagesChannel := make(chan dto.Message, 10)
	eventsChannel := make(chan dto.Event, 10)

	go func() {
		defer func() {
			close(eventsChannel)
			close(messagesChannel)
		}()

		var tmpMsg dto.Message
		var tmpEvent dto.Event

		for rawMessage := range redisChannel {
			if err := json.Unmarshal([]byte(rawMessage.Payload), &tmpMsg); err != nil {
				if err := json.Unmarshal([]byte(rawMessage.Payload), &tmpEvent); err != nil {
					log.Printf("problem unmarshaling redis message")
					continue
				} else {
					eventsChannel <- tmpEvent
				}
			} else {
				messagesChannel <- tmpMsg
			}
		}
	}()

	return messagesChannel, eventsChannel
}

func (q *redisQueue) Close() error {
	return q.pubSub.Close()
}

func (q *redisQueue) rawSend(ctx context.Context, actorId string, shipmentId string, payload interface{}) error {
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = client.Publish(ctx, getChannelNameForActor(actorId, shipmentId), string(bytes)).Result()
	return err
}

func getChannelNameForActor(actorId string, shipmentId string) string {
	name := "actor:" + actorId
	if shipmentId != "" {
		name += ":" + shipmentId
	}
	return name
}
