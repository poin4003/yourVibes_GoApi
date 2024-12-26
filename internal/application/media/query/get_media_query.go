package query

type MediaQuery struct {
	FileName string
}

type MediaQueryResult struct {
	FilePath       string
	ResultCode     int
	HttpStatusCode int
}
