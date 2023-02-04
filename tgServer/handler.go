package tgServer

import (
	"fmt"
	"gopkg.in/telebot.v3"
)

func (s *Server) start(ctx telebot.Context) error {
	err := ctx.Send("Hi! I will send you random hot content from all over reddit!", menu)
	if err != nil {
		return fmt.Errorf("could not send response: %v", err)
	}
	return nil
}

func (s *Server) getRandomPost(ctx telebot.Context) error {
	post, err := s.rdClient.FetchRandomPost()
	if err != nil {
		err = ctx.Send("Please try again later", menu)
		if err != nil {
			return fmt.Errorf("could not get postResponse and could not send response: %v", err)
		}
		return fmt.Errorf("error while retrieving post: %v ", err)
	}

	caption := fmt.Sprintf("%s\n\n%s", post.Title, post.Url)

	if post.ImageUrl != "" {
		photo := &telebot.Photo{File: telebot.FromURL(post.ImageUrl), Caption: caption}
		err = ctx.SendAlbum(telebot.Album{photo}, menu)
		if err != nil {
			return fmt.Errorf("could not send response\n url: %v\n error: %v", post.ImageUrl, err)
		}
		return nil
	}

	video := &telebot.Video{
		File:      telebot.FromURL(post.Video.Url),
		Caption:   caption,
		Streaming: true,
		Width:     post.Video.Width,
		Height:    post.Video.Height,
	}

	err = ctx.SendAlbum(telebot.Album{video}, menu)
	if err != nil {
		return fmt.Errorf("could not send response\n url: %v\n error: %v", post.Video, err)
	}

	return nil
}
