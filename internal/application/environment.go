package application

import (
	"log"
	"os"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"github.com/nats-io/nats.go"
)

type Environment struct {
	Port    string
	NatsUrl string
	EsdbUrl string
}

func GetEventStoreConnection(connectionString string) *esdb.Client {
	settings, err := esdb.ParseConnectionString(connectionString)

	if err != nil {
		panic(err)
	}

	db, err := esdb.NewClient(settings)

	if err != nil {
		panic(err)
	}

	log.Println("connected to eventstore with URL:", connectionString)

	return db
}

func GetNatsConnection(connectionString string) *nats.Conn {
	conn, err := nats.Connect(connectionString)

	if err != nil {
		panic(err)
	}

	log.Println("connected to nats with URL:", connectionString)

	return conn
}

func NewEnvironment() *Environment {
	env := &Environment{
		Port:    os.Getenv("PORT"),
		NatsUrl: os.Getenv("NATS_URL"),
		EsdbUrl: "esdb://localhost:2113?tls=false",
	}
	if env.Port == "" {
		env.Port = ":3000"
	}
	if env.NatsUrl == "" {
		env.NatsUrl = nats.DefaultURL
	}

	return env
}
