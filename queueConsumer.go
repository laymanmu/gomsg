package gomsg

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
	"github.com/twinj/uuid"
)

// OnMessageCallback is called once for each message received.
type OnMessageCallback func(message *Message) error

// QueueConsumer consumes a given queue.
type QueueConsumer struct {
	Name        string
	QueueConfig *QueueConfig
	Error       error
	startOnce   sync.Once
	stopOnce    sync.Once
	stopChan    chan (bool)
}

// NewQueueConsumer creates a QueueConsumer.
func NewQueueConsumer(queueConfig *QueueConfig) *QueueConsumer {
	return &QueueConsumer{Name: uuid.NewV4().String(), QueueConfig: queueConfig, stopChan: make(chan (bool))}
}

// Start will connect to the queue & start consuming messages.
func (qc *QueueConsumer) Start(callback OnMessageCallback) {
	qc.startOnce.Do(func() {
		qc.consume(callback)
	})
}

// Stop will stop the consumer & close its connection to the queue.
func (qc *QueueConsumer) Stop() {
	qc.stopOnce.Do(func() {
		qc.stopChan <- true
	})
}

func (qc *QueueConsumer) consume(callback OnMessageCallback) {
	log.Printf("[consumer:%v] started consuming", qc.Name)

	config := qc.QueueConfig
	connection, err := amqp.Dial(config.ConnectionString())
	if err != nil {
		qc.Error = fmt.Errorf("[consumer:%v] failed to connect to RabbitMQ server.  error: %v", qc.Name, err)
		return
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		qc.Error = fmt.Errorf("[consumer:%v] failed to open a channel: %v", qc.Name, err)
		return
	}
	defer channel.Close()

	queue, err := channel.QueueDeclare(config.Name, config.Durable, config.AutoDelete, config.Exclusive, config.NoWait, config.Arguments)
	if err != nil {
		qc.Error = fmt.Errorf("[consumer:%v] failed to declare a queue: %v", qc.Name, err)
		return
	}

	// create & consume a message chan from amqp:
	messageChan, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		qc.Error = fmt.Errorf("[consumer:%v] failed to create the messageChan: %v", qc.Name, err)
		return
	}

	go func() {
		for delivery := range messageChan {
			log.Printf("[consumer:%v] received delivery with correlationId: %v", qc.Name, delivery.CorrelationId)
			err := callback(NewMessage(delivery))
			if err != nil {
				log.Printf("[consumer:%v] stopping per caught an error from the callback: %v", qc.Name, err)
				qc.Error = err
				qc.Stop()
			}
		}
	}()

	// block until Shutdown() is called which sends to this channel:
	<-qc.stopChan
	close(qc.stopChan)
	log.Printf("[consumer:%v] stopped consuming", qc.Name)
}
