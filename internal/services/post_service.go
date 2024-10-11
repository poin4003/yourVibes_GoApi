package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/vo"
	"mime/multipart"
)

type (
	IPostUser interface {
		CreatePost(ctx context.Context, inPostData *vo.CreatePostInput, inMedia []multipart.File, userId uuid.UUID) (post *model.Post, resultCode int, err error)
		UpdatePost(ctx context.Context, in *vo.UpdatePostInput) (post *model.Post, resultCode int, err error)
		DeletePost(ctx context.Context, email string) (resultCode int, err error)
		GetPost(ctx context.Context, query interface{}, args ...interface{}) (post *model.Post, resultCode int, err error)
		GetAllPost(ctx context.Context) (posts []*model.Post, resultCode int, err error)
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
