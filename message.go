package gomsg

import (
	"encoding/json"

	"github.com/streadway/amqp"
	"github.com/twinj/uuid"
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
		ContentType:   d.ContentType,
	}
}

// CreateMessage creates a message from the given parms.
func CreateMessage(correlationID string, appID string, messageType string, contentType string, payload []byte) *Message {
	return &Message{
		ID:            uuid.NewV4().String(),
		CorrelationID: correlationID,
		AppID:         appID,
		Type:          messageType,
		ContentType:   contentType,
		Payload:       payload,
	}
}

// JSON returns a json string for this message.
func (m *Message) JSON() string {
	j, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(j)
}
