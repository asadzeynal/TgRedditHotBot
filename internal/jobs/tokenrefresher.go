package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/util"
)

type RedditAccessToken struct {
	AccessToken   string `json:"access_token"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	Scope         string `json:"scope"`
	nextRefreshAt time.Time
}

const AuthUrl = "https://www.reddit.com/api/v1/access_token"
const AuthParam = "grant_type=client_credentials"

func ScheduleTokenUpdate(config *util.Config, store db.Store) chan struct{} {
	configReady := make(chan struct{})
	err := LoadRedditConfig(store, config)
	if err != nil {
		fmt.Printf("Unable to load config: %v\n", err)
	}

	if config.RedditAccessToken != "" {
		configReady <- struct{}{}
	}

	expiresAt, err := time.Parse(time.RFC3339, config.TokenRefreshAt)
	if err != nil {
		config.TokenRefreshAt = ""
	}

	initialRefreshAfter := time.Until(expiresAt)
	if initialRefreshAfter < 1 {
		initialRefreshAfter = 1
	}

	var f util.Func[string, RedditAccessToken] = func(redditAuth string) (time.Duration, RedditAccessToken, error) {
		return getToken(redditAuth)
	}

	p := util.Schedule(initialRefreshAfter, config.RedditAuth, f)
	go func() {
		for {
			val, err := p.Get()
			if err != nil {
				fmt.Printf("could not refresh token: %v\n", err)
				continue
			}
			SaveRedditConfig(store, config, val.AccessToken, val.nextRefreshAt, config.EncryptionKey)
			close(configReady)
		}
	}()
	return configReady
}

func getToken(redditAuth string) (time.Duration, RedditAccessToken, error) {
	token, err := fetchToken(redditAuth)
	if err != nil {
		return 60 * time.Second, RedditAccessToken{}, fmt.Errorf("error while fetching reddit access token: %v", err)
	}
	token.nextRefreshAt = time.Now().Add(time.Duration(token.ExpiresIn))

	nextRefreshAfter := time.Until(token.nextRefreshAt) - 60*time.Second
	if nextRefreshAfter < 1 {
		nextRefreshAfter = 1
	}

	return nextRefreshAfter, token, nil
}

func fetchToken(redditAuth string) (RedditAccessToken, error) {
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

func SaveRedditConfig(store db.Store, config *util.Config, token string, refreshAt time.Time, key string) error {
	config.Set("TGRHB_REDDIT_ACCESS_TOKEN", token)
	config.Set("TGRHB_TOKEN_REFRESH_AT", refreshAt.Format(time.RFC3339))
	encryptedAK := util.Encrypt([]byte(token), key)
	store.UpdateConfig(context.Background(), db.UpdateConfigParams{
		ConfigType: "reddit",
		Data: db.ConfigData{
			RedditAccessToken:      encryptedAK,
			RedditTokenToRefreshAt: refreshAt.Format(time.RFC3339),
		},
	})
	return nil
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
