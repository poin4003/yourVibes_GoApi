package initialize

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/poin4003/yourVibes_GoApi/internal/utils/privacy"
)

func InitCustomValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("privacy_enum", privacy.ValidatePrivacy)
	}
}
