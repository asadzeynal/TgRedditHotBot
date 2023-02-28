package internal

import (
	"context"
	"fmt"
	"time"

	db "github.com/asadzeynal/TgRedditHotBot/db/sqlc"
	"github.com/asadzeynal/TgRedditHotBot/rdClient"
	"github.com/asadzeynal/TgRedditHotBot/util"
)

func ScheduleDbPopulation(store db.Store, client rdClient.Client, interval time.Duration) {
	initialDuration := 10 * time.Second
	type input struct {
		store  db.Store
		client rdClient.Client
	}
	i := input{
		store:  store,
		client: client,
	}

	var f util.Func[input, struct{}] = func(i input) (time.Duration, struct{}, error) {
		err := populate(i.store, i.client)
		if err != nil {
			return 60 * time.Second, struct{}{}, err
		}
		err = refreshPostsCount(i.store)
		if err != nil {
			return interval, struct{}{}, err
		}

		return interval, struct{}{}, nil
	}

	p := util.Schedule(initialDuration, i, f)
	go func() {
		for {
			_, err := p.Get()
			if err != nil {
				fmt.Printf("could not populate db: %v\n", err)
				continue
			}
		}
	}()

}

func populate(store db.Store, client rdClient.Client) error {
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

			if post.ContentType == "image" || post.ContentType == "gif" {
				_, err := queries.CreatePostImage(context.Background(), db.CreatePostImageParams{
					Post:  post.Id,
					Url:   post.ImageUrl,
					IsGif: post.ContentType == "gif",
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
					AudioUrl: post.Video.AudioUrl,
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

func refreshPostsCount(store db.Store) error {
	if sqlStore, ok := store.(*db.SQLStore); ok {
		err := sqlStore.RefreshPostsCount(context.Background())
		if err != nil {
			return fmt.Errorf("unable to refresh Materialized View posts_count: %v", err)
		}
		return nil
	}

	return fmt.Errorf("store is not of type *db.SQLStore")
}
