package util

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	RedditAuth        string `mapstructure:"REDDIT_AUTH"`
	TgToken           string `mapstructure:"TG_TOKEN"`
	RedditAccessToken string `mapstructure:"REDDIT_ACCESS_TOKEN"`
	TokenRefreshAt    string `mapstructure:"TOKEN_REFRESH_AT"`
	DBDriver          string `mapstructure:"DB_DRIVER"`
	DBSource          string `mapstructure:"DB_SOURCE"`
}

// In order for this to work with environment variables, the project has to have a .env file with all vars listed (can be empty)
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

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
