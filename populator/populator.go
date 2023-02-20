package populator

import (
	"context"
	"fmt"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
)

func Run(store db.Store, client *rdClient.Client) error {
	posts, err := client.FetchPosts()
	if err != nil {
		return fmt.Errorf("error while retrieving posts: %v ", err)
	}

	for i := range posts {
		store.ExecTx(context.Background(), func(queries *db.Queries) error {
			post := posts[i]
			_, err = queries.CreatePost(context.Background(), db.CreatePostParams{ID: post.Id, Title: post.Title, Url: post.Url})
			if err != nil {
				return fmt.Errorf("unable to store post: %v", err)
			}

			if post.ContentType == "image" {
				_, err := queries.CreatePostImage(context.Background(), db.CreatePostImageParams{
					Post: post.Id,
					Url:  post.ImageUrl,
				})
				if err != nil {
					return fmt.Errorf("unable to store image: %v", err)
				}

			} else if post.ContentType == "video" {
				_, err := queries.CreatePostVideo(context.Background(), db.CreatePostVideoParams{
					Post:     post.Id,
					Height:   int32(post.Video.Height),
					Width:    int32(post.Video.Width),
					Duration: int32(post.Video.Duration),
					Url:      post.Video.Url,
				})
				if err != nil {
					return fmt.Errorf("unable to store image: %v", err)
				}
			}

			return nil
		})
	}

	return nil

}
