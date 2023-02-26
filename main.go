package main

import (
	"context"
	"fmt"
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

	// var logger util.Logger = util.NewCustomLog()

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

	err = LoadRedditConfig(store, config)
	if err != nil {
		sugar.Infof("Unable to load config: %v\n", err)
	}

	internal.ScheduleTokenUpdate(config, store)

	client := rdClient.New(config, store)

	internal.ScheduleDbPopulation(store, client, dbPopulationInterval)
	if err != nil {
		log.Fatalf("failed to perform initial populator run: %v", err)
	}

	err = tgServer.Start(config, store)
	if err != nil {
		log.Fatalf("failed to start tg server: %v", err)
	}
}

func LoadRedditConfig(store db.Store, config *util.Config) error {
	redditConf, err := store.GetConfig(context.Background(), "reddit")
	if err != nil {
		return fmt.Errorf("could not load config from DB: %v", err)
	}
	token := util.Decrypt(redditConf.Data.RedditAccessToken, config.EncryptionKey)

	config.Set("TGRHB_REDDIT_ACCESS_TOKEN", token)
	config.Set("TGRHB_TOKEN_REFRESH_AT", redditConf.Data.RedditTokenToRefreshAt)
	return nil
}
