package repositories

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/application/user/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/notification/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
)

type (
	INotificationRepository interface {
		CreateOne(ctx context.Context, entity *entities.Notification) (*entities.Notification, error)
		CreateMany(ctx context.Context, entity []*entities.Notification) ([]*entities.Notification, error)
		UpdateOne(ctx context.Context, id uint, updateData *entities.NotificationUpdate) (*entities.Notification, error)
		UpdateMany(ctx context.Context, condition map[string]interface{}, updateData map[string]interface{}) error
		DeleteOne(ctx context.Context, id uint) (*entities.Notification, error)
		GetById(ctx context.Context, id uint) (*entities.Notification, error)
		GetOne(ctx context.Context, query interface{}, args ...interface{}) (*entities.Notification, error)
		GetMany(ctx context.Context, query *query.GetManyNotificationQuery) ([]*entities.Notification, *response.PagingResponse, error)
	}
)

var (
	localNotification INotificationRepository
)

func Notification() INotificationRepository {
	if localNotification == nil {
		panic("repository_implement localNotification not found for interface INotification")
	}

	return localNotification
}

func InitNotificationRepository(i INotificationRepository) {
	localNotification = i
}
