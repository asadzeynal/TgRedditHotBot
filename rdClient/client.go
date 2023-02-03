package rdClient

import (
	"encoding/json"
	"fmt"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const AuthUrl = "https://www.reddit.com/api/v1/access_token"
const RandomPostUrl = "https://oauth.reddit.com/r/random/hot"
const AuthParam = "grant_type=client_credentials"

type RedditAccessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	RefreshAt   time.Time
}

type Client struct {
	token           *RedditAccessToken
	config          util.Config
	routinePoolSize int
}

func New(config util.Config) (*Client, error) {
	client := Client{
		routinePoolSize: runtime.NumCPU(),
		config:          config,
	}

	err := client.fetchAccessToken()
	if err != nil {
		return &Client{}, fmt.Errorf("error while fetching reddit access token: %v", err)
	}

	client.scheduleTokenUpdate()

	log.Println("Initialized Reddit client")
	return &client, nil
}

// Schedules a reddit token refresh each n seconds, where n = token.ExpiresIn
func (c *Client) scheduleTokenUpdate() {
	ticker := time.NewTicker(c.token.RefreshAt.Sub(time.Now()))
	go func() {
		for {
			<-ticker.C
			err := c.fetchAccessToken()
			if err != nil {
				log.Println(err)
				// Try again in a minute
				ticker.Reset(60 * time.Second)
			}
			ticker.Reset(c.token.RefreshAt.Sub(time.Now()))
		}
	}()
}

func (c *Client) FetchRandomPost() (util.RedditPost, error) {
	resBody := util.RedditPost{}
	for {
		req, err := http.NewRequest(http.MethodGet, RandomPostUrl, nil)
		if err != nil {
			return util.RedditPost{}, fmt.Errorf("error while creating request: %v", err)
		}

		q := req.URL.Query()
		q.Add("limit", "1")
		req.URL.RawQuery = q.Encode()

		authString := fmt.Sprintf("Bearer %s", c.token.AccessToken)
		req.Header.Add("Authorization", authString)
		req.Header.Add("User-Agent", "TgRedditHot/0.0.1")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			return util.RedditPost{}, fmt.Errorf("error while making request: %v", err)
		}
		if res.StatusCode != http.StatusOK {
			log.Println("error:", res.StatusCode)
		}
		defer res.Body.Close()

		err = util.DecodeRedditPost(&res.Body, &resBody)
		if err != nil {
			fmt.Errorf("error from decoder: %v", err)
			continue
		}

		return resBody, nil
	}
}

// TODO: Implement timeout
// Fetches the Reddit access token and saves it to Server struct as RedditAccessToken
func (c *Client) fetchAccessToken() error {
	paramsReader := strings.NewReader(AuthParam)
	req, err := http.NewRequest(http.MethodPost, AuthUrl, paramsReader)
	if err != nil {
		return fmt.Errorf("error while creating request: %v", err)
	}

	authString := fmt.Sprintf("Basic %s", c.config.RedditAuth)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authString)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error while making request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch refresh token, status code: %v", res.StatusCode)
	}
	defer res.Body.Close()

	var token RedditAccessToken
	err = json.NewDecoder(res.Body).Decode(&token)

	token.RefreshAt = time.Now().Add(time.Duration(token.ExpiresIn-60) * time.Second)
	if err != nil {
		return fmt.Errorf("error when processing response: %v", err)
	}

	log.Println("Successfully updated reddit token")
	c.token = &token
	return nil
}
