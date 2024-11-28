package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	"time"
)

type CreateAdvertiseRequest struct {
	PostId      uuid.UUID `json:"post_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
	RedirectUrl string    `json:"redirect_url" binding:"required"`
}

func ValidateCreateAdvertiseRequest(req interface{}) error {
	dto, ok := req.(*CreateAdvertiseRequest)
	if !ok {
		return fmt.Errorf("validate CreateAdvertiseRequest failed")
	}

	today := time.Now().Truncate(24 * time.Hour)

	return validation.ValidateStruct(dto,
		validation.Field(&dto.PostId, validation.Required),
		validation.Field(&dto.StartDate, validation.Required, validation.By(func(value interface{}) error {
			startDate, ok := value.(time.Time)
			if !ok {
				return fmt.Errorf("invalid start date format")
			}
			if startDate.Before(today) {
				return fmt.Errorf("start date must be today or later")
			}
			return nil
		})),
		validation.Field(&dto.EndDate, validation.Required, validation.By(func(value interface{}) error {
			endDate, ok := value.(time.Time)
			if !ok {
				return fmt.Errorf("invalid end date format")
			}
			if endDate.Before(dto.StartDate) {
				return fmt.Errorf("end date must be before start date")
			}

			startDate := dto.StartDate.Truncate(24 * time.Hour)
			endDate = endDate.Truncate(24 * time.Hour)
			if endDate.Sub(startDate).Hours() > 30*24 {
				return fmt.Errorf("start date and end date must not be more than 30 days apart")
			}
			return nil
		})),
		validation.Field(&dto.RedirectUrl, validation.Required, is.URL),
	)
}

func (req *CreateAdvertiseRequest) ToCreateAdvertiseCommand() (*command.CreateAdvertiseCommand, error) {
	return &command.CreateAdvertiseCommand{
		PostId:      req.PostId,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		RedirectUrl: req.RedirectUrl,
	}, nil
}
