package main

import (
	"context"
	"fmt"
	"github.com/SteinerLabs/lms/backend/shared/events"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"time"
)

func main() {
	log.Println("Auth service")
	con, err := nats.Connect("nats://nats:4222")
	if err != nil {
		panic(err)
	}
	defer con.Close()

	js, err := jetstream.New(con)
	if err != nil {
		panic(err)
	}

	publisher := events.NewPublisher(js, "auth-service")
	for {
		err = publisher.Publish(
			context.Background(),
			events.NewEvent[any]("TEST.test", "source", "test", "corrId", "causId", "traceId"),
		)
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
	done := make(chan bool)
	<-done
	fmt.Println("Done")
}
