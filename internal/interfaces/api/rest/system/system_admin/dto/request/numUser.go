package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
)

type NumUsers struct {
	NumUsers int `json:"num_users"`
}

func ValidateNumUsers(req interface{}) error {
	dto, ok := req.(*NumUsers)
	if !ok {
		return fmt.Errorf("type assertion failed req.(*numUsers)")
	}

	return validation.ValidateStruct(dto)
}
