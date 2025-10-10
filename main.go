package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/d34ckgler/send-broker-event/rabbit"
	"github.com/d34ckgler/send-broker-event/types"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
)

// DECLARE CONST
const ENV_FILE = ".env"

// ********************
var config types.Config

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s : %s", msg, err)
	}
}

func SendNotification(data []byte) {
	// Connect: se conecta al servidor de RabbitMQ
	// Parametros:
	// - r: es un objeto de tipo Rabbit que contiene la informacion de conexion al servidor de RabbitMQ
	// Retorno:
	// - rabbit: es un objeto de tipo Rabbit que contiene la informacion de conexion al servidor de RabbitMQ
	// - error: si ocurre un error al conectarse al servidor de RabbitMQ, se retorna el error, de lo contrario se retorna nil
	rabbit, _ := rabbit.Connect(&rabbit.Rabbit{
		Host:     config.RabbitHost,  // el host del servidor de RabbitMQ
		Port:     config.RabbitPort,  // el puerto del servidor de RabbitMQ
		User:     config.RabbitUser,  // el usuario de autenticacion del servidor de RabbitMQ
		Password: config.RabbitPass,  // la contraseña de autenticacion del servidor de RabbitMQ
		Vhost:    config.RabbitVHost, // el vhost del servidor de RabbitMQ
	})

	// Cierra la conexión a RabbitMQ
	defer rabbit.Close()

	// CreateChannel: crea un canal de comunicación con RabbitMQ
	// Parametros:
	// - r: es un objeto de tipo Rabbit que contiene la informacion de conexion al servidor de RabbitMQ
	// Retorno:
	// - channel: es un objeto de tipo amqp.Channel que representa el canal de comunicación con RabbitMQ
	// - error: si ocurre un error al crear el canal de comunicación, se retorna el error, de lo contrario se retorna nil
	channel, _ := rabbit.CreateChannel()
	defer channel.Close()

	/*
		ExchangeDeclare: es un metodo de la estructura channel que se utiliza para declarar un exchange en RabbitMQ.
		Recibe los siguientes parametros:
		- name: es el nombre del exchange que se desea declarar.
		- kind: es el tipo de exchange que se desea declarar, en este caso se declarara un exchange de tipo fanout.
		- durable: determina si el exchange que se declaro es duradero o no. En este caso se declara como false.
		- autoDelete: determina si el exchange se borra automaticamente cuando ya no se usa. En este caso se declara como false.
		- internal: determina si el exchange es interno o no. En este caso se declara como false.
		- noWait: determina si se esperara a que el exchange se declare antes de continuar. En este caso se declara como false.
		- args: es un mapa que contiene argumentos adicionales que se pueden pasar al exchange. En este caso se pasa nil.
	*/
	err := channel.ExchangeDeclare(config.RabbitExchange, "direct", false, false, false, false, nil)
	failOnError(err, "Failed to declare a exchange")

	q, err := channel.QueueDeclare(
		config.RabbitQueue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.QueueBind(
		q.Name,                   // queue name
		config.RabbitExchangeKey, // routing key
		config.RabbitExchange,    // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	/*
		ContextWithTimeout: Crea un nuevo contexto con un tiempo de espera especificado.
		Recibe los siguientes parametros:
		- parent: es el contexto padre, si blackxto con el tiempo de espera especificado.
		- context.CancelFunc: es la funcion que se utiliza para cancelar el contexto.
	*/
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// product := Product{
	// 	Id:   132,
	// 	Name: "Chocolate",
	// }

	// body, err := json.Marshal(products)
	// failOnError(err, "Error al parsear el objeto")

	/*
		PublishWithContext: publica un mensaje en un exchange con un contexto especificado.
		Recibe los siguientes parametros:
		- ctx: es el contexto que se utilizará para determinar el tiempo de espera antes de cancelar la publicacion.
		- exchangeName: es el nombre del exchange al cual se enviara el mensaje.
		- routingKey: es la clave de enrutamiento que se utilizara para determinar a que cola se enviara el mensaje.
		- mandatory: determina si se debe enviar el mensaje a al menos una cola.
		- immediate: determina si se debe enviar el mensaje a una cola que este disponible inmediatamente.
		- msg: es el objeto de tipo amqp.Publishing que contiene la informacion del mensaje a enviar.
		Retorno:
		- error: si ocurre un error al publicar el mensaje, se retorna el error, de lo contrario se retorna nil.
	*/
	err = channel.PublishWithContext(ctx,
		config.RabbitExchange,    // exchange
		config.RabbitExchangeKey, // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			// Body:        []byte(message),
			Body: data,
		})
	failOnError(err, "Failed to publish a message")
	log.Println(" [x] Sent notify")
}

func init() {
	err := godotenv.Load(ENV_FILE)
	failOnError(err, "Cannot load .env file")

	err = envconfig.Process("CENTRAL", &config)
	failOnError(err, "Cannot load .env file")

	os.Setenv("TZ", "America/Caracas")
	time.Local = time.FixedZone("VET", -4*60*60)
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatal("Missing arguments")
	}

	binaryData := map[string]interface{}{
		"event":  args[0],
		"action": 1,
	}
	data, _ := json.Marshal(binaryData)

	SendNotification(data)
}
