package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type Comment struct {
	ID              uuid.UUID
	PostId          uuid.UUID
	UserId          uuid.UUID
	User            *User
	ParentId        *uuid.UUID
	Content         string
	LikeCount       int
	RepCommentCount int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          bool
}

type CommentUpdate struct {
	Content         *string
	LikeCount       *int
	RepCommentCount *int
	UpdatedAt       *time.Time
	Status          *bool
}

func (c *Comment) ValidateComment() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.PostId, validation.Required),
		validation.Field(&c.UserId, validation.Required),
		validation.Field(&c.CreatedAt, validation.Min(c.CreatedAt)),
		validation.Field(&c.Content, validation.RuneLength(1, 500)),
	)
}

func (cu *CommentUpdate) ValidateCommentUpdate() error {
	return validation.ValidateStruct(cu,
		validation.Field(&cu.Content, validation.RuneLength(1, 500)),
	)
}

func NewComment(
	commentId uuid.UUID,
	postId uuid.UUID,
	userId uuid.UUID,
	parentId *uuid.UUID,
	content string,
) (*Comment, error) {
	comment := &Comment{
		ID:              commentId,
		PostId:          postId,
		UserId:          userId,
		ParentId:        parentId,
		Content:         content,
		LikeCount:       0,
		RepCommentCount: 0,
		Status:          true,
	}
	if err := comment.ValidateComment(); err != nil {
		return nil, err
	}

	return comment, nil
}
