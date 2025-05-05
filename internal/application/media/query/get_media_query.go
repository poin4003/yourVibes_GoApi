package query

import (
	"io"
	"os"
	"time"
)

type MediaQuery struct {
	FileName    string
	RangeHeader string
}

type MediaQueryResult struct {
	File       io.ReadSeeker
	RawFile    *os.File
	ModTime    time.Time
	Headers    map[string]string
	StatusCode int
}
