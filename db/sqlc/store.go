package db

import (
	"context"
	"database/sql"
	"fmt"
)

type FullPost struct {
	Post
	Image       PostImage
	Video       PostVideo
	ContentType string
}

type Store interface {
	Querier
	FetchFullRandomPost(ctx context.Context) (FullPost, error)
	ExecTx(ctx context.Context, fn func(queries *Queries) error) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) FetchFullRandomPost(ctx context.Context) (FullPost, error) {
	p, err := store.GetRandomPost(ctx)
	if err != nil {
		return FullPost{}, fmt.Errorf("Unable to fetch random post: %v", err)
	}

	img, err := store.GetImagesByPost(ctx, p.ID)
	if err != nil {
		return FullPost{}, fmt.Errorf("Unable to fetch random post: %v", err)
	}

	vid, err := store.GetVideosByPost(ctx, p.ID)
	if err != nil {
		return FullPost{}, fmt.Errorf("Unable to fetch random post: %v", err)
	}

	var contentType string
	postImage := PostImage{}
	postVideo := PostVideo{}
	if len(img) != 0 {
		contentType = "image"
		postImage = img[0]
	} else if len(vid) != 0 {
		contentType = "video"
		postVideo = vid[0]
	} else {
		return FullPost{}, fmt.Errorf("Post without video or image: %v", p.ID)
	}

	return FullPost{Post: p, Image: postImage, Video: postVideo, ContentType: contentType}, nil

}

func (store *SQLStore) ExecTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
