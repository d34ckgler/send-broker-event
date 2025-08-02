package rabbit

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	Host     string
	Port     int
	User     string
	Password string
	Vhost    string
	conn     *amqp.Connection
}

type IRabbit interface {
	Connect()
	Client() *amqp.Connection
	CreateChannel() (*amqp.Channel, error)
	Close()
	failOnError(err error, msg string)
}

func failOnError(err error, msg string) error {
	if err != nil {
		log.Panicf("%s:%s", msg, err)
		return err
	}

	return nil
}

func Connect(r *Rabbit) (*Rabbit, error) {
	conString := fmt.Sprintf("amqp://%s:%s@%s:%d", r.User, r.Password, r.Host, r.Port)
	conn, err := amqp.Dial(conString)
	// Set Vhost
	conn.Config.Vhost = r.Vhost

	if err == nil {
		r.conn = conn
		return r, nil
	}

	return nil, failOnError(err, "Failed to connect to RabbitMQ")
}

func (r *Rabbit) Client() *amqp.Connection {
	return r.conn
}

func (r *Rabbit) CreateChannel() (*amqp.Channel, error) {
	channel, err := r.conn.Channel()
	if err == nil {
		return channel, nil
	}

	return nil, failOnError(err, "Failed to open a channel")
}

func (r *Rabbit) Close() {
	r.conn.Close()
}
