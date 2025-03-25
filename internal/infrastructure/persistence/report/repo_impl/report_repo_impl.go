package repo_impl

import (
	"context"
	"errors"
	"time"

	"github.com/poin4003/yourVibes_GoApi/internal/application/report/query"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/report/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/models"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/report/mapper"
	"gorm.io/gorm"
)

type rReport struct {
	db *gorm.DB
}

func NewReportRepositoryImplement(db *gorm.DB) *rReport {
	return &rReport{db: db}
}

func (r *rReport) GetManyUserReport(
	ctx context.Context,
	query *query.GetManyReportQuery,
) ([]*entities.UserReportEntity, *response.PagingResponse, error) {
	var userReportModels []*models.UserReport
	var total int64

	db := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Joins("JOIN reports ON reports.id = user_reports.report_id").
		Where("reports.deleted_at IS NULL").
		Preload("Report").
		Preload("Report.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		}).
		Preload("ReportedUser", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		}).
		Preload("Report.Admin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		})

	if query.Reason != "" {
		db = db.Where("reports.reason = ?", query.Reason)
	}

	if query.UserEmail != "" {
		db = db.Joins("JOIN users u ON u.id = reports.user_id").
			Where("u.email = ?", query.UserEmail)
	}

	if query.ReportedUserEmail != "" {
		db = db.Joins("JOIN users lu ON lu.id = user_reports.reported_user_id").
			Where("lu.email = ?", query.ReportedUserEmail)
	}

	if query.AdminEmail != "" {
		db = db.Joins("LEFT JOIN admins a ON a.id = reports.admin_id").
			Where("a.email = ?", query.AdminEmail)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("reports.created_at >= ?", query.FromDate)
	}

	if !query.ToDate.IsZero() {
		db = db.Where("reports.created_at <= ?", query.ToDate)
	}

	if !query.CreatedAt.IsZero() {
		createAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("reports.created_at = ?", createAt)
	}

	if query.Status != nil {
		db = db.Where("reports.status = ?", *query.Status)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "user_id":
			sortColumn = "reports.user_id"
		case "reported_user_id":
			sortColumn = "user_reports.reported_user_id"
		case "admin_id":
			sortColumn = "reports.admin_id"
		case "reason":
			sortColumn = "reports.reason"
		case "created_at":
			sortColumn = "reports.created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
			}
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
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
		Find(&userReportModels).
		Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var userReportEntities []*entities.UserReportEntity
	for _, userReportModel := range userReportModels {
		userReportEntity := mapper.FromUserReportModel(userReportModel)
		userReportEntities = append(userReportEntities, userReportEntity)
	}

	return userReportEntities, pagingResponse, nil
}

func (r *rReport) GetManyPostReport(
	ctx context.Context,
	query *query.GetManyReportQuery,
) ([]*entities.PostReportEntity, *response.PagingResponse, error) {
	var postReportModels []*models.PostReport
	var total int64

	db := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Joins("JOIN reports ON reports.id = post_reports.report_id").
		Where("reports.deleted_at IS NULL").
		Preload("Report").
		Preload("Report.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		}).
		Preload("Report.Admin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		})

	if query.Reason != "" {
		db = db.Where("reports.reason = ?", query.Reason)
	}

	if query.UserEmail != "" {
		db = db.Joins("JOIN users u ON u.id = reports.user_id").
			Where("u.email = ?", query.UserEmail)
	}

	if query.AdminEmail != "" {
		db = db.Joins("LEFT JOIN admins a ON a.id = reports.admin_id").
			Where("a.email = ?", query.AdminEmail)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("reports.created_at >= ?", query.FromDate)
	}

	if !query.ToDate.IsZero() {
		db = db.Where("reports.created_at <= ?", query.ToDate)
	}

	if !query.CreatedAt.IsZero() {
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("reports.created_at = ?", createdAt)
	}

	if query.Status != nil {
		db = db.Where("reports.status = ?", query.Status)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "user_id":
			sortColumn = "reports.user_id"
		case "reported_post_id":
			sortColumn = "post_reports.reported_post_id"
		case "admin_id":
			sortColumn = "reports.admin_id"
		case "reason":
			sortColumn = "reports.reason"
		case "created_at":
			sortColumn = "reports.created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
			}
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
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
		Find(&postReportModels).
		Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var postReportEntities []*entities.PostReportEntity
	for _, postReportModel := range postReportModels {
		postReportEntity := mapper.FromPostReportModel(postReportModel)
		postReportEntities = append(postReportEntities, postReportEntity)
	}

	return postReportEntities, pagingResponse, nil
}

func (r *rReport) GetManyCommentReport(
	ctx context.Context,
	query *query.GetManyReportQuery,
) ([]*entities.CommentReportEntity, *response.PagingResponse, error) {
	var commentReportModels []*models.CommentReport
	var total int64

	db := r.db.WithContext(ctx).
		Model(&models.CommentReport{}).
		Joins("JOIN reports ON reports.id = comment_reports.report_id").
		Where("reports.deleted_at IS NULL").
		Preload("Report").
		Preload("Report.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		}).
		Preload("Report.Admin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		})

	if query.Reason != "" {
		db = db.Where("reports.reason = ?", query.Reason)
	}

	if query.UserEmail != "" {
		db = db.Joins("JOIN users u ON u.id = reports.user_id").
			Where("u.email = ?", query.UserEmail)
	}

	if query.AdminEmail != "" {
		db = db.Joins("LEFT JOIN admins a ON a.id = reports.admin_id").
			Where("a.email = ?", query.AdminEmail)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("reports.created_at >= ?", query.FromDate)
	}

	if !query.ToDate.IsZero() {
		db = db.Where("reports.created_at <= ?", query.ToDate)
	}

	if !query.CreatedAt.IsZero() {
		createdAt := query.CreatedAt.Truncate(24 * time.Hour)
		db = db.Where("reports.created_at = ?", createdAt)
	}

	if query.Status != nil {
		db = db.Where("reports.status = ?", query.Status)
	}

	if query.SortBy != "" {
		sortColumn := ""
		switch query.SortBy {
		case "user_id":
			sortColumn = "reports.user_id"
		case "reported_comment_id":
			sortColumn = "comment_reports.reported_comment_id"
		case "admin_id":
			sortColumn = "reports.admin_id"
		case "reason":
			sortColumn = "reports.reason"
		case "created_at":
			sortColumn = "reports.created_at"
		}

		if sortColumn != "" {
			if query.IsDescending {
				db = db.Order(sortColumn + " DESC")
			} else {
				db = db.Order(sortColumn + " ASC")
			}
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
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
		Find(&commentReportModels).
		Error; err != nil {
		return nil, nil, response.NewServerFailedError(err.Error())
	}

	pagingResponse := &response.PagingResponse{
		Limit: limit,
		Page:  page,
		Total: total,
	}

	var commentReportEntities []*entities.CommentReportEntity
	for _, commentReportModel := range commentReportModels {
		commentReportEntity := mapper.FromCommentReportModel(commentReportModel)
		commentReportEntities = append(commentReportEntities, commentReportEntity)
	}

	return commentReportEntities, pagingResponse, nil
}

func (r *rReport) GetUserReportById(
	ctx context.Context,
	reportId uuid.UUID,
) (*entities.UserReportEntity, error) {
	var userReportModel models.UserReport

	if err := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Where("report_id = ?", reportId).
		Preload("Report.User").
		Preload("Report.Admin").
		Preload("ReportedUser").
		First(&userReportModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	if userReportModel.Report == nil || userReportModel.Report.DeletedAt.Valid {
		return nil, response.NewDataNotFoundError("report has been deleted")
	}

	return mapper.FromUserReportModel(&userReportModel), nil
}

func (r *rReport) GetPostReportById(
	ctx context.Context,
	reportId uuid.UUID,
) (*entities.PostReportEntity, error) {
	var postReportModel models.PostReport

	if err := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Where("report_id = ?", reportId).
		Preload("Report").
		Preload("Report.User").
		Preload("Report.Admin").
		Preload("ReportedPost").
		Preload("ReportedPost.Media").
		Preload("ReportedPost.User").
		Preload("ReportedPost.ParentPost.Media").
		Preload("ReportedPost.ParentPost.User").
		First(&postReportModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewDataNotFoundError(err.Error())
	}

	if postReportModel.Report == nil || postReportModel.Report.DeletedAt.Valid {
		return nil, response.NewDataNotFoundError("report has been deleted")
	}

	result := mapper.FromPostReportModel(&postReportModel)
	return result, nil
}

func (r *rReport) GetCommentReportById(
	ctx context.Context,
	reportId uuid.UUID,
) (*entities.CommentReportEntity, error) {
	var commentReportModel models.CommentReport

	if err := r.db.WithContext(ctx).
		Model(&models.CommentReport{}).
		Where("report_id = ?", reportId).
		Preload("Report.User").
		Preload("Report.Admin").
		Preload("ReportedComment").
		Preload("ReportedComment.User").
		Preload("ReportedComment.Post.User").
		Preload("ReportedComment.Post.Media").
		Preload("ReportedComment.Post.ParentPost.Media").
		Preload("ReportedComment.Post.ParentPost.User").
		First(&commentReportModel).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}

	if commentReportModel.Report == nil || commentReportModel.Report.DeletedAt.Valid {
		return nil, response.NewDataNotFoundError("report has been deleted")
	}

	return mapper.FromCommentReportModel(&commentReportModel), nil
}

func (r *rReport) CreatePostReport(
	ctx context.Context,
	entity *entities.PostReportEntity,
) error {
	postReportModel := mapper.ToPostReportModel(entity)
	// 1. Check exists
	postReportExists, _ := r.checkPostReportExists(ctx, entity.Report.UserId, entity.ReportedPostId)
	if postReportExists {
		return response.NewCustomError(response.ErrDataHasAlreadyExist, "post report already exits")
	}

	// 2. Check post exists
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", entity.ReportedPostId).
		Count(&count).Error; err != nil {
	}
	if count == 0 {
		return response.NewDataNotFoundError("post not exits!")
	}

	// 3. Create report and user report
	if err := r.db.WithContext(ctx).
		Create(&postReportModel).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (r *rReport) CreateUserReport(
	ctx context.Context,
	entity *entities.UserReportEntity,
) error {
	userReportModel := mapper.ToUserReportModel(entity)
	// 1. Check exits
	userReportExits, _ := r.checkUserReportExists(ctx, entity.Report.UserId, entity.ReportedUserId)
	if userReportExits {
		return response.NewCustomError(response.ErrDataHasAlreadyExist, "user report already exists")
	}

	// 2. Check user exists
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", entity.ReportedUserId).
		Count(&count).Error; err != nil {
	}
	if count == 0 {
		return response.NewDataNotFoundError("user not exits!")
	}

	// 3. Create report and comment report
	if err := r.db.WithContext(ctx).
		Create(&userReportModel).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (r *rReport) CreateCommentReport(
	ctx context.Context,
	entity *entities.CommentReportEntity,
) error {
	commentReportModel := mapper.ToCommentReportModel(entity)
	// 1. Check exists
	commentReportExists, _ := r.checkCommentReportExists(ctx, entity.Report.UserId, entity.ReportedCommentId)
	if commentReportExists {
		return response.NewCustomError(response.ErrDataHasAlreadyExist, "comment report already exists")
	}

	// 2. Check comment exists
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", entity.ReportedCommentId).
		Count(&count).Error; err != nil {
	}
	if count == 0 {
		return response.NewDataNotFoundError("comment not exits!")
	}

	// 3. Create report and comment report
	if err := r.db.WithContext(ctx).
		Create(&commentReportModel).
		Error; err != nil {
		return err
	}

	return nil
}

func (r *rReport) HandleUserReport(
	ctx context.Context,
	reportId, adminId uuid.UUID,
) (*entities.UserForReport, error) {
	userFound := &models.User{}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get user report
		userReportFound, err := r.getUserReportById(ctx, tx, reportId)
		if err != nil {
			return err
		}

		// 2. Check user exits
		if err = tx.WithContext(ctx).
			Select("id, family_name, name, avatar_url").
			First(&userFound, userReportFound.ReportedUserId).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 3. Check report status
		if userReportFound.Report.Status {
			return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
		}

		// 4. Update reported user status
		if err = r.updateUserStatus(ctx, tx, userReportFound.ReportedUserId, false); err != nil {
			return err
		}

		// 5. Update reportedUser posts status
		if err = r.updateUserPostStatus(ctx, tx, userReportFound.ReportedUserId, false); err != nil {
			return err
		}

		// 6. Update reported User comments status
		if err = r.updateUserCommentStatus(ctx, tx, userReportFound.ReportedUserId, false); err != nil {
			return err
		}

		// 7. Update report status
		if err = r.updateReportAdminIdAndStatus(ctx, tx, reportId, adminId, true); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.FromUserModel(userFound), nil
}

func (r *rReport) HandlePostReport(
	ctx context.Context,
	reportId, adminId uuid.UUID,
) (*entities.PostForReport, error) {
	postFound := &models.Post{}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get post report
		postReportFound, err := r.getPostReportById(ctx, tx, reportId)
		if err != nil {
			return err
		}

		// 2. Check post exists
		if err = tx.WithContext(ctx).
			Model(postFound).
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, family_name, name, avatar_url")
			}).
			First(&postFound, postReportFound.ReportedPostId).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 3. Check report status
		if postReportFound.Report.Status {
			return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
		}

		// 4. Update reported post status
		if err = r.updatePostStatus(ctx, tx, postReportFound.ReportedPostId, false); err != nil {
			return err
		}

		// 5. Update report status
		if err = r.updateReportAdminIdAndStatus(ctx, tx, reportId, adminId, true); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.FromPostModel(postFound), nil
}

func (r *rReport) HandleCommentReport(
	ctx context.Context,
	reportId, adminId uuid.UUID,
) (*entities.CommentForReport, error) {
	commentFound := &models.Comment{}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Get comment report
		commentReportFound, err := r.getCommentReportById(ctx, tx, reportId)
		if err != nil {
			return err
		}

		// 2. Get comment to check
		if err = tx.WithContext(ctx).
			Model(commentFound).
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, family_name, name, avatar_url")
			}).
			First(commentFound, commentReportFound.ReportedCommentId).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 3. Check report status
		if commentReportFound.Report.Status {
			return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
		}

		// 4. Update reported comment status
		if err = r.updateCommentStatus(ctx, tx, commentReportFound.ReportedCommentId, false); err != nil {
			return err
		}

		// 5. Update report status
		if err = r.updateReportAdminIdAndStatus(ctx, tx, reportId, adminId, true); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.FromCommentModel(commentFound), nil
}

func (r *rReport) ActivateUser(
	ctx context.Context,
	reportId uuid.UUID,
) (*entities.UserForReport, error) {
	userFound := &models.User{}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check user exists
		if err := tx.WithContext(ctx).
			Model(userFound).
			Joins("JOIN user_reports ON user_reports.reported_user_id = users.id").
			Where("user_reports.report_id = ?", reportId).
			Select("users.id, users.status, users.email, users.name, users.family_name").
			First(userFound).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 2. Check if user is already activate
		if userFound.Status {
			return response.NewCustomError(response.ErrCodeUserIsAlreadyActivated)
		}

		// 3. Update reported user status
		if err := r.updateUserStatus(ctx, tx, userFound.ID, true); err != nil {
			return err
		}

		// 4. Update reportedUser post status
		if err := tx.WithContext(ctx).Model(&models.Post{}).
			Where("user_id = ?", userFound.ID).
			Update("status", true).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// 5. Update reportedUser comment status
		if err := tx.WithContext(ctx).Model(&models.Comment{}).
			Where("user_id = ?", userFound.ID).
			Update("status", true).
			Error; err != nil {
			return response.NewServerFailedError(err.Error())
		}

		// 6. Delete report
		result := tx.WithContext(ctx).
			Where("id IN (?)",
				tx.Model(&models.UserReport{}).
					Select("report_id").
					Where("reported_user_id = ? AND report_id = ?", userFound.ID, reportId),
			).
			Delete(&models.Report{})

		if result.Error != nil {
			return response.NewServerFailedError(result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return response.NewDataNotFoundError("no report found to delete for this post and report id")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.FromUserModel(userFound), nil
}

func (r *rReport) ActivatePost(
	ctx context.Context,
	reportId uuid.UUID,
) (*entities.PostForReport, error) {
	postFound := &models.Post{}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check post exists
		if err := tx.WithContext(ctx).
			Model(postFound).
			Joins("JOIN post_reports ON post_reports.reported_post_id = posts.id").
			Where("post_reports.report_id = ?", reportId).
			Select("posts.id, posts.status, posts.user_id, posts.is_advertisement").
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, family_name, name, avatar_url")
			}).
			First(postFound).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 2. Check if post already activate
		if postFound.Status {
			return response.NewCustomError(response.ErrCodePostIsAlreadyActivated)
		}

		// 3. Update post status
		if err := r.updatePostStatus(ctx, tx, postFound.ID, true); err != nil {
			return err
		}

		// 4. Delete all related post reports
		result := tx.WithContext(ctx).
			Where("id IN (?)",
				tx.Model(&models.PostReport{}).
					Select("report_id").
					Where("reported_post_id = ? AND report_id = ?", postFound.ID, reportId),
			).
			Delete(&models.Report{})

		if result.Error != nil {
			return response.NewServerFailedError(result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return response.NewDataNotFoundError("no report found to delete for this post and report id")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.FromPostModel(postFound), nil
}

func (r *rReport) ActivateComment(
	ctx context.Context,
	reportId uuid.UUID,
) (*entities.CommentForReport, error) {
	commentFound := &models.Comment{}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Check comment exists
		if err := tx.WithContext(ctx).
			Model(commentFound).
			Joins("JOIN comment_reports ON comment_reports.reported_comment_id = comments.id").
			Where("comment_reports.report_id = ?", reportId).
			Select("comments.id, comments.status, comments.user_id").
			Preload("User", func(db *gorm.DB) *gorm.DB {
				return db.Select("id, family_name, name, avatar_url")
			}).
			First(commentFound).
			Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return response.NewDataNotFoundError(err.Error())
			}
			return response.NewServerFailedError(err.Error())
		}

		// 2. Check if comment already activate
		if commentFound.Status {
			return response.NewCustomError(response.ErrCodeCommentIsAlreadyActivated)
		}

		// 3. Update comment status
		if err := r.updateCommentStatus(ctx, tx, commentFound.ID, true); err != nil {
			return err
		}

		// 4. Delete all related comment reports
		result := tx.WithContext(ctx).
			Where("id IN (?)",
				tx.Model(&models.CommentReport{}).
					Select("report_id").
					Where("reported_comment_id = ? AND report_id = ?", commentFound.ID, reportId),
			).
			Delete(&models.Report{})

		if result.Error != nil {
			return response.NewServerFailedError(result.Error.Error())
		}
		if result.RowsAffected == 0 {
			return response.NewDataNotFoundError("no report found to delete for this post and report id")
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return mapper.FromCommentModel(commentFound), nil
}

func (r *rReport) DeleteReportById(
	ctx context.Context,
	reportId uuid.UUID,
) error {
	// 1. Get report
	reportModel := &models.Report{}
	if err := r.db.WithContext(ctx).
		Model(reportModel).
		First(reportModel, "id = ?", reportId).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.NewDataNotFoundError(err.Error())
		}
		return response.NewServerFailedError(err.Error())
	}

	// 2. Check if report is already handle (can't delete handled report)
	if reportModel.Status {
		return response.NewCustomError(response.ErrCodeReportIsAlreadyHandled)
	}

	// 3. Delete report
	if err := r.db.WithContext(ctx).
		Delete(&models.Report{}, reportId).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}

func (r *rReport) getUserReportById(
	ctx context.Context,
	tx *gorm.DB,
	reportId uuid.UUID,
) (*models.UserReport, error) {
	userReport := &models.UserReport{}
	if err := tx.WithContext(ctx).
		Model(userReport).
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, status")
		}).
		First(userReport, "report_id = ?", reportId).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}
	return userReport, nil
}

func (r *rReport) getPostReportById(
	ctx context.Context,
	tx *gorm.DB,
	reportId uuid.UUID,
) (*models.PostReport, error) {
	postReport := &models.PostReport{}
	if err := tx.WithContext(ctx).
		Model(postReport).
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, status")
		}).
		First(postReport, "report_id = ?", reportId).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}
	return postReport, nil
}

func (r *rReport) getCommentReportById(
	ctx context.Context,
	tx *gorm.DB,
	reportId uuid.UUID,
) (*models.CommentReport, error) {
	commentReport := &models.CommentReport{}
	if err := tx.WithContext(ctx).
		Model(commentReport).
		Preload("Report", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, status")
		}).
		First(commentReport, "report_id = ?", reportId).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, response.NewDataNotFoundError(err.Error())
		}
		return nil, response.NewServerFailedError(err.Error())
	}
	return commentReport, nil
}

func (r *rReport) updateUserStatus(
	ctx context.Context,
	tx *gorm.DB,
	userId uuid.UUID,
	status bool,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Update("status", status).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (r *rReport) updatePostStatus(
	ctx context.Context,
	tx *gorm.DB,
	postId uuid.UUID,
	status bool,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.Post{}).
		Where("id = ?", postId).
		Update("status", status).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (r *rReport) updateUserPostStatus(
	ctx context.Context,
	tx *gorm.DB,
	userId uuid.UUID,
	status bool,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.Post{}).
		Where("user_id = ?", userId).
		Update("status", status).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (r *rReport) updateCommentStatus(
	ctx context.Context,
	tx *gorm.DB,
	commentId uuid.UUID,
	status bool,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.Comment{}).
		Where("id = ?", commentId).
		Update("status", status).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (r *rReport) updateUserCommentStatus(
	ctx context.Context,
	tx *gorm.DB,
	userId uuid.UUID,
	status bool,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.Comment{}).
		Where("user_id = ?", userId).
		Update("status", status).
		Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (r *rReport) updateReportAdminIdAndStatus(
	ctx context.Context,
	tx *gorm.DB,
	reportId, adminId uuid.UUID,
	status bool,
) error {
	if err := tx.WithContext(ctx).
		Model(&models.Report{}).
		Where("id = ?", reportId).
		Updates(map[string]interface{}{
			"admin_id": adminId,
			"status":   status,
		}).Error; err != nil {
		return response.NewServerFailedError(err.Error())
	}
	return nil
}

func (r *rReport) checkPostReportExists(
	ctx context.Context,
	userId, reportedPostId uuid.UUID,
) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.PostReport{}).
		Joins("JOIN reports ON reports.id = post_reports.report_id").
		Where("reports.user_id = ? AND post_reports.reported_post_id = ?", userId, reportedPostId).
		Count(&count).Error; err != nil {
	}
	return count > 0, nil
}

func (r *rReport) checkUserReportExists(
	ctx context.Context,
	userId, reportedUserId uuid.UUID,
) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.UserReport{}).
		Joins("JOIN reports ON reports.id = user_reports.report_id").
		Where("reports.user_id = ? AND user_reports.reported_user_id = ?", userId, reportedUserId).
		Count(&count).Error; err != nil {
	}
	return count > 0, nil
}

func (r *rReport) checkCommentReportExists(
	ctx context.Context,
	userId, reportedCommentId uuid.UUID,
) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.CommentReport{}).
		Joins("JOIN reports ON reports.id = comment_reports.report_id").
		Where("reports.user_id = ? AND comment_reports.reported_comment_id = ?", userId, reportedCommentId).
		Count(&count).Error; err != nil {
	}
	return count > 0, nil
}
