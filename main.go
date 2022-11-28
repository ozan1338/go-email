package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "pkc-ew3qg.asia-southeast2.gcp.confluent.cloud:9092",
		"security.protocol": "SASL_SSL",
		"sasl.username": "ZRITSDHTM4YORCX3",
		"sasl.password": "iOJVSZ5sHVRnmunF7VvCw+lC1iADXyNZGeYuVlZZfUlvcvUn4fotwbsxRoW2WY2W",
		"sasl.mechanism": "PLAIN",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	fmt.Println("START")

	if err != nil {
		panic(err)
	}

	

	consumer.SubscribeTopics([]string{"email_topic"}, nil)
	
	defer consumer.Close()

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			return
		}

		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
		// var message map[string]interface{}
		var message struct{
			AmbassadorEmail string
			AmbassadorRevenue float32
			Code string
			Id int
			AdminRevenue float32
		}

		json.Unmarshal(msg.Value, &message)

		// emailAmbassador, ok := message["ambassadorEmail"].(string)
		// if !ok {
		// 	emailAmbassador = "unknown"
		// }

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", message.AmbassadorRevenue, message.Code))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{message.AmbassadorEmail}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed",message.Id, message.AdminRevenue))

		smtp.SendMail("host.docker.internal:1025", nil, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)
	}

	


}