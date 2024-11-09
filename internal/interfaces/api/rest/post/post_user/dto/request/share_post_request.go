package request

import (
	"github.com/poin4003/yourVibes_GoApi/internal/consts"
)

type SharePostInput struct {
	Content  string              `form:"content,omitempty" binding:"omitempty"`
	Privacy  consts.PrivacyLevel `form:"privacy,omitempty" binding:"omitempty,privacy_enum"`
	Location string              `form:"location,omitempty"`
}
