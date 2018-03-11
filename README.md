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
2018/03/10 23:20:58 iteration 0
2018/03/10 23:20:58 [producer:c977d8e9-257c-4825-a334-60cea016e7d9] started producing
2018/03/10 23:20:58 [consumer:c382aae7-aab7-42fa-ae55-eb1805468bb5] started consuming
2018/03/10 23:20:58 [producer:c977d8e9-257c-4825-a334-60cea016e7d9] produced a message with correlationId: 91b8589f-79c8-4fee-811a-8644f5413f35
2018/03/10 23:20:58 [consumer:c382aae7-aab7-42fa-ae55-eb1805468bb5] received delivery with correlationId: 91b8589f-79c8-4fee-811a-8644f5413f35
2018/03/10 23:20:58 callback got a message: {"ID":"34fab67d-3ab7-48ef-912a-990646169d93","Type":"TestMessage","Payload":"eyJjb2xvciI6InJlZCIsIndlaWdodCI6NDR9","ContentType":"application/json","CorrelationID":"91b8589f-79c8-4fee-811a-8644f5413f35","AppID":"app1"}
2018/03/10 23:20:58   payload: {"color":"red","weight":44}

2018/03/10 23:20:59 iteration 1
2018/03/10 23:20:59 [producer:c977d8e9-257c-4825-a334-60cea016e7d9] produced a message with correlationId: 62fcfb61-cf4e-4708-854a-47d72811d789
2018/03/10 23:20:59 [consumer:c382aae7-aab7-42fa-ae55-eb1805468bb5] received delivery with correlationId: 62fcfb61-cf4e-4708-854a-47d72811d789
2018/03/10 23:20:59 callback got a message: {"ID":"452deab5-a0d8-4f31-adb9-be6700dafc40","Type":"TestMessage","Payload":"eyJjb2xvciI6InJlZCIsIndlaWdodCI6NDR9","ContentType":"application/json","CorrelationID":"62fcfb61-cf4e-4708-854a-47d72811d789","AppID":"app1"}
2018/03/10 23:20:59   payload: {"color":"red","weight":44}

2018/03/10 23:21:00 iteration 2
2018/03/10 23:21:00 [producer:c977d8e9-257c-4825-a334-60cea016e7d9] produced a message with correlationId: 3d07cfa9-4ff3-46eb-8111-b8d2d33a6a7c
2018/03/10 23:21:00 [consumer:c382aae7-aab7-42fa-ae55-eb1805468bb5] received delivery with correlationId: 3d07cfa9-4ff3-46eb-8111-b8d2d33a6a7c
2018/03/10 23:21:00 callback got a message: {"ID":"d24d7a31-8de1-4376-8dbd-7130eddb7014","Type":"TestMessage","Payload":"eyJjb2xvciI6InJlZCIsIndlaWdodCI6NDR9","ContentType":"application/json","CorrelationID":"3d07cfa9-4ff3-46eb-8111-b8d2d33a6a7c","AppID":"app1"}
2018/03/10 23:21:00   payload: {"color":"red","weight":44}

2018/03/10 23:21:01 [consumer:c382aae7-aab7-42fa-ae55-eb1805468bb5] stopped consuming
2018/03/10 23:21:01 [producer:c977d8e9-257c-4825-a334-60cea016e7d9] stopped producing
```