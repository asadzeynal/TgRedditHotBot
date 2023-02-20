package main

import (
	"database/sql"
	"log"
	"time"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/populator"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/tgServer"
	"github.com/asadzeynal/TgRedditHotBot/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)

	client, err := rdClient.New(config, store)
	if err != nil {
		log.Fatalf("failed to start reddit client: %v", err)
	}

	err = scheduleDbPopulation(store, client, time.Hour)
	if err != nil {
		log.Fatalf("failed to perform initial populator run: %v", err)
	}

	err = tgServer.Start(config, client, store)
	if err != nil {
		log.Fatalf("failed to start tg server: %v", err)
	}
}

func scheduleDbPopulation(store db.Store, client *rdClient.Client, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	go func(s db.Store, c *rdClient.Client) {
		for {
			<-ticker.C
			err := populator.Run(s, c)
			if err != nil {
				log.Println(err)
				// Try again in a minute
				ticker.Reset(60 * time.Second)
			}
			ticker.Reset(interval)
		}
	}(store, client)

	return nil
}
