package consts

type ReportType int16

const (
	USER_REPORT    ReportType = 0
	POST_REPORT    ReportType = 1
	COMMENT_REPORT ReportType = 2
)

var ReportTypes = []interface{}{
	USER_REPORT,
	POST_REPORT,
	COMMENT_REPORT,
}
