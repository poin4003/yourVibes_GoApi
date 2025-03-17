package common

import "github.com/google/uuid"

type UserShortVerResult struct {
	ID         uuid.UUID
	FamilyName string
	Name       string
	AvatarUrl  string
}
