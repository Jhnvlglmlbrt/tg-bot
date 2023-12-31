package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/Jhnvlglmlbrt/tg-bot/clients/telegram"
	event_consumer "github.com/Jhnvlglmlbrt/tg-bot/consumer/event-consumer"
	"github.com/Jhnvlglmlbrt/tg-bot/events/telegram"
	"github.com/Jhnvlglmlbrt/tg-bot/storage/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

const (
	storagePath       = "files_storage"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	// s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err = s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(mustTokenHost()),
		s,
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
