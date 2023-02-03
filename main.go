package main

import (
	"github.com/asadzeynal/TgRedditHotBot/tgServer"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config: %v", err)
	}
	//err = rdClient.Start(config)
	//if err != nil {
	//	log.Fatal("Failed to start reddit client")
	//}

	err = tgServer.Start(config)
	if err != nil {
		log.Fatal("Failed to start reddit client", err)
	}
}
