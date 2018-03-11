# gomessage

### Example RabbitMQ queue consumer:
```golang
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	config := NewQueueConfig("hello")
	consumer := NewQueueConsumer(config)

	go consumer.Consume(func(message *Message) {
		log.Printf("callback got a message: %v", message)
	})

	iterations := 3
	for i := 0; i < iterations; i++ {
		log.Printf("iteration %v", i)
		failOnError(consumer.Error, "consumer failed")
		time.Sleep(3 * time.Second)
	}

	consumer.Shutdown()
	time.Sleep(2 * time.Second)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
```
