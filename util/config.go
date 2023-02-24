package util

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	RedditAuth        string `koanf:"TGRHB_REDDIT_AUTH"`
	TgToken           string `koanf:"TGRHB_TG_TOKEN"`
	RedditAccessToken string `koanf:"TGRHB_REDDIT_ACCESS_TOKEN"`
	TokenRefreshAt    string `koanf:"TGRHB_TOKEN_REFRESH_AT"`
	DBDriver          string `koanf:"TGRHB_DB_DRIVER"`
	DBSource          string `koanf:"TGRHB_DB_SOURCE"`
}

const envVariableName = "TGRHB_ENV"

var k *koanf.Koanf

// In order for this to work with environment variables, the project has to have a .env file with all vars listed (can be empty)
func LoadConfig(path string) (Config, error) {
	k = koanf.New(path)
	// local or prod for now
	environment := os.Getenv(envVariableName)

	switch environment {
	case "prod":
		e := env.Provider("TGRHB_", ".", nil)
		if err := k.Load(e, nil); err != nil {
			return Config{}, fmt.Errorf("error loading config: %v", err)
		}
	default:
		environment = "local"
		f := file.Provider("app.dev.env")
		if err := k.Load(f, dotenv.Parser()); err != nil {
			return Config{}, fmt.Errorf("error loading config: %v", err)
		}
	}

	fmt.Printf("Loading config on environment: %v\n", environment)

	config := Config{}
	k.Unmarshal("", &config)

	return config, nil
}

func (c *Config) Set(key string, value string) {
	k.Set(key, value)
}
