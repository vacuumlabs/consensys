package dto

import "errors"

type MessageType string

const (
	MessageTypeEventNotification      = "EventNotification"
	MessageTypeMessageAcknowledgement = "MessageAcknowledgement"
)

var (
	ErrUnknownMessageType = errors.New("unknown message type")
)

func (m MessageType) MarshalText() (text []byte, err error) {
	return []byte(m), nil
}

func (m *MessageType) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	switch string(text) {
	case MessageTypeEventNotification:
		*m = MessageTypeEventNotification
	case MessageTypeMessageAcknowledgement:
		*m = MessageTypeMessageAcknowledgement

	default:
		return ErrUnknownMessageType
	}

	return nil
}

type Message struct {
	Id        string         `json:"id"`
	Nonce     uint64         `json:"nonce"`
	Timestamp int64          `json:"timestamp"`
	Payload   MessagePayload `json:"payload"`
	Hash      string         `json:"hash"`
	Signature string         `json:"signature"`
}

type MessagePayload struct {
	Type           MessageType `json:"MessageType"`
	ActorId        string      `json:"ActorId"`
	MessageContent string      `json:"MessageContent"`
	Shipment       Shipment    `json:"shipment"`
}
