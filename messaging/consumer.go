package messaging

import (
	"encoding/json"
	"fmt"
	"github.com/nuveo/log"
	"github.com/streadway/amqp"
	"gitlab.com/emanueljoivo/workflows/storage"
)

type RabbitMQConsumer struct {
	Conn *amqp.Connection
}

func NewConsumer(serverAddr string, serverPort string, serverUser string, serverPassword string) *RabbitMQConsumer {
	log.Println("Consumer connecting with Messaging Server")
	consumerConfig := fmt.Sprintf("amqp://%s:%s@%s:%s/", serverUser, serverPassword, serverAddr, serverPort)

	conn, err := amqp.Dial(consumerConfig)

	HandlerError(err, "Unable to create connection")

	log.Println("Consumer connection with Messaging Server done")
	return &RabbitMQConsumer{Conn: conn}
}

func (rc RabbitMQConsumer) Consume() *storage.Workflow {
	ch, err := rc.Conn.Channel()

	HandlerError(err, "Unable to create channel")

	queue, err := ch.QueueDeclare(
		DefaultQueueName,
		true,
		false,
		false,
		false,
		nil)

	HandlerError(err,"Unable to declare queue")

	msgs, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
		)
	HandlerError(err, "Failed to register a consumer")

	var workflow *storage.Workflow

	msg := <- msgs
	err = json.Unmarshal(msg.Body, &workflow)

	if err != nil {
		log.Errorf("%s\n", err.Error())
	}

	return workflow
}
