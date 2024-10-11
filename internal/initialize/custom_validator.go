package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
	"mime/multipart"
)

func InitCustomValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("privacy_enum", validatePrivacy)
		v.RegisterValidation("file", validateFile)
	}
}

func validatePrivacy(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == string(consts.PUBLIC) || value == string(consts.PRIVATE) || value == string(consts.FRIEND_ONLY) {
		return true
	}

	return false
}

func validateFile(fl validator.FieldLevel) bool {
	files := fl.Field().Interface().([]multipart.FileHeader)

	for _, file := range files {
		if file.Size == 0 {
			return false
		}
	}

	return true
}
