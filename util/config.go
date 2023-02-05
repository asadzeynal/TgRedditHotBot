package util

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	RedditAuth        string `mapstructure:"reddit-auth"`
	TgToken           string `mapstructure:"tg-token"`
	RedditAccessToken string `mapstructure:"reddit-access-token"`
	TokenRefreshAt    string `mapstructure:"token-refresh-at"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func (c *Config) Set(key string, value string) error {
	viper.Set(key, value)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("unable to write config: %v", err)
	}
	err = viper.Unmarshal(&c)
	return err
}
