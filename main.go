package main

import (
	"context"
	"log"
	"time"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	internal "github.com/asadzeynal/TgRedditHotBot/internal/jobs"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/tgServer"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const (
	dbPopulationInterval = time.Hour
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	sugar := logger.Sugar()

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	sugar.Infow("starting app", "event", "configLoaded", "env", config.Environment)

	conn, err := pgx.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	store := db.NewStore(conn, sugar)

	configReady := internal.ScheduleTokenUpdate(config, store)
	<-configReady

	client := rdClient.New(config, store)

	internal.ScheduleDbPopulation(store, client, dbPopulationInterval)

	err = tgServer.Start(config, store)
	if err != nil {
		log.Fatalf("failed to start tg server: %v", err)
	}
}
