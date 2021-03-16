package memory

import (
	"github.com/google/uuid"
	"sync"
	"time"
	"vax/pkg/model/dao"
	"vax/pkg/model/dto"
)

type shipmentsDao struct {
	data map[string]dto.Shipment
	mu   sync.Mutex
}

func NewShipmentsDao() dao.Shipment {
	return &shipmentsDao{
		data: make(map[string]dto.Shipment),
	}
}

func (dao *shipmentsDao) Create(vaccineName string, quantity uint64, expirationDays int64, manufacturerId string, authorityId string, customerId string) dto.Shipment {
	now := time.Now().UTC()
	shipment := dto.Shipment{
		Id:                uuid.NewString(),
		VaccineName:       vaccineName,
		Quantity:          quantity,
		ManufacturingDate: now.Unix(),
		ManufacturerId:    manufacturerId,
		ExpirationDate:    now.Add(time.Duration(expirationDays*24) * time.Hour).Unix(),
		AuthorityId:       authorityId,
		CustomerId:        customerId,
	}

	dao.mu.Lock()
	defer dao.mu.Unlock()

	if _, contains := dao.data[shipment.Id]; contains {
		// TODO
		panic("handle collision")
	}

	dao.data[shipment.Id] = shipment
	return shipment
}

func (dao *shipmentsDao) Get(id string) *dto.Shipment {
	dao.mu.Lock()
	defer dao.mu.Unlock()

	if shipment, contains := dao.data[id]; contains {
		return &shipment
	}

	return nil
}

func (dao *shipmentsDao) Update(shipment dto.Shipment) {
	if shipment.Id == "" {
		return
	}

	dao.mu.Lock()
	defer dao.mu.Unlock()

	dao.data[shipment.Id] = shipment
}
