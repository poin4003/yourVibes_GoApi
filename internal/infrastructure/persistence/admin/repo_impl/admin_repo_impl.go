package repo_impl

import (
	"context"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/admin/query"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/admin/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/admin/mapper"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"gorm.io/gorm"
	"time"
)

type rAdmin struct {
	db *gorm.DB
}

func NewAdminRepositoryImplement(db *gorm.DB) *rAdmin {
	return &rAdmin{db: db}
}

func (r *rAdmin) GetById(
	ctx context.Context,
	id uuid.UUID,
) (*entities.Admin, error) {
	var adminModel models.Admin
	if err := r.db.WithContext(ctx).
		First(&adminModel, id).
		Error; err != nil {
		return nil, err
	}

	return mapper.FromAdminModel(&adminModel), nil
}

func (r *rAdmin) GetStatusById(
	ctx context.Context,
	id uuid.UUID,
) (bool, error) {
	var adminStatus bool
	if err := r.db.WithContext(ctx).
		Model(&models.Admin{}).
		Select("status").
		First(&adminStatus, id).
		Error; err != nil {
		return false, err
	}
	return adminStatus, nil
}

func (r *rAdmin) CreateOne(
	ctx context.Context,
	entity *entities.Admin,
) (*entities.Admin, error) {
	adminModel := mapper.ToAdminModel(entity)

	if err := r.db.WithContext(ctx).
		Create(adminModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, adminModel.ID)
}

func (r *rAdmin) UpdateOne(
	ctx context.Context,
	id uuid.UUID,
	updateData *entities.AdminUpdate,
) (*entities.Admin, error) {
	updates := map[string]interface{}{}

	if updateData.FamilyName != nil {
		updates["family_name"] = *updateData.FamilyName
	}

	if updateData.Name != nil {
		updates["name"] = *updateData.Name
	}

	if updateData.PhoneNumber != nil {
		updates["phone_number"] = *updateData.PhoneNumber
	}

	if updateData.IdentityId != nil {
		updates["identity_id"] = *updateData.IdentityId
	}

	if updateData.Birthday != nil {
		updates["birthday"] = *updateData.Birthday
	}

	if updateData.Status != nil {
		updates["status"] = *updateData.Status
	}

	if updateData.Role != nil {
		updates["role"] = *updateData.Role
	}

	if updateData.Password != nil {
		updates["password"] = *updateData.Password
	}

	if err := r.db.WithContext(ctx).
		Model(&models.Admin{}).
		Where("id = ?", id).
		Updates(updates).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *rAdmin) GetOne(
	ctx context.Context,
	query interface{},
	args ...interface{},
) (*entities.Admin, error) {
	var adminModel models.Admin

	if err := r.db.WithContext(ctx).
		Model(&adminModel).
		Where(query, args...).
		First(&adminModel).
		Error; err != nil {
		return nil, err
	}

	return r.GetById(ctx, adminModel.ID)
}

func (r *rAdmin) GetMany(
	ctx context.Context,
	query *query.GetManyAdminQuery,
) ([]*entities.Admin, *response.PagingResponse, error) {
	var adminModels []models.Admin
	var total int64

	db := r.db.WithContext(ctx).Model(&models.Admin{})

	if query.Name != "" {
		db = db.Where("unaccent(family_name || ' ' || name) ILIKE unaccent(?)", "%"+query.Name+"%")
	}

	if query.Email != "" {
		db = db.Where("email = ?", query.Email)
	}

	if query.PhoneNumber != "" {
		db = db.Where("phone_number = ?", query.PhoneNumber)
	}

	if query.IdentityId != "" {
		db = db.Where("identity_id = ?", query.IdentityId)
	}

	if !query.Birthday.IsZero() {
		birthday := query.Birthday.Truncate(24 * time.Hour)
		db = db.Where("birthday = ?", birthday)
	}

	if !query.CreatedAt.IsZero() {
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("created_at = ?", createdAt)
	}

	if query.Status != nil {
		if *query.Status {
			db = db.Where("status = ?", true)
		} else {
			db = db.Where("status = ?", false)
		}
	}

	if query.Role != nil {
		if *query.Role {
			db = db.Where("role = ?", true)
		} else {
			db = db.Where("role = ?", false)
		}
	}

	if query.SortBy != "" {
		switch query.SortBy {
		case "id":
			if query.IsDescending {
				db = db.Order("id DESC")
			} else {
				db = db.Order("id ASC")
			}
		case "name":
			combinedName := "unaccent(family_name || ' ' name)"
			if query.IsDescending {
				db = db.Order(combinedName + "DESC")
			} else {
				db = db.Order(combinedName + "ASC")
			}
		case "email":
			if query.IsDescending {
				db = db.Order("email DESC")
			} else {
				db = db.Order("email ASC")
			}
		case "phone_number":
			if query.IsDescending {
				db = db.Order("phone_number DESC")
			} else {
				db = db.Order("phone_number ASC")
			}
		case "identity_id":
			if query.IsDescending {
				db = db.Order("identity_id DESC")
			} else {
				db = db.Order("identity_id ASC")
			}
		case "birthday":
			if query.IsDescending {
				db = db.Order("birthday DESC")
			} else {
				db = db.Order("birthday ASC")
			}
		case "created_at":
			if query.IsDescending {
				db = db.Order("created_at DESC")
			} else {
				db = db.Order("created_at ASC")
			}
		}
	}

	if err := db.Count(&total).
		Error; err != nil {
		return nil, nil, err
	}

	limit := query.Limit
	page := query.Page
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if err := db.Offset(offset).
		Limit(limit).
		Find(&adminModels).
		Error; err != nil {
		return nil, nil, err
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var adminsEntities []*entities.Admin
	for _, admin := range adminModels {
		adminEntity := mapper.FromAdminModel(&admin)
		adminsEntities = append(adminsEntities, adminEntity)
	}

	return adminsEntities, pagingResponse, nil
}

func (r *rAdmin) CheckAdminExistByEmail(
	ctx context.Context,
	email string,
) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).
		Model(&models.Admin{}).
		Where("email = ?", email).
		Count(&count).
		Error; err != nil {
	}

	return count > 0, nil
}
