package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/advertise/mapper"
	"gorm.io/gorm"
)

type rBill struct {
	db *gorm.DB
}

func NewBillRepositoryImplement(db *gorm.DB) *rBill {
	return &rBill{db: db}
}

func (r *rBill) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Bill, error) {
	var billModel models.Bill
	if err := r.db.WithContext(ctx).
		Preload("Advertise").
		First(&billModel, id).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromBillModel(&billModel), nil
}

func (r *rBill) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Bill, error) {
	var billModel models.Bill

	if err := r.db.WithContext(ctx).
		Model(&billModel).
		Where(query, args...).
		First(&billModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, billModel.ID)
}

func (r *rBill) CreateOne(
	ctx context.Context,
	entity *entities.Bill,
) (*entities.Bill, error) {
	billModel := mapper.ToBillModel(entity)

	if err := r.db.WithContext(ctx).
		Create(billModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, billModel.ID)
}

func (r *rBill) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.BillUpdate,
) (*entities.Bill, error) {
	updates := map[string]interface{}{}

	if updateData.Price != nil {
		updates["price"] = *updateData.Price
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Bill{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *rBill) DeleteOne(
	ctx context.Context,
	id uuid.UUID,
) error {
	if err := r.db.WithContext(ctx).
		Delete(&models.Bill{}, id).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rBill) CheckExists(
	ctx context.Context,
	postId uuid.UUID,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&models.Bill{}).
		Joins("JOIN advertises ON advertises.id = bills.advertise_id").
		Where("advertises.post_id = ?", postId).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
