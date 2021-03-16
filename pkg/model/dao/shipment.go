package dao

import "vax/pkg/model/dto"

type Shipment interface {
	Create(vaccineName string, quantity uint64, expirationDays int64, manufacturerId string, authorityId string, customerId string) dto.Shipment
	Get(id string) *dto.Shipment
	Update(shipment dto.Shipment)
}
