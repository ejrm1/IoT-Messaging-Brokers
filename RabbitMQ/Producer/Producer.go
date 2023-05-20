// producer.go

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
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

func (d *Device) SendData(Device_type int, ch *amqp.Channel) {
	data := RandStringRunes(Device_type)
	payload := Data{Value: data, Timestamp: float64(time.Now().UnixNano())}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ch.Publish(
		"",
		"deviceData",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonPayload,
		})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Device %d Sending: %s\n", d.ID, string(jsonPayload))

}

func main() {
	// Number of devices
	deviceCount := 3

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

	_, err = ch.QueueDeclare(
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

	for i := 0; i < deviceCount; i++ {
		go func(id int) {
			for {
				device := Device{ID: id}
				device.SendData(10, ch)
				time.Sleep(time.Second)
			}

		}(i)
	}

	select {}
}