package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"mime/multipart"
)

type (
	IPostUser interface {
		CreatePost(ctx context.Context, postModel *model.Post, inMedia []multipart.File) (post *model.Post, resultCode int, err error)
		UpdatePost(ctx context.Context, postId uuid.UUID, updateData map[string]interface{}, deleteMediaIds []uint, inMedia []multipart.File) (post *model.Post, resultCode int, err error)
		DeletePost(ctx context.Context, email string) (resultCode int, err error)
		GetPost(ctx context.Context, postId uuid.UUID) (post *model.Post, resultCode int, err error)
		GetManyPost(ctx context.Context, query *query_object.PostQueryObject) (posts []*model.Post, resultCode int, err error)
	}
)

var (
	localPostUser IPostUser
	//localUserInfo  IUserInfo
)

func PostUser() IPostUser {
	if localPostUser == nil {
		panic("repository_implement localPostUser not found for interface IPostUser")
	}

	return localPostUser
}

func InitPostUser(i IPostUser) {
	localPostUser = i
}
