package command

import (
	"github.com/google/uuid"
)

type UpsertStatisticCommand struct {
	PostId     uuid.UUID
	Reach      int
	Clicks     int
	Impression int
}
