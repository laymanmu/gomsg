package gomsg

import (
	"fmt"

	"github.com/streadway/amqp"
)

// QueueConfig holds configuration for connecting to a rabbitMQ queue.
type QueueConfig struct {
	Name       string
	UserName   string
	Password   string
	HostName   string
	Port       int
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Arguments  amqp.Table
}

// NewQueueConfig creates a QueueConfig instance.
func NewQueueConfig(name string) *QueueConfig {
	return &QueueConfig{
		Name:       name,
		UserName:   "guest",
		Password:   "guest",
		HostName:   "localhost",
		Port:       5672,
		Durable:    false,
		AutoDelete: false,
		Exclusive:  false,
		NoWait:     false}
}

// ConnectionString returns the connection string for amqp.
func (q *QueueConfig) ConnectionString() string {
	return fmt.Sprintf("amqp://%v:%v@%v:%v", q.UserName, q.Password, q.HostName, q.Port)
}
