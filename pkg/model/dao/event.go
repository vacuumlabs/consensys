package dao

import "vax/pkg/model/dto"

type Event interface {
	Create(hash string) *dto.Event
	Get(id string) *dto.Event
}
