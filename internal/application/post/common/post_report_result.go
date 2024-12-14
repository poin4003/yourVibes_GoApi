package common

import (
	"github.com/google/uuid"
	"time"
)

type PostReportResult struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
	AdminId        uuid.UUID
	User           *UserResult
	ReportedPost   *PostResult
	Admin          *AdminResult
	Reason         string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PostReportShortVerResult struct {
	UserId         uuid.UUID
	ReportedPostId uuid.UUID
	AdminId        uuid.UUID
	Reason         string
	Status         bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
