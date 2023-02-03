package main

import (
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/tgServer"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}

	client, err := rdClient.New(config)
	if err != nil {
		log.Fatalf("failed to start reddit client: %v", err)
	}

	err = tgServer.Start(config, client)
	if err != nil {
		log.Fatalf("failed to start tg server: %v", err)
	}
}
