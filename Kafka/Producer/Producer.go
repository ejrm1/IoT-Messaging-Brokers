package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type Data struct {
	Timestamp float64 `json:"timestamp"`
	Value     string  `json:"value"`
}

type Device struct {
	ID int
}

func (d *Device) SendData(Device_type int, w *kafka.Writer) {
	data := RandStringRunes(Device_type)
	payload := Data{Value: data, Timestamp: float64(time.Now().UnixNano())}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("Key-%d", d.ID)),
		Value: jsonPayload,
	})

	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("Device %d Sending: %s\n", d.ID, string(jsonPayload))

}

func main() {
	// Number of devices
	deviceCount := 10000

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "deviceData",
	})

	for i := 0; i < deviceCount; i++ {
		go func(id int) {
			for {
				device := Device{ID: id}
				device.SendData(10, w)
				time.Sleep(time.Second / 10)
			}

		}(i)
	}

	select {}
}
