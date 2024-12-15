package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"time"
)

type Post struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *User
	ParentId        *uuid.UUID
	ParentPost      *Post
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement bool
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Media           []*Media
}

type PostWithLiked struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *User
	ParentId        *uuid.UUID
	ParentPost      *Post
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement bool
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Media           []*Media
	IsLiked         bool
}

type PostUpdate struct {
	Content         *string
	LikeCount       *int
	CommentCount    *int
	Privacy         *consts.PrivacyLevel
	Location        *string
	IsAdvertisement *bool
	Status          *bool
	UpdatedAt       *time.Time
}

type PostForReport struct {
	ID              uuid.UUID
	UserId          uuid.UUID
	User            *User
	ParentId        *uuid.UUID
	ParentPost      *PostForReport
	Content         string
	LikeCount       int
	CommentCount    int
	Privacy         consts.PrivacyLevel
	Location        string
	IsAdvertisement bool
	Status          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Media           []*Media
}

func (p *Post) ValidatePost() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Content, validation.Length(2, 1000)),
		validation.Field(&p.Privacy, validation.In(consts.PRIVATE, consts.PUBLIC, consts.FRIEND_ONLY)),
		validation.Field(&p.LikeCount, validation.Min(0)),
		validation.Field(&p.CommentCount, validation.Min(0)),
		validation.Field(&p.UpdatedAt, validation.Min(p.CreatedAt)),
	)
}

func (p *PostUpdate) ValidatePostUpdate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.Content, validation.Length(2, 1000)),
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
