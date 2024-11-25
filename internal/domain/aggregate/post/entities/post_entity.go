package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type Post struct {
	ID              uuid.UUID           `validate:"omitempty,uuid4"`
	UserId          uuid.UUID           `validate:"omitempty,uuid4"`
	User            *User               `validate:"omitempty"`
	ParentId        *uuid.UUID          `validate:"omitempty,uuid4"`
	ParentPost      *Post               `validate:"omitempty"`
	Content         string              `validate:"required"`
	LikeCount       int                 `validate:"omitempty"`
	CommentCount    int                 `validate:"omitempty"`
	Privacy         consts.PrivacyLevel `validate:"omitempty,oneof=public private friend_only"`
	Location        string              `validate:"omitempty"`
	IsAdvertisement bool                `validate:"omitempty"`
	Status          bool                `validate:"omitempty"`
	CreatedAt       time.Time           `validate:"omitempty"`
	UpdatedAt       time.Time           `validate:"omitempty,gtefield=CreatedAt"`
	Media           []*Media            `validate:"omitempty"`
}

type PostUpdate struct {
	Content         *string              `validate:"omitempty"`
	LikeCount       *int                 `validate:"omitempty"`
	CommentCount    *int                 `validate:"omitempty"`
	Privacy         *consts.PrivacyLevel `validate:"omitempty,oneof=public private friend_only"`
	Location        *string              `validate:"omitempty"`
	IsAdvertisement *bool                `validate:"omitempty"`
	Status          *bool                `validate:"omitempty"`
	UpdatedAt       *time.Time           `validate:"omitempty,gtefield=CreatedAt"`
}

func (p *Post) ValidatePost() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Content, validation.Min(2)),
		validation.Field(&p.Privacy, validation.In(consts.PRIVATE, consts.PUBLIC, consts.FRIEND_ONLY)),
		validation.Field(&p.LikeCount, validation.Min(0)),
		validation.Field(&p.CommentCount, validation.Min(0)),
		validation.Field(&p.UpdatedAt, validation.Min(p.CreatedAt)),
	)
}

func (p *PostUpdate) ValidatePostUpdate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Content, validation.Min(2)),
		validation.Field(&p.Privacy, validation.In(consts.PRIVATE, consts.PUBLIC, consts.FRIEND_ONLY)),
		validation.Field(&p.LikeCount, validation.Min(0)),
		validation.Field(&p.CommentCount, validation.Min(0)),
	)
}

func NewPost(
	userId uuid.UUID,
	content string,
	privacy consts.PrivacyLevel,
	location string,
) (*Post, error) {
	post := &Post{
		ID:              uuid.New(),
		UserId:          userId,
		Content:         content,
		LikeCount:       0,
		CommentCount:    0,
		Privacy:         privacy,
		Location:        location,
		IsAdvertisement: false,
		Status:          true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := post.ValidatePost(); err != nil {
		return nil, err
	}

	return post, nil
}

func NewPostForShare(
	userId uuid.UUID,
	content string,
	privacy consts.PrivacyLevel,
	location string,
	parentId *uuid.UUID,
) (*Post, error) {
	post := &Post{
		ID:              uuid.New(),
		UserId:          userId,
		Content:         content,
		LikeCount:       0,
		CommentCount:    0,
		ParentId:        parentId,
		Privacy:         privacy,
		Location:        location,
		IsAdvertisement: false,
		Status:          true,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := post.ValidatePost(); err != nil {
		return nil, err
	}

	return post, nil
}
