package notary

import (
	"vax/pkg/hash"
	"vax/pkg/model/dao"
	"vax/pkg/model/dto"
)

func VerifyEvent(event dto.Event, localDao dao.Shipment, canHaveLocalCopy bool) bool {
	// TODO: verify the signature
	if event.Hash == "" || event.Id == "" {
		return false
	} else if localEvent := localDao.Get(event.Payload.Shipment.Id); localEvent != nil {
		return canHaveLocalCopy
	} else if payloadHash, err := hash.EventPayload(*event.Payload); err != nil || event.Hash != payloadHash {
		return false
	} else if notarizedEvent, err := GetEvent(event.Id); err != nil {
		return false
	} else if notarizedEvent.Nonce == event.Nonce && notarizedEvent.Timestamp == event.Timestamp {
		return true
	} else {
		return false
	}
}
