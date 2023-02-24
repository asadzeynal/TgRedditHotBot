package db

type ConfigData struct {
	RedditAccessToken      string `json:"RedditAccessToken"`
	RedditTokenToRefreshAt string `json:"RedditTokenToRefreshAt"`
}
