package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/asadzeynal/TgRedditHotBot/util"
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
	db     *sql.DB
	logger *util.Logger
}

func NewStore(db *sql.DB, logger *util.Logger) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
		logger:  logger,
	}
}

func (store *SQLStore) RefreshPostsCount(ctx context.Context) error {
	_, err := store.db.ExecContext(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY posts_count;")
	if err != nil {
		return fmt.Errorf("Unable to refresh posts_count: %v", err)
	}
	store.logger.Info("Refreshed posts_count materialized view")
	return nil
}

func (store *SQLStore) FetchFullRandomPost(ctx context.Context) (FullPost, error) {
	postsCount, err := store.GetTotalCount(ctx)
	if err != nil {
		return FullPost{}, fmt.Errorf("Unable to fetch posts count: %v", err)
	}

	postRowNum := util.RandomInRange(0, int(postsCount))
	store.logger.Info("total count: %v, random num: %v", postsCount, postRowNum)

	p, err := store.GetRandomPost(ctx, int32(postRowNum))
	if err != nil {
		return FullPost{}, fmt.Errorf("Unable to fetch random post: %v", err)
	}
	store.logger.Info("post at random num %v: %v", postRowNum, p)

	imgs, err := store.GetImagesByPost(ctx, p.ID)
	if err != nil {
		return FullPost{}, fmt.Errorf("Unable to fetch random post: %v", err)
	}
	store.logger.Info("postId: %v, postImages: %v", p.ID, imgs)

	vids, err := store.GetVideosByPost(ctx, p.ID)
	if err != nil {
		return FullPost{}, fmt.Errorf("Unable to fetch random post: %v", err)
	}
	store.logger.Info("postId: %v, postVideos: %v", p.ID, vids)

	var contentType string
	postImage := PostImage{}
	postVideo := PostVideo{}
	if len(imgs) != 0 {
		if imgs[0].IsGif {
			contentType = "gif"
		} else {
			contentType = "image"
		}
		postImage = imgs[0]
	} else if len(vids) != 0 {
		contentType = "video"
		postVideo = vids[0]
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
