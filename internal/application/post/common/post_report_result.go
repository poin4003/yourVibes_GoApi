package common

import (
	"time"

	"github.com/google/uuid"
)

type PostReportResult struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
	AdminId        *uuid.UUID
	User           *UserForReportResult
	ReportedPost   *PostForReportResult
	Admin          *AdminResult
	Reason         string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PostReportShortVerResult struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
	AdminId        *uuid.UUID
	Reason         string
	UserEmail      string
	AdminEmail     *string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
