package query

import (
	"os"
	"time"
)

type MediaQuery struct {
	FileName    string
	RangeHeader string
}

type MediaQueryResult struct {
	File       *os.File
	ModTime    time.Time
	Headers    map[string]string
	StatusCode int
}
