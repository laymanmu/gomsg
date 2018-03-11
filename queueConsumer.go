package gomsg

import (
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
	"github.com/twinj/uuid"
)

// OnMessageCallback is called once for each message received.
type OnMessageCallback func(message *Message)

// QueueConsumer consumes a given queue.
type QueueConsumer struct {
	Name         string
	QueueConfig  *QueueConfig
	Error        error
	once         sync.Once
	shutdownChan chan (bool)
	isDead       bool
}

// NewQueueConsumer creates a QueueConsumer.
func NewQueueConsumer(queueConfig *QueueConfig) *QueueConsumer {
	return &QueueConsumer{Name: uuid.NewV4().String(), QueueConfig: queueConfig, shutdownChan: make(chan (bool)), isDead: false}
}

// Consume subscribes to a queue then blocks until Shutdown() is called.
func (qc *QueueConsumer) Consume(callback OnMessageCallback) {
	if qc.isDead {
		qc.Error = fmt.Errorf("[consumer:%v] Consume() called for dead consumer", qc.Name)
		return
	}

	log.Printf("[consumer:%v] consumer is alive", qc.Name)

	// create a rabbitMQ connection:
	config := qc.QueueConfig
	connection, err := amqp.Dial(config.ConnectionString())
	if err != nil {
		qc.Error = fmt.Errorf("[consumer:%v] failed to connect to RabbitMQ server using connection string: %v. error: %v", qc.Name, config.ConnectionString(), err)
		return
	}
	defer connection.Close()

	// create a rabbitMQ channel:
	channel, err := connection.Channel()
	if err != nil {
		qc.Error = fmt.Errorf("[consumer:%v] failed to open a channel: %v", qc.Name, err)
		return
	}
	defer channel.Close()

	// create a rabbitMQ queue:
	queue, err := channel.QueueDeclare(config.Name, config.Durable, config.AutoDelete, config.Exclusive, config.NoWait, config.Arguments)
	if err != nil {
		qc.Error = fmt.Errorf("[consumer:%v] failed to declare a queue: %v", qc.Name, err)
		return
	}

	// create a go messageChan for consuming:
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

	// start a go routine to consume to the messageChan:
	go func() {
		for delivery := range messageChan {
			log.Printf("[consumer:%v] received delivery: %v", qc.Name, delivery)
			callback(NewMessage(delivery))
		}
	}()

	// block until Shutdown() is called:
	<-qc.shutdownChan
}

// Shutdown will stop the consumer & close the connections to the queue.
func (qc *QueueConsumer) Shutdown() {
	qc.once.Do(func() {
		qc.shutdownChan <- true
		qc.isDead = true
		log.Printf("[consumer:%v] consumer is dead.", qc.Name)
	})
}
