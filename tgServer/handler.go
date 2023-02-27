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
	var toSend telebot.Sendable

	storeFileId := func(postId string, fileId string) error {
		switch post.ContentType {
		case "video":
			return s.store.SetVideoFileId(context.Background(), db.SetVideoFileIdParams{Post: postId, TgFileID: fileId})
		default:
			return s.store.SetImageFileId(context.Background(), db.SetImageFileIdParams{Post: postId, TgFileID: fileId})
		}
	}

	switch post.ContentType {
	case "image":
		var file telebot.File
		if post.Image.TgFileID != "" {
			file = telebot.File{FileID: post.Image.TgFileID}
		} else {
			file = telebot.FromURL(post.Image.Url)
		}
		toSend = &telebot.Photo{File: file, Caption: caption}
	case "gif":
		var file telebot.File
		if post.Image.TgFileID != "" {
			file = telebot.File{FileID: post.Image.TgFileID}
		} else {
			file = telebot.FromURL(post.Image.Url)
		}
		toSend = &telebot.Animation{File: file, Caption: caption}
	case "video":
		var file telebot.File
		if post.Video.TgFileID != "" {
			file = telebot.File{FileID: post.Video.TgFileID}
		} else {
			file = telebot.FromURL(post.Video.Url)
		}
		toSend = constructVideo(post, caption, file)

	default:
		return fmt.Errorf("unknown content type: %v postId: %v", post.ContentType, post.ID)
	}

	err = ctx.Send(toSend, menu)
	if err != nil {
		ctx.Send("Please try again later", menu)
		return fmt.Errorf("could not send response. post: %v, error: %v", post.ID, err)
	}

	var fileId, storedFileId string
	if toSendImg, ok := toSend.(*telebot.Photo); ok {
		fileId = toSendImg.FileID
		storedFileId = post.Image.TgFileID
	}
	if toSendVid, ok := toSend.(*telebot.Video); ok {
		fileId = toSendVid.FileID
		storedFileId = post.Video.TgFileID
	}
	if toSendGif, ok := toSend.(*telebot.Animation); ok {
		fileId = toSendGif.FileID
		storedFileId = post.Image.TgFileID
	}

	if storedFileId == "" && fileId != "" {
		err = storeFileId(post.ID, fileId)
		if err != nil {
			fmt.Printf("could not update fileId: %v ", err)
		}
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
