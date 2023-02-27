package tgServer

import (
	"context"
	"fmt"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
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
	post, err := s.store.FetchFullRandomPost(context.Background())
	if err != nil {
		errSend := ctx.Send("Please try again later", menu)
		if errSend != nil {
			return fmt.Errorf("could not get postResponse and could not send response: %v", errSend)
		}
		return fmt.Errorf("error while retrieving post: %v ", err)
	}

	caption := fmt.Sprintf("%s\n\n%s", post.Title, post.Url)

	switch post.ContentType {
	case "image":
		var file telebot.File
		if post.Image.TgFileID != "" {
			file = telebot.File{FileID: post.Image.TgFileID}
		} else {
			file = telebot.FromURL(post.Image.Url)
		}
		photo := &telebot.Photo{File: file, Caption: caption}
		err = ctx.SendAlbum(telebot.Album{photo}, menu)
		if err != nil {
			return fmt.Errorf("could not send response\n url: %v\n error: %v", post.Image.Url, err)
		}

		if post.Image.TgFileID == "" && photo.FileID != "" {
			s.store.SetImageFileId(context.Background(), db.SetImageFileIdParams{Post: post.ID, TgFileID: photo.FileID})
		}
	case "gif":
		var file telebot.File
		if post.Image.TgFileID != "" {
			file = telebot.File{FileID: post.Image.TgFileID}
		} else {
			file = telebot.FromURL(post.Image.Url)
		}
		gif := &telebot.Animation{File: file, Caption: caption}
		err = ctx.Send(gif, menu)
		if err != nil {
			return fmt.Errorf("could not send response\n url: %v\n error: %v", post.Image.Url, err)
		}
		if post.Image.TgFileID == "" && gif.FileID != "" {
			s.store.SetImageFileId(context.Background(), db.SetImageFileIdParams{Post: post.ID, TgFileID: gif.FileID})
		}
	case "video":
		var file telebot.File
		if post.Video.TgFileID != "" {
			file = telebot.File{FileID: post.Image.TgFileID}
		} else {
			file = telebot.FromURL(post.Video.Url)
		}
		video := constructVideo(post, caption, file)
		err = ctx.SendAlbum(telebot.Album{video}, menu)
		if err != nil {
			return fmt.Errorf("could not send response\n url: %v\n error: %v", post.Video, err)
		}
		if post.Video.TgFileID == "" && video.FileID != "" {
			s.store.SetVideoFileId(context.Background(), db.SetVideoFileIdParams{Post: post.ID, TgFileID: video.FileID})
		}
	default:
		return fmt.Errorf("unknown content type: %v postId: %v", post.ContentType, post.ID)
	}

	return nil
}

func constructVideo(post db.FullPost, caption string, file telebot.File) *telebot.Video {
	return &telebot.Video{
		File:      file,
		Caption:   caption,
		Streaming: true,
		Width:     int(post.Video.Width),
		Height:    int(post.Video.Height),
	}
}
