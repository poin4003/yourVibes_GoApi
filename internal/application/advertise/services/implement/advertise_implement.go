package implement

import (
	"context"
	"github.com/poin4003/yourVibes_GoApi/internal/domain/cache"
	"time"

	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	voucherEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/voucher/entities"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/payment"

	advertiseCommand "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/common"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/mapper"
	advertiseQuery "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/query"
	advertiseEntity "github.com/poin4003/yourVibes_GoApi/internal/domain/aggregate/advertise/entities"
	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
)

type sAdvertise struct {
	advertiseRepo repository.IAdvertiseRepository
	billRepo      repository.IBillRepository
	voucherRepo   repository.IVoucherRepository
	postCache     cache.IPostCache
}

func NewAdvertiseImplement(
	advertiseRepo repository.IAdvertiseRepository,
	billRepo repository.IBillRepository,
	voucherRepo repository.IVoucherRepository,
	postCache cache.IPostCache,
) *sAdvertise {
	return &sAdvertise{
		advertiseRepo: advertiseRepo,
		billRepo:      billRepo,
		voucherRepo:   voucherRepo,
		postCache:     postCache,
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
		return nil, response.NewServerFailedError(err.Error())
	}

	// 1.2. If bill has exists
	if billStatus {
		// 1.2.1. Get latest ad
		latestAds, err := s.advertiseRepo.GetLatestAdsByPostId(ctx, command.PostId)
		if err != nil {
			return nil, response.NewServerFailedError(err.Error())
		}

		// 1.2.2. Check payment status
		if latestAds.Bill.Status {
			// 1.2.2.1. Check ads expiration date
			today := time.Now()
			if !today.After(latestAds.EndDate) {
				return nil, response.NewCustomError(response.ErrAdsExpired)
			}
		}
	}

	// 2. Calculate bill
	duration := command.EndDate.Sub(command.StartDate)
	durationDate := int(duration.Seconds() / 86400)
	price := durationDate*30000 + int(float64(durationDate)*30000*0.1)

	var voucher *voucherEntity.VoucherEntity
	if command.VoucherCode != nil {
		// Get and redeem voucher
		voucher, err = s.voucherRepo.RedeemVoucher(ctx, *command.VoucherCode)
		if err != nil {
			return nil, err
		}

		// Check and calculate discount
		discount := 0
		if voucher.Type == consts.PERCENTAGE {
			discount = price * voucher.Value / 100
		} else {
			discount = price - voucher.Value
		}

		// Apply discount to price
		price -= discount
	}

	// 3. Apply voucher into bill if it exists
	var voucherId *uuid.UUID
	if voucher != nil {
		voucherId = &voucher.ID
	}

	// 4. Create advertise
	newAdvertise, err := advertiseEntity.NewAdvertise(
		command.PostId,
		command.StartDate,
		command.EndDate,
	)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	advertiseCreated, err := s.advertiseRepo.CreateOne(ctx, newAdvertise)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	billEntity, err := advertiseEntity.NewBill(
		advertiseCreated.ID,
		price,
		voucherId,
	)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	newBill, err := s.billRepo.CreateOne(ctx, billEntity)
	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 5. Call momo api to handle payment
	payUrl, err := payment.SendRequestToMomo(
		newBill.ID.String(),
		newBill.ID.String(),
		price,
		"Payment by momo",
		"This is momo payment",
		command.RedirectUrl,
	)

	if err != nil {
		return nil, response.NewServerFailedError(err.Error())
	}

	// 6. Delete post cache
	s.postCache.DeletePost(ctx, command.PostId)

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
		return nil, response.NewServerFailedError(err.Error())
	}

	if advertise == nil {
		return nil, response.NewDataNotFoundError("advertise not found")
	}

	return &advertiseQuery.GetOneAdvertiseResult{
		Advertise: mapper.NewAdvertiseDetailFromEntity(advertise),
	}, nil
}

func (s *sAdvertise) GetAdvertiseWithStatistic(
	ctx context.Context,
	AdvertiseId uuid.UUID,
) (result *common.AdvertiseForStatisticResult, err error) {
	advertise, err := s.advertiseRepo.GetDetailAndStatisticOfAdvertise(ctx, AdvertiseId)
	if err != nil {
		return nil, err
	}

	return mapper.NewAdvertiseDetailAndStatisticResult(advertise), nil
}

func (s *sAdvertise) GetShortAdvertiseByUserId(
	ctx context.Context,
	query *advertiseQuery.GetManyAdvertiseByUserId,
) (result *advertiseQuery.GetManyAdvertiseResultsByUserId, err error) {
	advertiseEntities, paging, err := s.advertiseRepo.GetAdvertiseByUserId(ctx, query)
	if err != nil {
		return nil, err
	}

	var advertiseResults []*common.ShortAdvertiseResult
	for _, advertise := range advertiseEntities {
		advertiseResults = append(advertiseResults, mapper.NewShortAdvertiseResult(advertise))
	}

	return &advertiseQuery.GetManyAdvertiseResultsByUserId{
		Advertises:     advertiseResults,
		PagingResponse: paging,
	}, nil
}
