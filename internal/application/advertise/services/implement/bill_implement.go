package implement

import (
	"context"
	"errors"
	bill_command "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	post_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	bill_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	bill_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
	"gorm.io/gorm"
	"net/http"
)

type sBill struct {
	advertiseRepo    bill_repo.IAdvertiseRepository
	billRepo         bill_repo.IBillRepository
	postRepo         bill_repo.IPostRepository
	notificationRepo bill_repo.INotificationRepository
}

func NewBillImplement(
	advertiseRepo bill_repo.IAdvertiseRepository,
	billRepo bill_repo.IBillRepository,
	postRepo bill_repo.IPostRepository,
	notificationRepo bill_repo.INotificationRepository,
) *sBill {
	return &sBill{
		advertiseRepo:    advertiseRepo,
		billRepo:         billRepo,
		postRepo:         postRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sBill) ConfirmPayment(
	ctx context.Context,
	command *bill_command.ConfirmPaymentCommand,
) (result *bill_command.ConfirmPaymentResult, err error) {
	result = &bill_command.ConfirmPaymentResult{}
	// 1. Find bill
	billFound, err := s.billRepo.GetById(ctx, *command.BillId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 2. Update status bill to paid
	updateBillData := &bill_entity.BillUpdate{
		Status: pointer.Ptr(true),
	}

	err = updateBillData.ValidateUpdateBill()
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	_, err = s.billRepo.UpdateOne(ctx, billFound.ID, updateBillData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 3. Update post to isAdvertisement
	// 3.1. Find post
	postFound, err := s.postRepo.GetById(ctx, billFound.Advertise.PostId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 3.2. Update isAdvertisement
	updatePostData := &post_entity.PostUpdate{
		IsAdvertisement: pointer.Ptr(true),
	}

	err = updatePostData.ValidatePostUpdate()
	if err != nil {
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePostData)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusBadRequest
			return result, err
		}
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
