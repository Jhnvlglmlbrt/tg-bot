package main

import (
	"flag"
	"log"

	tgClient "github.com/Jhnvlglmlbrt/tg-bot/clients/telegram"
	event_consumer "github.com/Jhnvlglmlbrt/tg-bot/consumer/event-consumer"
	"github.com/Jhnvlglmlbrt/tg-bot/events/telegram"
	"github.com/Jhnvlglmlbrt/tg-bot/storage/files"
)

const (
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(mustTokenHost()),
		files.New(storagePath),
	)

	log.Print("Service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustTokenHost() (string, string) {

	host := flag.String("h", "", "host of telegram api")

	token := flag.String("t", "", "token for access to tg bot")

	flag.Parse()

	if *host == "" || *token == "" {
		log.Fatal("Host and token must be specified")
	}

	return *host, *token
}
