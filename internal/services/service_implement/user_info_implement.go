package service_implement

import (
	"context"
	"errors"
	"fmt"

	"github.com/poin4003/yourVibes_GoApi/internal/model"
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"gorm.io/gorm"
)

type sUserInfo struct {
	repo repository.IUserRepository
}

func NewUserInfoImplement(repo repository.IUserRepository) *sUserInfo {
	return &sUserInfo{repo: repo}
}

func (s *sUserInfo) GetInfoByUserId(ctx context.Context, id string) (*model.User, error) {
	userFound, err := s.repo.GetUser(ctx, "id = ?", id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userFound, nil
}

func (s *sUserInfo) GetUsersByName(ctx context.Context, keyword string, limit, page int) ([]*model.User, int64, error) {
	query := "unaccent(family_name || ' ' || name) ILIKE unaccent(?)"
	args := []interface{}{"%" + keyword + "%"}

	users, total, err := s.repo.GetManyUser(ctx, limit, page, query, args...)
	if err != nil {
		return nil, 0, err
	}

	if users == nil {
		users = []*model.User{}
	}

	return users, total, nil
}
