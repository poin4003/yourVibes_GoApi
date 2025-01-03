package entities

import "time"

type Revenue struct {
	Month time.Time
	Total int64
}
