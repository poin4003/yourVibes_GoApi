package implement

import (
	"context"

	billCommand "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	billEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	postEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/post/entities"
	billRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/pointer"
)

type sBill struct {
	advertiseRepo    billRepo.IAdvertiseRepository
	billRepo         billRepo.IBillRepository
	postRepo         billRepo.IPostRepository
	notificationRepo billRepo.INotificationRepository
}

func NewBillImplement(
	advertiseRepo billRepo.IAdvertiseRepository,
	billRepo billRepo.IBillRepository,
	postRepo billRepo.IPostRepository,
	notificationRepo billRepo.INotificationRepository,
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
	command *billCommand.ConfirmPaymentCommand,
) error {
	if command == nil {
		return response.NewServerFailedError("command confirm payment not found")
	}

	// 1. Find bill
	billFound, err := s.billRepo.GetById(ctx, *command.BillId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if billFound == nil {
		return response.NewDataNotFoundError("bill not found")
	}

	// 2. Update status bill to paid
	updateBillData := &billEntity.BillUpdate{
		Status: pointer.Ptr(true),
	}

	err = updateBillData.ValidateUpdateBill()
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.billRepo.UpdateOne(ctx, billFound.ID, updateBillData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	// 3. Update post to isAdvertisement
	// 3.1. Find post
	postFound, err := s.postRepo.GetById(ctx, billFound.Advertise.PostId)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	if postFound == nil {
		return response.NewDataNotFoundError("post not found")
	}

	// 3.2. Update isAdvertisement
	updatePostData := &postEntity.PostUpdate{
		IsAdvertisement: pointer.Ptr(true),
	}

	err = updatePostData.ValidatePostUpdate()
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	_, err = s.postRepo.UpdateOne(ctx, postFound.ID, updatePostData)
	if err != nil {
		return response.NewServerFailedError(err.Error())
	}

	return nil
}
