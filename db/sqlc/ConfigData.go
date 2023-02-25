package db

type ConfigData struct {
	RedditAccessToken      []byte `json:"RedditAccessToken"`
	RedditTokenToRefreshAt string `json:"RedditTokenToRefreshAt"`
}
