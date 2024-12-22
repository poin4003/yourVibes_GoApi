package implement

import (
	"context"
	"errors"
	"fmt"
	advertise_command "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/mapper"
	advertise_query "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	advertise_entity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	"github.com/poin4003/yourVibes_GoApi/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/payment"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type sAdvertise struct {
	advertiseRepo    advertise_repo.IAdvertiseRepository
	billRepo         advertise_repo.IBillRepository
	notificationRepo advertise_repo.INotificationRepository
}

func NewAdvertiseImplement(
	advertiseRepo advertise_repo.IAdvertiseRepository,
	billRepo advertise_repo.IBillRepository,
	notificationRepo advertise_repo.INotificationRepository,
) *sAdvertise {
	return &sAdvertise{
		advertiseRepo:    advertiseRepo,
		billRepo:         billRepo,
		notificationRepo: notificationRepo,
	}
}

func (s *sAdvertise) CreateAdvertise(
	ctx context.Context,
	command *advertise_command.CreateAdvertiseCommand,
) (result *advertise_command.CreateAdvertiseResult, err error) {
	result = &advertise_command.CreateAdvertiseResult{}
	// 1. Check previous ad status
	// 1.1. Check if the post has had any ads before by bill
	billStatus, err := s.billRepo.CheckExists(ctx, command.PostId)
	if err != nil {
		result.PayUrl = ""
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		return result, err
	}

	// 1.2. If bill has exists
	if billStatus {
		// 1.2.1. Get latest ad
		latestAds, err := s.advertiseRepo.GetLatestAdsByPostId(ctx, command.PostId)
		if err != nil {
			result.PayUrl = ""
			result.ResultCode = response.ErrServerFailed
			result.HttpStatusCode = http.StatusInternalServerError
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

	// 3. Create bill
	duration := command.EndDate.Sub(command.StartDate)
	durationDate := int(duration.Seconds() / 86400)
	price := durationDate*30000 + int(float64(durationDate)*30000*0.1)

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

func (s *sAdvertise) GetManyAdvertise(
	ctx context.Context,
	query *advertise_query.GetManyAdvertiseQuery,
) (result *advertise_query.GetManyAdvertiseResults, err error) {
	result = &advertise_query.GetManyAdvertiseResults{}

	advertiseEntities, paging, err := s.advertiseRepo.GetMany(ctx, query)
	if err != nil {
		result.Advertises = nil
		result.ResultCode = response.ErrServerFailed
		result.HttpStatusCode = http.StatusInternalServerError
		result.PagingResponse = nil
		return result, err
	}

	var advertiseResults []*common.AdvertiseWithBillResult
	for _, advertiseEntity := range advertiseEntities {
		advertiseResults = append(advertiseResults, mapper.NewAdvertiseWithBillResultFromEntity(advertiseEntity))
	}

	result.Advertises = advertiseResults
	result.ResultCode = response.ErrCodeSuccess
	result.HttpStatusCode = http.StatusOK
	result.PagingResponse = paging
	return result, nil
}

func (s *sAdvertise) GetAdvertise(
	ctx context.Context,
	query *advertise_query.GetOneAdvertiseQuery,
) (result *advertise_query.GetOneAdvertiseResult, err error) {
	result = &advertise_query.GetOneAdvertiseResult{}
	result.Advertise = nil
	result.ResultCode = response.ErrServerFailed
	result.HttpStatusCode = http.StatusInternalServerError
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
