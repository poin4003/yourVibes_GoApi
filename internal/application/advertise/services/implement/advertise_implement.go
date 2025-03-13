package implement

import (
	"context"
	response2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/payment"
	"time"

	advertiseCommand "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/mapper"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	advertiseEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	advertiseRepo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
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
	// 1. Check previous ad status
	// 1.1. Check if the post has had any ads before by bill
	billStatus, err := s.billRepo.CheckExists(ctx, command.PostId)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	// 1.2. If bill has exists
	if billStatus {
		// 1.2.1. Get latest ad
		latestAds, err := s.advertiseRepo.GetLatestAdsByPostId(ctx, command.PostId)
		if err != nil {
			return nil, response2.NewServerFailedError(err.Error())
		}

		// 1.2.2. Check payment status
		if latestAds.Bill.Status {
			// 1.2.2.1. Check ads expiration date
			today := time.Now()
			if !today.After(latestAds.EndDate) {
				return nil, response2.NewCustomError(response2.ErrAdsExpired)
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
		return nil, response2.NewServerFailedError(err.Error())
	}

	advertiseCreated, err := s.advertiseRepo.CreateOne(ctx, newAdvertise)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
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
		return nil, response2.NewServerFailedError(err.Error())
	}

	newBill, err := s.billRepo.CreateOne(ctx, billEntity)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
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
		return nil, response2.NewServerFailedError(err.Error())
	}

	return &advertiseCommand.CreateAdvertiseResult{
		PayUrl: payUrl,
	}, nil
}

func (s *sAdvertise) GetManyAdvertise(
	ctx context.Context,
	query *advertiseQuery.GetManyAdvertiseQuery,
) (result *advertiseQuery.GetManyAdvertiseResults, err error) {
	advertiseEntities, paging, err := s.advertiseRepo.GetMany(ctx, query)
	if err != nil {
		return result, err
	}

	var advertiseResults []*common.AdvertiseWithBillResult
	for _, advertise := range advertiseEntities {
		advertiseResults = append(advertiseResults, mapper.NewAdvertiseWithBillResultFromEntity(advertise))
	}

	return &advertiseQuery.GetManyAdvertiseResults{
		Advertises:     advertiseResults,
		PagingResponse: paging,
	}, nil
}

func (s *sAdvertise) GetAdvertise(
	ctx context.Context,
	query *advertiseQuery.GetOneAdvertiseQuery,
) (result *advertiseQuery.GetOneAdvertiseResult, err error) {
	// 1. Get advertise detail
	advertise, err := s.advertiseRepo.GetOne(ctx, query.AdvertiseId)
	if err != nil {
		return nil, response2.NewServerFailedError(err.Error())
	}

	if advertise == nil {
		return nil, response2.NewDataNotFoundError("advertise not found")
	}

	return &advertiseQuery.GetOneAdvertiseResult{
		Advertise: mapper.NewAdvertiseDetailFromEntity(advertise),
	}, nil
}
