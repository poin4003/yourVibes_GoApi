package implement

import (
	"context"
	"errors"
	"net/http"
	"time"

	advertiseCommand "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/mapper"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	advertiseEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertiseRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/payment"
	"gorm.io/gorm"
)

type sAdvertise struct {
	advertiseRepo    advertiseRepo.IAdvertiseRepository
	billRepo         advertiseRepo.IBillRepository
	notificationRepo advertiseRepo.INotificationRepository
}

func NewAdvertiseImplement(
	advertiseRepo advertiseRepo.IAdvertiseRepository,
	billRepo advertiseRepo.IBillRepository,
	notificationRepo advertiseRepo.INotificationRepository,
) *sAdvertise {
	return &sAdvertise{
		advertiseRepo:    advertiseRepo,
		billRepo:         billRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sAdvertise) CreateAdvertise(
	ctx context.Context,
	command *advertiseCommand.CreateAdvertiseCommand,
) (result *advertiseCommand.CreateAdvertiseResult, err error) {
	result = &advertiseCommand.CreateAdvertiseResult{
		PayUrl:         "",
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Check previous ad status
	// 1.1. Check if the post has had any ads before by bill
	billStatus, err := s.billRepo.CheckExists(ctx, command.PostId)
	if err != nil {
		return result, err
	}

	// 1.2. If bill has exists
	if billStatus {
		// 1.2.1. Get latest ad
		latestAds, err := s.advertiseRepo.GetLatestAdsByPostId(ctx, command.PostId)
		if err != nil {
			return result, err
		}

		// 1.2.2. Check payment status
		if latestAds.Bill.Status == true {
			// 1.2.2.1. Check ads expiration date
			today := time.Now()
			if !today.After(latestAds.EndDate) {
				result.PayUrl = ""
				result.ResultCode = response.ErrAdsExpired
				result.HttpStatusCode = http.StatusBadRequest
				return result, err
			}
		}
	}

	// 2. Create advertise
	newAdvertise, err := advertiseEntity.NewAdvertise(
		command.PostId,
		command.StartDate,
		command.EndDate,
	)
	if err != nil {
		return result, err
	}

	advertiseCreated, err := s.advertiseRepo.CreateOne(ctx, newAdvertise)
	if err != nil {
		return result, err
	}

	// 3. Create bill
	duration := command.EndDate.Sub(command.StartDate)
	durationDate := int(duration.Seconds() / 86400)
	price := durationDate*30000 + int(float64(durationDate)*30000*0.1)

	billEntity, err := advertiseEntity.NewBill(
		advertiseCreated.ID,
		price,
	)
	if err != nil {
		return result, err
	}

	newBill, err := s.billRepo.CreateOne(ctx, billEntity)
	if err != nil {
		return result, err
	}

	// 4. Call momo api to handle payment
	payUrl, err := payment.SendRequestToMomo(
		newBill.ID.String(),
		newBill.ID.String(),
		price,
		"Payment by momo",
		"This is momo payment",
		command.RedirectUrl,
	)

	if err != nil {
		return result, err
	}

	result.PayUrl = payUrl
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}

func (s *sAdvertise) GetManyAdvertise(
	ctx context.Context,
	query *advertiseQuery.GetManyAdvertiseQuery,
) (result *advertiseQuery.GetManyAdvertiseResults, err error) {
	result = &advertiseQuery.GetManyAdvertiseResults{
		Advertises:     nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
		PagingResponse: nil,
	}
	advertiseEntities, paging, err := s.advertiseRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	var advertiseResults []*common.AdvertiseWithBillResult
	for _, advertise := range advertiseEntities {
		advertiseResults = append(advertiseResults, mapper.NewAdvertiseWithBillResultFromEntity(advertise))
	}

	result.Advertises = advertiseResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}

func (s *sAdvertise) GetAdvertise(
	ctx context.Context,
	query *advertiseQuery.GetOneAdvertiseQuery,
) (result *advertiseQuery.GetOneAdvertiseResult, err error) {
	result = &advertiseQuery.GetOneAdvertiseResult{
		Advertise:      nil,
		ResultCode:     response.ErrServerFailed,
		HttpStatusCode: http.StatusInternalServerError,
	}
	// 1. Get advertise detail
	advertise, err := s.advertiseRepo.GetOne(ctx, query.AdvertiseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.ResultCode = response.ErrDataNotFound
			result.HttpStatusCode = http.StatusNotFound
			return result, nil
		}
		return result, err
	}

	result.Advertise = mapper.NewAdvertiseDetailFromEntity(advertise)
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	return result, nil
}
