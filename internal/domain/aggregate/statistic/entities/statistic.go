package entities

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"time"
)

type StatisticEntity struct {
	ID         uuid.UUID
	PostId     uuid.UUID
	Reach      int
	Clicks     int
	Impression int
	Status     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (vs *StatisticEntity) ValidateStatisticEntity() error {
	return validation.ValidateStruct(vs,
		validation.Field(&vs.PostId, validation.Required),
		validation.Field(&vs.Reach, validation.Min(0)),
		validation.Field(&vs.Clicks, validation.Min(0)),
		validation.Field(&vs.Impression, validation.Min(0)),
	)
}

func NewStatisticEntity(
	postId uuid.UUID,
	reach, clicks, impression int,
) (*StatisticEntity, error) {
	statistics := &StatisticEntity{
		ID:         uuid.New(),
		PostId:     postId,
		Reach:      reach,
		Clicks:     clicks,
		Impression: impression,
		Status:     true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := statistics.ValidateStatisticEntity(); err != nil {
		return nil, err
	}

	return statistics, nil
}
