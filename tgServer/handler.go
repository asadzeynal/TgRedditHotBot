package tgServer

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"strings"
)

func (s *Server) start(ctx telebot.Context) error {
	err := ctx.Send("Hi! I will send you random hot content from all over reddit!", menu)
	if err != nil {
		return fmt.Errorf("could not send response: %v", err)
	}
	return nil
}

func (s *Server) getRandomPost(ctx telebot.Context) error {
	postResponse, err := s.rdClient.FetchRandomPost()
	if err != nil {
		err = ctx.Send("Please try again later", menu)
		if err != nil {
			return fmt.Errorf("could not get postResponse and could not send response: %v", err)
		}
		return fmt.Errorf("error while retrieving post: %v ", err)
	}

	parsedUrl := strings.ReplaceAll(postResponse.ImageUrl, "&amp;", "&")
	caption := fmt.Sprintf("%s\n%s", postResponse.Title, postResponse.Url)

	photo := &telebot.Photo{File: telebot.FromURL(parsedUrl), Caption: caption}
	err = ctx.SendAlbum(telebot.Album{photo}, menu)
	if err != nil {
		return fmt.Errorf("could not send response\n oldUrl: %v\n url: %v\n error: %v", postResponse.ImageUrl, parsedUrl, err)
	}
	return nil
}
