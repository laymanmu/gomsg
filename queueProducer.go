package gomsg

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
	"github.com/twinj/uuid"
)

// QueueProducer produces messages in a queue.
type QueueProducer struct {
	Name        string
	QueueConfig *QueueConfig
	Error       error
	startOnce   sync.Once
	stopOnce    sync.Once
	stopChan    chan (bool)
	msgChan     chan (Message)
}

// NewQueueProducer creates a QueueProducer.
func NewQueueProducer(queueConfig *QueueConfig) *QueueProducer {
	return &QueueProducer{Name: uuid.NewV4().String(), QueueConfig: queueConfig, stopChan: make(chan (bool)), msgChan: make(chan (Message))}
}

// Start will connect to the queue & start producing messages.
func (qp *QueueProducer) Start() {
	qp.startOnce.Do(func() {
		qp.produce()
	})
}

// Stop will stop the producer & close its connection to the queue.
func (qp *QueueProducer) Stop() {
	qp.stopOnce.Do(func() {
		qp.stopChan <- true
	})
}

// Send will queue up a given Message to be sent.
func (qp *QueueProducer) Send(message Message) {
	qp.msgChan <- message
}

func (qp *QueueProducer) produce() {
	log.Printf("[producer:%v] started producing", qp.Name)

	config := qp.QueueConfig
	connection, err := amqp.Dial(config.ConnectionString())
	if err != nil {
		qp.Error = fmt.Errorf("[producer:%v] failed to connect to RabbitMQ server.  error: %v", qp.Name, err)
		return
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		qp.Error = fmt.Errorf("[producer:%v] failed to open a channel: %v", qp.Name, err)
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(config.Name, config.Durable, config.AutoDelete, config.Exclusive, config.NoWait, config.Arguments)
	if err != nil {
		qp.Error = fmt.Errorf("[producer:%v] failed to declare a queue: %v", qp.Name, err)
		return
	}

	go func() {
		for m := range qp.msgChan {
			err = channel.Publish(
				"",         // exchange
				queue.Name, // routing key
				false,      // mandatory
				false,      // immediate
				amqp.Publishing{
					MessageId:     m.ID,
					CorrelationId: m.CorrelationID,
					AppId:         m.AppID,
					Type:          m.Type,
					ContentType:   m.ContentType,
					Body:          []byte(m.Payload),
				})
			if err != nil {
				log.Printf("[producer:%v] stopping per caught an error from the publisher: %v", qp.Name, err)
				qp.Error = err
				qp.Stop()
			} else {
				log.Printf("[producer:%v] produced a message with correlationId: %v", qp.Name, m.CorrelationID)
			}
		}
	}()

	// block until Shutdown() is called which sends to this channel:
	<-qp.stopChan
	close(qp.stopChan)
	log.Printf("[producer:%v] stopped producing", qp.Name)
}
