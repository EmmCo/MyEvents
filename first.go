package dblayer

import (
	"awesomeProject/configuration"
	"awesomeProject/eventservice/rest"
	"awesomeProject/peresistence/dblayer"
	"flag"
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"os"

	//_ "github.com/go-sql-driver/mysql"
	//"github.com/gorilla/mux"
	//"net/http"

	//"github.com/jmoiron/sqlx"
)

func main() {
	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672"
	}
	connection, err := amqp.Dial(amqpURL)
	if err != nil {
		panic("could not establish AMQP connection: " + err.Error())
	}
	channel, err := connection.Channel()
	if err != nil {
		panic("could not open channel: " + err.Error())
	}
	channel.ExchangeDeclare("events", "topic",
		true, false, false, false, nil)
	message := amqp.Publishing{
		Body : []byte("Hello world"),
	}
	err = channel.Publish("event", "some-routing-key",
		false, false, message)
	if err != nil {
		panic("error while publishing message: " + err.Error())
	}
	_, err = channel.QueueDeclare("my_queue",
		true, false, false, false, nil)
	if err != nil {
		panic("error while declaring the queue: " + err.Error())
	}

	err = channel.QueueBind("my_queue", "#", "events",
		false, nil)
	if err != nil {
		panic("error while binding the queue: " + err.Error())
	}

	msgs, err := channel.Consume("my_queue", "",
		false, false, false, false, nil)
	if err != nil {
		panic("error while consuming the queue: " + err.Error())
	}
	for msg := range msgs {
		fmt.Println("message received: " + string(msg.Body))
		msg.Ack(false)
	}

	confPath := flag.String("conf", `.\configuration\config.json`, "flag to set the path to the configuration json file")
	flag.Parse()
	//extract configuration
	config, _ := configuration.ExtractConfiguration(*confPath)
	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	//RESTful API start
	httpErrChan, httptlsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndPint, dbhandler)
	select {
	case err := <-httpErrChan:
		log.Fatal("HTTP Error: ", err)
	case err := <-httptlsErrChan:
		log.Fatal("HTTPS Error: ", err)
	}
}
