package main

import (
	"encoding/json"
	"fmt"
	"net/smtp"

	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	BOOTSRAP_SERVER = "BOOTSTRAP_SERVERS"
	SERCURITY_PROTOCOL = "SECURITY_PROTOCOL"
	SASL_USERNAME = "SASL_USERNAME"
	SASL_PASSWORD = "SASL_PASSWORD"
	SASL_MECHANISM = "SASL_MECHANISM"
	KAFKA_TOPIC = "KAFKA_TOPIC"
	EMAIL_HOST = "EMAIL_HOST"
	EMAIL_PORT = "EMAIL_PORT"
	EMAIL_USERNAME = "EMAIL_USERNAME"
	EMAIL_PASSWORD = "EMAIL_PASSWORD"
)

var (
	bootstrap_server = os.Getenv(BOOTSRAP_SERVER)
	security_protocol = os.Getenv(SERCURITY_PROTOCOL)
	sasl_username = os.Getenv(SASL_USERNAME)
	sasl_password = os.Getenv(SASL_PASSWORD)
	sasl_mechanism = os.Getenv(SASL_MECHANISM)
	kafka_topic = os.Getenv(KAFKA_TOPIC)
	email_host = os.Getenv(EMAIL_HOST)
	email_port = os.Getenv(EMAIL_PORT)
	email_username = os.Getenv(EMAIL_USERNAME)
	email_password = os.Getenv(EMAIL_PASSWORD)
)

func main() {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrap_server,
		"security.protocol": security_protocol,
		"sasl.username": sasl_username,
		"sasl.password": sasl_password,
		"sasl.mechanism": sasl_mechanism,
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	fmt.Println("START")

	if err != nil {
		panic(err)
	}

	

	consumer.SubscribeTopics([]string{kafka_topic}, nil)
	
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

		auth := smtp.PlainAuth("", email_username,email_password, email_host)

		ambassadorMessage := []byte(fmt.Sprintf("You earned $%f from the link #%s", message.AmbassadorRevenue, message.Code))

		smtp.SendMail(fmt.Sprintf("%s:%s",email_host,email_port), auth, "no-reply@email.com", []string{message.AmbassadorEmail}, ambassadorMessage)

		adminMessage := []byte(fmt.Sprintf("Order #%d with a total of $%f has been completed",message.Id, message.AdminRevenue))

		smtp.SendMail(fmt.Sprintf("%s:%s",email_host,email_port), auth, "no-reply@email.com", []string{"admin@admin.com"}, adminMessage)
	}

	


}