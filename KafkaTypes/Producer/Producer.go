//producer.go

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
	ID         int
	DeviceType int
}

func (d *Device) SendData(writers []*kafka.Writer) {
	// generate random data of length d.DeviceType
	data := RandStringRunes(d.DeviceType)
	payload := Data{Value: data, Timestamp: float64(time.Now().UnixNano())}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := writers[d.DeviceType-1]

	err = w.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(fmt.Sprintf("Key-%d", d.ID)),
		Value: jsonPayload,
	})

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Device %d by %d Sending: %s\n", d.ID, d.DeviceType, string(jsonPayload))

}

func main() {
	// Number of devices
	deviceCount := 1000

	writers := make([]*kafka.Writer, 5)
	for i := 0; i < 5; i++ {
		writers[i] = kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{"kafka:9092"},
			Topic:   fmt.Sprintf("deviceData%d", i+1),
		})
	}

	for i := 0; i < deviceCount; i++ {
		go func(id int) {
			for {
				device := Device{ID: id, DeviceType: rand.Intn(5) + 1}
				device.SendData(writers)
				time.Sleep(time.Second / 10)
			}

		}(i)
	}

	select {}
}
