package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/nuveo/log"
	"github.com/streadway/amqp"
)

const (
	DefaultQueueName string = "workflows"
)

type RabbitMQSender struct {
	Conn *amqp.Connection
}

func NewSender(serverAddr string, serverPort string, serverUser string, serverPassword string) *RabbitMQSender {
	log.Println("Sender connecting with Messaging Server")
	senderConfig := fmt.Sprintf("amqp://%s:%s@%s:%s/", serverUser, serverPassword, serverAddr, serverPort)

	conn, err := amqp.Dial(senderConfig)
	HandlerError(err, "Failed to connect to RabbitMQ")

	log.Println("Sender connection with Messaging Server done")
	return &RabbitMQSender{
		Conn:    conn,
	}
}

func (rs RabbitMQSender) Send(i interface{}) {
	w, err := json.Marshal(i)

	if err != nil {
		log.Errorf("Error in type conversion\n")
	}

	ch, err := rs.Conn.Channel()
	defer ch.Close()

	HandlerError(err, "Failed to open a channel")

	queue, err := ch.QueueDeclare(
		DefaultQueueName,
		true,
		false,
		false,
		false,
		nil,
	)

	HandlerError(err, "Failed to declared queue workflows")

	err = ch.Publish(
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "application/json",
			Body: w,
		})

	if err != nil {
		log.Errorf("%s\n", err)
	} else {
		log.Println("Workflow added in queue")
	}
}

func HandlerError(err error, msg string) {
	if err != nil {
		log.Errorf("%s: %s", msg, err)
	}
}