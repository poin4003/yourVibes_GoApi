package extensions

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(ctx *gin.Context) (uuid.UUID, error) {
	userId, exists := ctx.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("unauthorized: user ID not found in context")
	}

	userUUID, ok := userId.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid user id format")
	}

	return userUUID, nil
}

func GetAdminID(ctx *gin.Context) (uuid.UUID, error) {
	adminId, exists := ctx.Get("adminId")
	if !exists {
		return uuid.Nil, fmt.Errorf("unauthorized: admin ID not found in context")
	}

	userUUID, ok := adminId.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid admin id format")
	}

	return userUUID, nil
}
