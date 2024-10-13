package post_dto

import (
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
)

type UpdatePostInput struct {
	Title    *string                `form:"title,omitempty"`
	Content  *string                `form:"content,omitempty"`
	Privacy  *consts.PrivacyLevel   `form:"privacy,omitempty" binding:"privacy_enum"`
	Location *string                `form:"location,omitempty"`
	MediaIDs []uint                 `form:"media_ids,omitempty"`
	Media    []multipart.FileHeader `form:"media,omitempty" binding:"file"`
}
