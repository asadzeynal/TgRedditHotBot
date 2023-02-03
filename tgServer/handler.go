package tgServer

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
)

func (s *Server) start(ctx telebot.Context) error {
	err := ctx.Send("Hi! I will send you random hot content from all over reddit!")
	if err != nil {
		return fmt.Errorf("could not send response: %v", err)
	}
	return nil
}

func (s *Server) getRandomPost(ctx telebot.Context) {
	err := ctx.Send("Hi! I will send you random hot content from all over reddit!")
	if err != nil {
		log.Println(fmt.Errorf("could not send response: %v", err))
	}
}
