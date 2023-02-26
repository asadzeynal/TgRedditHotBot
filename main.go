package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

type RedditAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

const AuthUrl = "https://www.reddit.com/api/v1/access_token"
const RandomPostUrl = "https://oauth.reddit.com/r/all/top"
const AuthParam = "grant_type=client_credentials"

var logger util.Logger = util.NewCustomLog()

func main() {
	var err error
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

	LoadRedditConfig()
	scheduleTokenUpdate()

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

func scheduleTokenUpdate(config *util.Config) {
	expiresAt, err := time.Parse(time.RFC3339, config.TokenRefreshAt)
	if err != nil {
		config.TokenRefreshAt = ""
	}

	initialRefreshAfter := time.Until(expiresAt)
	if initialRefreshAfter < 1 {
		initialRefreshAfter = 1
	}

	var f util.Func[string, RedditAccessToken] = func(redditAuth string) (time.Duration, RedditAccessToken, error) {
		return updateToken(redditAuth)
	}

	p := util.Schedule(initialRefreshAfter, config.RedditAuth, f)
	SaveRedditConfig(token.AccessToken, nextRefreshAt)

}

func updateToken(redditAuth string) (time.Duration, RedditAccessToken, error) {
	token, err := fetchAccessToken(redditAuth)
	if err != nil {
		return 60 * time.Second, RedditAccessToken{}, fmt.Errorf("error while fetching reddit access token: %v", err)
	}
	nextRefreshAt := time.Now().Add(time.Duration(token.ExpiresIn))

	nextRefreshAfter := time.Until(nextRefreshAt) - 60*time.Second
	if nextRefreshAfter < 1 {
		nextRefreshAfter = 1
	}

	return nextRefreshAfter, token, nil
}

func fetchAccessToken(redditAuth string) (RedditAccessToken, error) {
	paramsReader := strings.NewReader(AuthParam)
	req, err := http.NewRequest(http.MethodPost, AuthUrl, paramsReader)
	if err != nil {
		return RedditAccessToken{}, fmt.Errorf("error while creating request: %v", err)
	}

	authString := fmt.Sprintf("Basic %s", redditAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authString)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return RedditAccessToken{}, fmt.Errorf("error while making request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return RedditAccessToken{}, fmt.Errorf("failed to fetch refresh token, status code: %v", res.StatusCode)
	}
	defer res.Body.Close()

	var token RedditAccessToken
	err = json.NewDecoder(res.Body).Decode(&token)

	if err != nil {
		return RedditAccessToken{}, fmt.Errorf("error when processing response: %v", err)
	}
	return token, nil
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

func LoadRedditConfig() error {
	redditConf, err := store.GetConfig(context.Background(), "reddit")
	if err != nil {
		return fmt.Errorf("could not load config from DB: %v", err)
	}
	token := util.Decrypt(redditConf.Data.RedditAccessToken, config.EncryptionKey)

	config.Set("TGRHB_REDDIT_ACCESS_TOKEN", token)
	config.Set("TGRHB_TOKEN_REFRESH_AT", redditConf.Data.RedditTokenToRefreshAt)
	return nil
}

func SaveRedditConfig(token string, refreshAt time.Time) error {
	encryptedAK := util.Encrypt([]byte(token), config.EncryptionKey)
	store.UpdateConfig(context.Background(), db.UpdateConfigParams{
		ConfigType: "reddit",
		Data: db.ConfigData{
			RedditAccessToken:      encryptedAK,
			RedditTokenToRefreshAt: refreshAt.Format(time.RFC3339),
		},
	})
	return nil
}
