package gomsg

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

// Message model.
type Message struct {
	ID            string
	Payload       []byte
	CorrelationID string
	AppID         string
	ContentType   string
}

// NewMessage creates a Message.
func NewMessage(d amqp.Delivery) *Message {
	return &Message{ID: d.MessageId, Payload: d.Body, CorrelationID: d.CorrelationId, AppID: d.AppId, ContentType: d.ContentType}
}

// JSON returns the message data as a json string.
func (m *Message) JSON() string {
	b, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
