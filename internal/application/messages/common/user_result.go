package common

import (
	"github.com/google/uuid"
)

type UserResult struct {
	ID         uuid.UUID
	FamilyName string
	Name       string
	AvatarUrl  string
}
