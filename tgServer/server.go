package tgServer

import (
	"fmt"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"gopkg.in/telebot.v3"
	"log"
	"time"
)

var (
	// Universal markup builders.
	menu = &telebot.ReplyMarkup{ResizeKeyboard: true}
	// Reply buttons.
	btnHelp = menu.Text("ℹ Help")
)

type Server struct {
	rdClient *rdClient.Client
}

func Start(config util.Config, client *rdClient.Client) error {
	server := Server{client}

	pref := telebot.Settings{
		Token:  config.TgToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		return fmt.Errorf("failed to start bot: %v", err)
	}

	bot.Handle("/start", server.start)
	bot.Handle(&btnHelp, server.getRandomPost)

	bot.Start()

	log.Println("Successfully started reddit bot server")
	return nil
}
