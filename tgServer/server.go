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
	btnMorePosts = menu.Text("Give me a post please!")
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

	menu.Reply(menu.Row(btnMorePosts))

	bot.Handle("/start", server.start)
	bot.Handle(&btnMorePosts, server.getRandomPost)

	log.Println("starting reddit tg server")
	bot.Start()

	return nil
}
