package implement

import (
	"context"
	"fmt"
	advertise_command "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	advertise_query "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	advertise_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/payment"
	"net/http"
)

type sAdvertise struct {
	advertiseRepo    advertise_repo.IAdvertiseRepository
	billRepo         advertise_repo.IBillRepository
	postRepo         advertise_repo.IPostRepository
	newFeedRepo      advertise_repo.INewFeedRepository
	notificationRepo advertise_repo.INotificationRepository
}

func NewAdvertiseImplement(
	advertiseRepo advertise_repo.IAdvertiseRepository,
	billRepo advertise_repo.IBillRepository,
	postRepo advertise_repo.IPostRepository,
	newFeedRepo advertise_repo.INewFeedRepository,
	notificationRepo advertise_repo.INotificationRepository,
) *sAdvertise {
	return &sAdvertise{
		advertiseRepo:    advertiseRepo,
		billRepo:         billRepo,
		postRepo:         postRepo,
		newFeedRepo:      newFeedRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sAdvertise) CreateAdvertise(
	ctx context.Context,
	command *advertise_command.CreateAdvertiseCommand,
) (result *advertise_command.CreateAdvertiseResult, err error) {
	result = &advertise_command.CreateAdvertiseResult{}
	// 1. Create advertise
	advertiseEntity, err := advertise_entity.NewAdvertise(
		command.PostId,
		command.StartDate,
		command.EndDate,
	)
	if err != nil {
		result.PayUrl = ""
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	newAdvertise, err := s.advertiseRepo.CreateOne(ctx, advertiseEntity)
	if err != nil {
		result.PayUrl = ""
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 2. Create bill
	duration := command.EndDate.Sub(command.StartDate)
	durationDate := int(duration.Seconds() / 86400)
	price := durationDate*30000 + int(float64(durationDate) * 30000 * 0.1)

	fmt.Println(price)

	billEntity, err := advertise_entity.NewBill(
		newAdvertise.ID,
		price,
	)
	if err != nil {
		result.PayUrl = ""
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	newBill, err := s.billRepo.CreateOne(ctx, billEntity)
	if err != nil {
		result.PayUrl = ""
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 3. Call momo api to handle payment
	payUrl, err := payment.SendRequestToMomo(
		newBill.ID.String(),
		fmt.Sprintf("Pay for bill %s", newBill.ID),
		price,
		"Payment by momo",
		"This is momo payment",
		command.RedirectUrl,
	)

	if err != nil {
		result.PayUrl = ""
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	result.PayUrl = payUrl
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sAdvertise) GetAdvertise(
	ctx context.Context,
	query *advertise_query.GetOneAdvertiseQuery,
) (result *advertise_query.GetOneAdvertiseResult, err error) {
	result = &advertise_query.GetOneAdvertiseResult{}
	return
}
