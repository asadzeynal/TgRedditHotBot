package tgServer

import (
	"fmt"
	"log"
	"net/http"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/util"
	"gopkg.in/telebot.v3"
)

var (
	// Universal markup builders.
	menu = &telebot.ReplyMarkup{ResizeKeyboard: true}
	// Reply buttons.
	btnMorePosts = menu.Text("Give me a post please!")
)

type Server struct {
	rdClient *rdClient.Client
	store    db.Store
}

func Start(config util.Config, client *rdClient.Client, store db.Store) error {
	server := Server{client, store}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { fmt.Fprintf(w, "ok") })
	http.ListenAndServe(":8090", nil)

	pref := telebot.Settings{
		Token: config.TgToken,
		// Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := telebot.NewBot(pref)
	if err != nil {
		return fmt.Errorf("failed to start bot: %v", err)
	}

	menu.Reply(menu.Row(btnMorePosts))

	bot.Handle("/start", server.start)
	bot.Handle(&btnMorePosts, server.getRandomPost)
	bot.SetWebhook(&telebot.Webhook{Listen: ":8080"})

	log.Println("starting reddit tg server")
	bot.Start()

	return nil
}
