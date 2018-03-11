# gomsg

### Example producer & consumer for RabbitMQ queue using amqp package.
```golang
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/laymanmu/gomsg"
	"github.com/twinj/uuid"
)

func main() {
	config := gomsg.NewQueueConfig("hello")

	// consumer:
	consumer := gomsg.NewQueueConsumer(config)
	callback := func(m *gomsg.Message) error {
		log.Printf("callback got a message: %v", m.JSON())
		log.Printf("  payload: %v", string(m.Payload))
		// returning an error will stop the consumer:
		return nil
	}
	go consumer.Start(callback)

	// producer:
	producer := gomsg.NewQueueProducer(config)
	go producer.Start()

	type TestPayload struct {
		Color  string `json:"color"`
		Weight int    `json:"weight"`
	}

	iterations := 3
	for i := 0; i < iterations; i++ {
		log.Printf("iteration %v", i)
		if consumer.Error != nil {
			log.Fatalf("consumer error: %v", consumer.Error)
		}
		if producer.Error != nil {
			log.Fatalf("producer error: %v", producer.Error)
		}

		jsonPayload, _ := json.Marshal(TestPayload{Color: "red", Weight: 44})

		m := gomsg.CreateMessage(
			uuid.NewV4().String(),
			"app1",
			"TestMessage",
			"application/json",
			jsonPayload,
		)
		producer.Send(*m)
		time.Sleep(1 * time.Second)
		fmt.Println()
	}

	producer.Stop()
	consumer.Stop()
	time.Sleep(1 * time.Second)
}
```
Output:
```
2018/03/10 23:45:00 iteration 0
2018/03/10 23:45:00 [consumer:912480d5-853a-450c-8203-7d1d9cc93e02] started consuming
2018/03/10 23:45:00 [producer:083a8f4c-9edb-4eff-acf5-1b54ba249842] started producing
2018/03/10 23:45:00 [producer:083a8f4c-9edb-4eff-acf5-1b54ba249842] produced a message with correlationId: 9e79ab9d-5131-49cd-9a83-a5a1da2dd0f0
2018/03/10 23:45:00 [consumer:912480d5-853a-450c-8203-7d1d9cc93e02] consumed a message with correlationId: 9e79ab9d-5131-49cd-9a83-a5a1da2dd0f0
2018/03/10 23:45:00 callback got a message: {"ID":"7de38c03-4a63-4b98-8f78-c2453cd2d2c8","Type":"TestMessage","Payload":"eyJjb2xvciI6InJlZCIsIndlaWdodCI6NDR9","ContentType":"application/json","CorrelationID":"9e79ab9d-5131-49cd-9a83-a5a1da2dd0f0","AppID":"app1"}
2018/03/10 23:45:00   payload: {"color":"red","weight":44}

2018/03/10 23:45:01 iteration 1
2018/03/10 23:45:01 [producer:083a8f4c-9edb-4eff-acf5-1b54ba249842] produced a message with correlationId: bd572bd6-c941-462c-be9f-c9a0cab48a8c
2018/03/10 23:45:01 [consumer:912480d5-853a-450c-8203-7d1d9cc93e02] consumed a message with correlationId: bd572bd6-c941-462c-be9f-c9a0cab48a8c
2018/03/10 23:45:01 callback got a message: {"ID":"55cea3f6-a67c-4946-b2d1-4718876720cd","Type":"TestMessage","Payload":"eyJjb2xvciI6InJlZCIsIndlaWdodCI6NDR9","ContentType":"application/json","CorrelationID":"bd572bd6-c941-462c-be9f-c9a0cab48a8c","AppID":"app1"}
2018/03/10 23:45:01   payload: {"color":"red","weight":44}

2018/03/10 23:45:02 iteration 2
2018/03/10 23:45:02 [producer:083a8f4c-9edb-4eff-acf5-1b54ba249842] produced a message with correlationId: 36720063-657e-4a00-9099-d5052c2997fe
2018/03/10 23:45:02 [consumer:912480d5-853a-450c-8203-7d1d9cc93e02] consumed a message with correlationId: 36720063-657e-4a00-9099-d5052c2997fe
2018/03/10 23:45:02 callback got a message: {"ID":"669e05bc-8bca-47e8-a092-e1878c094545","Type":"TestMessage","Payload":"eyJjb2xvciI6InJlZCIsIndlaWdodCI6NDR9","ContentType":"application/json","CorrelationID":"36720063-657e-4a00-9099-d5052c2997fe","AppID":"app1"}
2018/03/10 23:45:02   payload: {"color":"red","weight":44}

2018/03/10 23:45:03 [producer:083a8f4c-9edb-4eff-acf5-1b54ba249842] stopped producing
2018/03/10 23:45:03 [consumer:912480d5-853a-450c-8203-7d1d9cc93e02] stopped consuming
```