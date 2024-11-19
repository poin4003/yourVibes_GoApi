package query

import "github.com/google/uuid"

type GetPostLikeQuery struct {
	PostId uuid.UUID
	Limit  int
	Page   int
}
