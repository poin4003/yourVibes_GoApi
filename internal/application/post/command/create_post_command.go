package command

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/common"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type CreatePostCommand struct {
	UserId   uuid.UUID
	Content  string
	Privacy  consts.PrivacyLevel
	Location string
	Media    []multipart.FileHeader
}

// Message from python server
type ContentModerationResult struct {
	Label        string `json:"label"`
	CensoredText string `json:"censored_text"`
}

type MediaModerationResult struct {
	Label string `json:"label"`
}

type PostModerationResult struct {
	PostID  uuid.UUID               `json:"post_id"`
	Content ContentModerationResult `json:"content"`
	Media   MediaModerationResult   `json:"media"`
}

type ApprovePostCommand struct {
	PostId       uuid.UUID
	CensoredText string
}

type RejectPostCommand struct {
	PostId uuid.UUID
	Label  string
}

type CreatePostCommandResult struct {
	Post *common.PostResult
}
