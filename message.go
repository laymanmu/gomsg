package gomsg

import (
	"github.com/streadway/amqp"
)

// Message model.
type Message struct {
	ID            string
	Type          string
	Payload       []byte
	ContentType   string
	CorrelationID string
	AppID         string
}

// NewMessage creates a Message.
func NewMessage(d amqp.Delivery) *Message {
	return &Message{
		ID:            d.MessageId,
		Type:          d.Type,
		Payload:       d.Body,
		CorrelationID: d.CorrelationId,
		AppID:         d.AppId,
		ContentType:   d.ContentType}
}
