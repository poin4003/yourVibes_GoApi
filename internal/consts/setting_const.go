package consts

type Language string

const (
	VI Language = "vi"
	EN Language = "en"
)

var Languages = []interface{}{
	VI,
	EN,
}
