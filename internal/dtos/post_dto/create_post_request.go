package post_dto

import (
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
)

type CreatePostInput struct {
	Title    string                 `form:"title" binding:"required"`
	Content  string                 `form:"content" binding:"required"`
	Privacy  consts.PrivacyLevel    `form:"privacy" binding:"privacy_enum"`
	Location string                 `form:"location,omitempty"`
	Media    []multipart.FileHeader `form:"media,omitempty" binding:"file"`
}
