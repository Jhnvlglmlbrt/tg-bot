package main

import (
	"flag"
	"log"

	"github.com/Jhnvlglmlbrt/tg-bot/clients/telegram"
)

func main() {

	tgClient = telegram.New(mustTokenHost())

	// consumer.Start(fetcher, processor)

	// fetcher = fetcher.New()

	// processor = fetcher.New()

}

func mustTokenHost() (string, string) {

	token := flag.String("t", "", "token for access to tg bot")
	host := flag.String("h", "", "host of telegram api")

	flag.Parse()

	if *host == "" || *token == "" {
		log.Fatal("Host and token must be specified")
	}

	return *token, *host
}
