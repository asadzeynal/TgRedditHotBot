package tgServer

import (
	"fmt"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"gopkg.in/telebot.v3"
	"time"
)

func Start(config util.Config) error {
	pref := telebot.Settings{
		Token:  config.TgToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		return fmt.Errorf("failed to start bot: %v", err)
	}

	bot.Handle("/start", func(ctx telebot.Context) error {
		return ctx.Send("Hello!")
	})

	bot.Start()

	return nil
}
