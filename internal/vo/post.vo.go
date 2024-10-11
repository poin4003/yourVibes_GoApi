package vo

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

type UpdatePostInput struct {
	Title    *string              `json:"title" binding:"required"`
	Content  *string              `json:"content" binding:"required"`
	Privacy  *consts.PrivacyLevel `json:"privacy" binding:"required,privacy_enum"`
	Location *string              `json:"location,omitempty"`
}
