package entities

import (
	"fmt"
	"time"
	"unicode/utf8"

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
		validation.Field(&c.Content, validation.By(func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid content type")
			}

			length := utf8.RuneCountInString(str)
			if length < 1 || length > 500 {
				return fmt.Errorf("content length must be between 1 and 500 characters, but got %d", length)
			}
			return nil
		})),
	)
}

func (cu *CommentUpdate) ValidateCommentUpdate() error {
	return validation.ValidateStruct(cu,
		validation.Field(&cu.Content, validation.By(func(value interface{}) error {
			str, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid content type")
			}

			length := utf8.RuneCountInString(str)
			if length < 1 || length > 500 {
				return fmt.Errorf("content length must be between 1 and 500 characters, but got %d", length)
			}
			return nil
		})),
	)
}

func NewComment(
	postId uuid.UUID,
	userId uuid.UUID,
	parentId *uuid.UUID,
	content string,
) (*Comment, error) {
	comment := &Comment{
		ID:              uuid.New(),
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
