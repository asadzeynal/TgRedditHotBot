package rdClient

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/util"
)

type RedditVideo struct {
	Height   int
	Width    int
	Duration int
	Url      string
}

type RedditPost struct {
	Id          string
	ImageUrl    string
	Title       string
	Url         string
	Video       RedditVideo
	ContentType string
}

type Client struct {
	config *util.Config
	store  db.Store
}

const RandomPostUrl = "https://oauth.reddit.com/r/all/top"

var client Client

func New(config *util.Config, store db.Store) Client {
	client = Client{
		config: config,
		store:  store,
	}

	log.Println("Initialized Reddit client")
	return client
}

// FetchRandomPost Fetches 100 top posts from a /r/all subreddit
func (c Client) FetchPosts() ([]*RedditPost, error) {
	var resBody []*RedditPost
	req, err := http.NewRequest(http.MethodGet, RandomPostUrl, nil)
	if err != nil {
		return []*RedditPost{}, fmt.Errorf("error while creating request: %v", err)
	}

	q := req.URL.Query()
	q.Add("limit", "100")
	req.URL.RawQuery = q.Encode()

	authString := fmt.Sprintf("Bearer %s", c.config.RedditAccessToken)
	req.Header.Add("Authorization", authString)
	req.Header.Add("User-Agent", "TgRedditHot/0.0.1")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return []*RedditPost{}, fmt.Errorf("error while making request: %v", err)
	}
	if res.StatusCode != http.StatusOK {
		log.Println("error:", res.StatusCode)
	}
	defer res.Body.Close()

	resBody, err = DecodeRedditResponse(&res.Body)
	if err != nil {
		return []*RedditPost{}, fmt.Errorf("error from decoder: %v", err)
	}

	return resBody, nil
}
