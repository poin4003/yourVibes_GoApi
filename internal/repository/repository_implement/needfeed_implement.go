package repository_implement

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/query_object"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
)

type rNewFeed struct {
	db *gorm.DB
}

func NewNewFeedRepositoryImplement(db *gorm.DB) *rNewFeed {
	return &rNewFeed{db: db}
}

func (r *rNewFeed) CreateManyNewFeed(
	ctx context.Context,
	newFeed *model.NewFeed,
) error {
	return nil
}

func (r *rNewFeed) DeleteNewFeed(
	ctx context.Context,
	newFeed *model.NewFeed,
) error {
	return nil
}

func (r *rNewFeed) GetManyNewFeed(
	ctx context.Context,
	userId uuid.UUID,
	query *query_object.NewFeedQueryObject,
) ([]*model.Post, *response.PagingResponse, error) {
	return nil, nil, nil
}
