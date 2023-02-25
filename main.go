package main

import (
	"context"
	"log"
	"time"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/populator"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/tgServer"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"github.com/jackc/pgx/v5"
)

const (
	dbPopulationInterval = time.Hour
)

var logger util.Logger = util.NewCustomLog()

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	store := db.NewStore(conn, logger)

	client, err := rdClient.New(config, store)
	if err != nil {
		log.Fatalf("failed to start reddit client: %v", err)
	}

	err = scheduleDbPopulation(store, client, dbPopulationInterval)
	if err != nil {
		log.Fatalf("failed to perform initial populator run: %v", err)
	}

	err = tgServer.Start(config, client, store)
	if err != nil {
		log.Fatalf("failed to start tg server: %v", err)
	}
}

func scheduleDbPopulation(store db.Store, client *rdClient.Client, interval time.Duration) error {
	ticker := time.NewTicker(10 * time.Second)
	go func(s db.Store, c *rdClient.Client) {
		for {
			<-ticker.C
			err := populator.Run(s, c)
			if err != nil {
				logger.Error("Unable to populate db: %v ", err)
				// Try again in a minute
				ticker.Reset(60 * time.Second)
			}
			err = populator.RefreshPostsCount(s)
			if err != nil {
				logger.Warn("Unable to refresh posts count: %v ", err)
			}

			ticker.Reset(interval)
		}
	}(store, client)

	return nil
}
