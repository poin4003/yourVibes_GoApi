package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
	"time"
)

type CreateAdvertiseRequest struct {
	PostId      uuid.UUID `json:"post_id" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required" validate:"gtfield=StartDate"`
	RedirectUrl string    `json:"redirect_url" binding:"required" validate:"url"`
}

func (req *CreateAdvertiseRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}

func (req *CreateAdvertiseRequest) ToCreateAdvertiseCommand() (*command.CreateAdvertiseCommand, error) {
	return &command.CreateAdvertiseCommand{
		PostId:      req.PostId,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		RedirectUrl: req.RedirectUrl,
	}, nil
}
