package main

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func main() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "deviceData",
	})

	consumerCount := 3

	for i := 0; i < consumerCount; i++ {
		go func(id int) {
			for {
				m, err := r.ReadMessage(context.Background())
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Printf("Consumer %d received: %s\n", id, string(m.Value))
			}
		}(i)
	}

	select {}
}
