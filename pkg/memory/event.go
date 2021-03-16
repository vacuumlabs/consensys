package memory

import (
	"github.com/google/uuid"
	"sync"
	"time"
	"vax/pkg/model/dao"
	"vax/pkg/model/dto"
)

var contains = struct{}{}

type eventsDao struct {
	data   map[string]dto.Event
	hashes map[string]struct{}
	mu     sync.Mutex

	nonceProvider dao.NonceProvider
}

func NewEventsDao(nonceProvider dao.NonceProvider) dao.Event {
	return &eventsDao{
		data:          make(map[string]dto.Event, 1000),
		hashes:        make(map[string]struct{}, 1000),
		nonceProvider: nonceProvider,
	}
}

func (dao *eventsDao) Get(id string) *dto.Event {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	if event, contains := dao.data[id]; contains {
		return &event
	}

	return nil
}

func (dao *eventsDao) Create(hash string) *dto.Event {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	if _, contains := dao.hashes[hash]; contains {
		return nil
	}

	event := dto.Event{
		Id:        uuid.NewString(),
		Nonce:     dao.nonceProvider.Provide(),
		Timestamp: time.Now().UTC().Unix(),
		Payload:   nil, // TODO
		Hash:      hash,
		Signature: "", // TODO
	}

	dao.data[event.Id] = event
	dao.hashes[hash] = contains

	return &event
}
