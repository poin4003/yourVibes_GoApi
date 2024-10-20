package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"mime/multipart"
)

type (
	IPostUser interface {
		CreatePost(ctx context.Context, postModel *model.Post, inMedia []multipart.File) (post *model.Post, resultCode int, err error)
		UpdatePost(ctx context.Context, postId uuid.UUID, updateData map[string]interface{}, deleteMediaIds []uint, inMedia []multipart.File) (post *model.Post, resultCode int, err error)
		DeletePost(ctx context.Context, postId uuid.UUID) (resultCode int, err error)
		GetPost(ctx context.Context, postId uuid.UUID) (post *model.Post, resultCode int, err error)
		GetManyPosts(ctx context.Context, query *query_object.PostQueryObject) (posts []*model.Post, resultCode int, response *response.PagingResponse, err error)
	}
	IPostLike interface {
		LikePost(ctx context.Context, likeUserPost *model.LikeUserPost) error
		DeleteLikePost(ctx context.Context, likeUserPost *model.LikeUserPost) error
		GetUsersOnLikes(ctx context.Context, postId uuid.UUID) ([]*model.User, error)
	}
)

var (
	localPostUser     IPostUser
	localLikeUserPost IPostLike
)

func PostUser() IPostUser {
	if localPostUser == nil {
		panic("repository_implement localPostUser not found for interface IPostUser")
	}

	return localPostUser
}

func LikeUserPost() IPostLike {
	if localLikeUserPost == nil {
		panic("repository_implement localLikeUserPost not found for interface ILikeUserPost")
	}

	return localLikeUserPost
}

func InitPostUser(i IPostUser) {
	localPostUser = i
}

func InitLikeUserPost(i IPostLike) {
	localLikeUserPost = i
}
