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
	RedditAuth        string `koanf:"REDDIT_AUTH"`
	TgToken           string `koanf:"TG_TOKEN"`
	RedditAccessToken string `koanf:"REDDIT_ACCESS_TOKEN"`
	TokenRefreshAt    string `koanf:"TOKEN_REFRESH_AT"`
	DBDriver          string `koanf:"DB_DRIVER"`
	DBSource          string `koanf:"DB_SOURCE"`
}

const envVariableName = "TGREDDITHOTBOT_ENV"

var k = koanf.New(".")

// In order for this to work with environment variables, the project has to have a .env file with all vars listed (can be empty)
func LoadConfig(path string) (Config, error) {
	// local or prod for now
	environment := os.Getenv(envVariableName)

	switch environment {
	case "prod":
		e := env.Provider("", "", nil)
		if err := k.Load(e, dotenv.Parser()); err != nil {
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
