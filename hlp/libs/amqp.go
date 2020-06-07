package libs

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/Mateus-pilo/go-whats-opt/hlp"
)

var connection *session

type session struct {
	*amqp.Connection
	*amqp.Channel
}

func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}


func ConnectionMqp() (session) {

	if connection != nil {
		log.Println("Caiu")
		return *connection
	}

	url := hlp.Config.GetString("AMQP_URL")
	conn, err := amqp.Dial(url)

	if err != nil {
		panic("could not establish connection with RabbitMQ:" + err.Error())
	}
	defer conn.Close()

	channel, err := conn.Channel()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()


	if err != nil {
			panic("could not open RabbitMQ channel:" + err.Error())
	}

	err = channel.ExchangeDeclare("msgSend", "topic", true, false, false, false, nil)

	if err != nil {
			panic(err)
	}


	if err != nil {
		panic("error publishing a message to the queue:" + err.Error())
	}

	// We create a queue named Test
	_, err = channel.QueueDeclare("msgSend", true, false, false, false, nil)
	_, err = channel.QueueDeclare("msgReceive", true, false, false, false, nil)


	if err != nil {
		panic("error declaring the queue: " + err.Error())
	}

	// We bind the queue to the exchange to send and receive data from the queue
	err = channel.QueueBind("msgSend", "#", "events", false, nil)
	err = channel.QueueBind("msgReceive", "#", "events", false, nil)

	if err != nil {
		panic("error binding to the queue: " + err.Error())
	}

	connection = &session{conn,channel}
	return *connection;
}

func failOnError(err error, msg string) {
  if err != nil {
    log.Fatalf("%s: %s", msg, err)
  }
}
