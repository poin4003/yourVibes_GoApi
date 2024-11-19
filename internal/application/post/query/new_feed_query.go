package query

import "github.com/google/uuid"

type GetNewFeedQuery struct {
	UserId uuid.UUID
	Limit  int
	Page   int
}
