package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	sugar.Infow("database connection created", "event", "dbConnected")

	store := db.NewStore(conn, sugar)

	configReadyFromStore, configReadyFromFetch := internal.ScheduleTokenUpdate(config, store)

	select {
	case <-configReadyFromStore:
	case <-configReadyFromFetch:
	}

	sugar.Infow("reddit config ready", "event", "dbConfigReady")

	client := rdClient.New(config, store)

	sugar.Infow("reddit client instance created", "event", "redditClientReady")

	internal.ScheduleDbPopulation(store, client, dbPopulationInterval)

	sugar.Infow("scheduled db population", "event", "dbPopulationScheduled")

	go func() error {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "ok") })
		err := http.ListenAndServe(":8090", nil)
		if err != nil {
			return fmt.Errorf("error while creating http server: %v", err)
		}
		return nil
	}()

	err = tgServer.Start(config, store)
	if err != nil {
		log.Fatalf("failed to start tg server: %v", err)
	}
}
