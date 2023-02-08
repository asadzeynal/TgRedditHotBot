package main

import (
	"database/sql"
	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/tgServer"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"log"
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

	err = tgServer.Start(config, client)
	if err != nil {
		log.Fatalf("failed to start tg server: %v", err)
	}
}
