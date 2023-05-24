// consumer.go

package main

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/streadway/amqp"
)

type Message struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"deviceData",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}
	var messageCount int64
	var errorCount int64
	var latencytotal int64

	consumerCount := 50

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			currentMessageCount := atomic.LoadInt64(&messageCount)
			currentErrorCount := atomic.LoadInt64(&errorCount)
			var currentavgLatency float64

			if currentMessageCount != 0 {
				currentavgLatency = float64(atomic.LoadInt64(&latencytotal)) / float64(currentMessageCount)
			}

			fmt.Printf("Current latency: %f nanoseconds\n", currentavgLatency)
			fmt.Printf("Current throughput: %d messages per second\n", currentMessageCount)
			fmt.Printf("Current error rate: %d errors per second\n", currentErrorCount)
			atomic.StoreInt64(&latencytotal, 0)
			atomic.StoreInt64(&messageCount, 0)
			atomic.StoreInt64(&errorCount, 0)
		}
	}()

	for i := 0; i < consumerCount; i++ {

		go func(id int) {
			for d := range msgs {
				var msg Message
				err := json.Unmarshal(d.Body, &msg)
				if err != nil {
					fmt.Println("Error unmarshaling JSON:", err)
					atomic.AddInt64(&errorCount, 1)
					continue
				}
				atomic.AddInt64(&messageCount, 1)

				receivedTime := time.Now()
				latency := receivedTime.Sub(time.Unix(0, msg.Timestamp))
				atomic.AddInt64(&latencytotal, latency.Nanoseconds())

				//fmt.Printf("Consumer %d received: %s\n", id, d.Body)

				//fmt.Printf("Latency: %s\n", latency)

			}
		}(i)

	}

	select {}
}
