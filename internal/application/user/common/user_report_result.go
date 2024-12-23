package common

import (
	"time"

	"github.com/google/uuid"
)

type UserReportResult struct {
	UserId         uuid.UUID
	ReportedUserId uuid.UUID
	AdminId        *uuid.UUID
	User           *UserForReportResult
	ReportedUser   *UserForReportResult
	Admin          *AdminResult
	Reason         string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserReportShortVerResult struct {
	UserId            uuid.UUID
	ReportedUserId    uuid.UUID
	AdminId           *uuid.UUID
	Reason            string
	UserEmail         string
	ReportedUserEmail string
	AdminEmail        *string
	Status            bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
