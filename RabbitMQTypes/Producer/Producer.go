// producer.go

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

const (
	Category1 = "celular"
	Category2 = "lavadora"
	Category3 = "refrigerador"
	Category4 = "tablet"
	Category5 = "television"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Data struct {
	Timestamp int64  `json:"timestamp"`
	Value     string `json:"value"`
}

type Device struct {
	ID   int
	Type string // Add the device type
}

func (d *Device) SendData(ch *amqp.Channel) {
	categories := []string{Category1, Category2, Category3, Category4, Category5}
	deviceTypeIndex := 0
	for i, category := range categories {
		if d.Type == category {
			deviceTypeIndex = i + 1
			break
		}
	}

	data := RandString(deviceTypeIndex)
	payload := Data{Value: data, Timestamp: time.Now().UnixNano()}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ch.Publish(
		"",
		d.Type, // use device type as the routing key
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonPayload,
		})

	if err != nil {
		fmt.Println(err)
	}
	//fmt.Printf("Device %d Sending: %s\n", d.ID, string(jsonPayload))
}

func main() {
	// Number of devices
	deviceCount := 20000

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

	categories := []string{Category1, Category2, Category3, Category4, Category5}

	for i := 0; i < deviceCount; i++ {
		go func(id int) {
			for {
				device := Device{ID: id, Type: categories[id%5]}
				device.SendData(ch)
				time.Sleep(time.Second / 10)
			}

		}(i)
	}

	select {}
}
