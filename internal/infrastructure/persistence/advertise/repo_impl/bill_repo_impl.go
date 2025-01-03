package repo_impl

import (
	"context"
	"time"

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

func (r *rBill) GetMonthlyRevenue(ctx context.Context, date time.Time) ([]string, []int64, error) {
	months := make([]string, 13)
	revenueByMonth := make([]int64, 13)

	startDate := date.AddDate(0, -12, 0)

	currentDate := startDate
	for i := 0; i < 13; i++ {
		months[i] = currentDate.Format("01/2006")
		currentDate = currentDate.AddDate(0, 1, 0)
	}

	var revenues []entities.Revenue
	if err := r.db.WithContext(ctx).
		Model(&models.Bill{}).
		Select("DATE_TRUNC('month', created_at) AS month, SUM(price) AS total").
		Where("created_at >= ?", startDate).
		Where("status = true").
		Group("month").
		Order("month").
		Scan(&revenues).
		Error; err != nil {
		return nil, nil, err
	}

	monthIndexMap := make(map[string]int)
	for i, month := range months {
		monthIndexMap[month] = i
	}

	for _, revenue := range revenues {
		monthStr := revenue.Month.Format("01/2006")
		if idx, exists := monthIndexMap[monthStr]; exists {
			revenueByMonth[idx] = revenue.Total
		}

	}

	return months, revenueByMonth, nil
}

func (r *rBill) GetRevenueForMonth(ctx context.Context, date time.Time) (int64, error) {
	var total int64
	startDate := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	endDate := startDate.AddDate(0, 1, 0)

	if err := r.db.WithContext(ctx).
		Model(&models.Bill{}).
		Where("created_at >= ? AND created_at < ?", startDate, endDate).
		Where("status = true").
		Select("COALESCE(CAST(SUM(price) AS INT), 0) AS total").
		Scan(&total).
		Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *rBill) GetRevenueForDay(ctx context.Context, date time.Time) (int64, error) {
	var total int64
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	if err := r.db.WithContext(ctx).
		Model(&models.Bill{}).
		Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
		Where("status = true").
		Select("COALESCE(CAST(SUM(price) AS INT), 0) AS total").
		Scan(&total).
		Error; err != nil {
		return 0, err
	}

	return total, nil
}
