package entities

import (
	"github.com/go-playground/validator/v10"
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

func (p *Post) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

func (p *PostUpdate) ValidatePostUpdate() error {
	validate := validator.New()
	return validate.Struct(p)
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
	if err := post.Validate(); err != nil {
		return nil, err
	}

	return post, nil
}

func NewPostUpdate(
	updateData *PostUpdate,
) (*PostUpdate, error) {
	postUpdate := &PostUpdate{
		Content:         updateData.Content,
		Privacy:         updateData.Privacy,
		Location:        updateData.Location,
		LikeCount:       updateData.LikeCount,
		CommentCount:    updateData.CommentCount,
		IsAdvertisement: updateData.IsAdvertisement,
		Status:          updateData.Status,
	}
	if err := postUpdate.ValidatePostUpdate(); err != nil {
		return nil, err
	}

	return postUpdate, nil
}
