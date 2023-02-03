package tgServer

import (
	"fmt"
	"gopkg.in/telebot.v3"
)

func (s *Server) start(ctx telebot.Context) error {
	menu.Reply(menu.Row(btnHelp))

	err := ctx.Send("Hi! I will send you random hot content from all over reddit!", menu)
	if err != nil {
		return fmt.Errorf("could not send response: %v", err)
	}
	return nil
}

func (s *Server) getRandomPost(ctx telebot.Context) error {
	err := ctx.Send("")
	if err != nil {
		return fmt.Errorf("could not send response: %v", err)
	}
	return nil
}
