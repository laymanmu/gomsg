# gomsg

### Example RabbitMQ queue consumer:
```golang
package main

import (
	"log"
	"time"

	"github.com/laymanmu/gomsg"
)

func main() {
	config := gomsg.NewQueueConfig("hello")
	consumer := gomsg.NewQueueConsumer(config)

	callback := func(m *gomsg.Message) error {
		log.Printf("callback got a message: %v", m)
		// returning an error will stop the consumer:
		return nil
	}

	go consumer.Start(callback)

	iterations := 3
	for i := 0; i < iterations; i++ {
		log.Printf("iteration %v", i)
		if consumer.Error != nil {
			log.Fatalf("%v", consumer.Error)
		}
		time.Sleep(3 * time.Second)
	}

	consumer.Stop()
	time.Sleep(1 * time.Second)
}
```
