package object

import (
	"github.com/poin4003/yourVibes_GoApi/internal/application/message/common"
)

type MessageUserObject struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      UserObject `json:"user"`
}

type UserObject struct {
	ID         string `json:"id"`
	FamilyName string `json:"family_name"`
	Name       string `json:"name"`
	AvatarUrl  string `json:"avatar_url"`
}

func ToMessageUserObject(messageResult *common.MessageResult) *MessageUserObject {
	return &MessageUserObject{
		ID:        messageResult.ID.String(),
		Content:   messageResult.Content,
		CreatedAt: messageResult.CreatedAt.Format("2006-01-02 15:04:05"),
		User:      *ToUserObject(messageResult.User),
	}
}

func ToUserObject(userResult *common.UserResult) *UserObject {
	return &UserObject{
		ID:         userResult.ID.String(),
		FamilyName: userResult.FamilyName,
		Name:       userResult.Name,
		AvatarUrl:  userResult.AvatarUrl,
	}
}