package db

import (
	"context"
	"fmt"

	"github.com/asadzeynal/TgRedditHotBot/util"
	"github.com/jackc/pgx/v5"
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
	db     *pgx.Conn
	logger util.Logger
}

func NewStore(db *pgx.Conn, logger util.Logger) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
		logger:  logger,
	}
}

func (store *SQLStore) RefreshPostsCount(ctx context.Context) error {
	_, err := store.db.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY posts_count;")
	if err != nil {
		return fmt.Errorf("unable to refresh posts_count: %v", err)
	}
	store.logger.Info("Refreshed posts_count materialized view")
	return nil
}

func (store *SQLStore) FetchFullRandomPost(ctx context.Context) (FullPost, error) {
	postsCount, err := store.GetTotalCount(ctx)
	if err != nil {
		return FullPost{}, fmt.Errorf("unable to fetch posts count: %v", err)
	}

	postRowNum := util.RandomInRange(0, int(postsCount))
	store.logger.Infow("fetched total num and generated a random num", "totalPostCount", postsCount, "random num", postRowNum)

	p, err := store.GetRandomPost(ctx, int32(postRowNum))
	if err != nil {
		return FullPost{}, fmt.Errorf("unable to fetch random post: %v", err)
	}
	store.logger.Info("post fetched", "postOffset", postRowNum, "post", p)

	imgs, err := store.GetImagesByPost(ctx, p.ID)
	if err != nil {
		return FullPost{}, fmt.Errorf("unable to fetch random post: %v", err)
	}

	vids, err := store.GetVideosByPost(ctx, p.ID)
	if err != nil {
		return FullPost{}, fmt.Errorf("unable to fetch random post: %v", err)
	}
	store.logger.Info("Fetched images and photos", "event", "mediaFetched", "postId", p.ID, "postImages", imgs, "postVideos", vids)

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
	tx, err := store.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
