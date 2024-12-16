package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	ID              uuid.UUID  `validate:"omitempty,uuid4"`
	PostId          uuid.UUID  `validate:"required,uuid4"`
	UserId          uuid.UUID  `validate:"required,uuid4"`
	User            *User      `validate:"omitempty"`
	ParentId        *uuid.UUID `validate:"omitempty,uuid4"`
	Content         string     `validate:"required,min=2,max=500"`
	LikeCount       int        `validate:"omitempty"`
	RepCommentCount int        `validate:"omitempty"`
	CommentLeft     int        `validate:"omitempty"`
	CommentRight    int        `validate:"omitempty"`
	CreatedAt       time.Time  `validate:"omitempty"`
	UpdatedAt       time.Time  `validate:"omitempty,gtefield=CreatedAt"`
	Status          bool       `validate:"omitempty"`
}

type CommentUpdate struct {
	Content         *string    `validate:"omitempty,min=2,max=500"`
	LikeCount       *int       `validate:"omitempty"`
	RepCommentCount *int       `validate:"omitempty"`
	UpdatedAt       *time.Time `validate:"omitempty,gtefield=CreatedAt"`
	Status          *bool      `validate:"omitempty"`
}

type CommentForReport struct {
	ID              uuid.UUID
	PostId          uuid.UUID
	UserId          uuid.UUID
	User            *UserForReport
	ParentId        *uuid.UUID
	Content         string
	LikeCount       int
	RepCommentCount int
	CommentLeft     int
	CommentRight    int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Status          bool
}

func (c *Comment) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (cu *CommentUpdate) ValidateCommentUpdate() error {
	validate := validator.New()
	return validate.Struct(cu)
}

func NewComment(
	postId uuid.UUID,
	userId uuid.UUID,
	parentId *uuid.UUID,
	content string,
	commentLeft int,
	commentRight int,
) (*Comment, error) {
	comment := &Comment{
		ID:              uuid.New(),
		PostId:          postId,
		UserId:          userId,
		ParentId:        parentId,
		Content:         content,
		LikeCount:       0,
		RepCommentCount: 0,
		CommentLeft:     commentLeft,
		CommentRight:    commentRight,
		Status:          true,
	}
	if err := comment.Validate(); err != nil {
		return nil, err
	}

	return comment, nil
}
