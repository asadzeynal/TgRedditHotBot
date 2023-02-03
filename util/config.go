package util

import "github.com/spf13/viper"

type Config struct {
	RedditAuth string `mapstructure:"REDDIT_AUTH"`
	TgToken    string `mapstructure:"TG_TOKEN"`
}

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
