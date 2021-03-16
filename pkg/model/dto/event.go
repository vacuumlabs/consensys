package dto

import "errors"

type EventType string

const (
	EventTypeShipmentCreated  = "ShipmentCreated"
	EventTypeShipmentShipped  = "ShipmentShipped"
	EventTypeShipmentApproved = "ShipmentApproved"
	EventTypeShipmentBlocked  = "ShipmentBlocked"
	EventTypeShipmentAccepted = "ShipmentAccepted"
	EventTypeShipmentRejected = "ShipmentRejected"
)

var (
	ErrUnknownEventType = errors.New("unknown message type")
)

func (m EventType) MarshalText() (text []byte, err error) {
	return []byte(m), nil
}

func (m *EventType) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	switch string(text) {
	case EventTypeShipmentCreated:
		*m = EventTypeShipmentCreated
	case EventTypeShipmentShipped:
		*m = EventTypeShipmentShipped
	case EventTypeShipmentApproved:
		*m = EventTypeShipmentApproved
	case EventTypeShipmentBlocked:
		*m = EventTypeShipmentBlocked
	case EventTypeShipmentAccepted:
		*m = EventTypeShipmentAccepted
	case EventTypeShipmentRejected:
		*m = EventTypeShipmentRejected

	default:
		return ErrUnknownEventType
	}

	return nil
}

type Event struct {
	Id        string        `json:"id"`
	Nonce     uint64        `json:"nonce"`
	Timestamp int64         `json:"timestamp"`
	Payload   *EventPayload `json:"payload,omitempty"`
	Hash      string        `json:"hash"`
	Signature string        `json:"signature"`
}

type EventPayload struct {
	Type     EventType `json:"MessageType"`
	ActorId  string    `json:"ActorId"`
	Shipment Shipment  `json:"shipment"`
}
