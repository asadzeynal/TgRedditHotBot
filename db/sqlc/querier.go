// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.0

package db

import (
	"context"
)

type Querier interface {
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreatePostImage(ctx context.Context, arg CreatePostImageParams) (PostImage, error)
	CreatePostVideo(ctx context.Context, arg CreatePostVideoParams) (PostVideo, error)
	GetImagesByPost(ctx context.Context, post string) ([]PostImage, error)
	GetRandomPost(ctx context.Context) (Post, error)
	GetVideosByPost(ctx context.Context, post string) ([]PostVideo, error)
}

var _ Querier = (*Queries)(nil)
