package main

import (
	"context"
	"fmt"

	"sync/atomic"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "deviceData",
	})

	var messageCount int64
	var latencytotal int64
	consumerCount := 10000

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			currentMessageCount := atomic.LoadInt64(&messageCount)
			//currentErrorCount := atomic.LoadInt64(&errorCount)
			var currentavgLatency float64

			if currentMessageCount != 0 {
				currentavgLatency = float64(atomic.LoadInt64(&latencytotal)) / float64(currentMessageCount)
			}

			fmt.Printf("Current latency: %f nanoseconds\n", currentavgLatency)
			fmt.Printf("Current throughput: %d messages per second\n", currentMessageCount)
			//fmt.Printf("Current error rate: %d errors per second\n", currentErrorCount)
			atomic.StoreInt64(&latencytotal, 0)
			atomic.StoreInt64(&messageCount, 0)
			//atomic.StoreInt64(&errorCount, 0)
		}
	}()

	for i := 0; i < consumerCount; i++ {
		go func(id int) {
			for {
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					fmt.Println(err)
					continue
				}
				atomic.AddInt64(&messageCount, 1)
				latency := time.Since(m.Time)
				atomic.AddInt64(&latencytotal, latency.Nanoseconds())
				//fmt.Printf("Consumer %d received: %s (Latency:%s)\n", id, string(m.Value), latency.String())
			}
		}(i)
	}

	select {}
}
