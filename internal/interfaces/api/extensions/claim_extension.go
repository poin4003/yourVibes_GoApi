package extensions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(ctx *gin.Context) (uuid.UUID, error) {
	userId, exists := ctx.Get("userId")
	if !exists {
		return uuid.Nil, fmt.Errorf("Unauthorized: user ID not found in context")
	}

	userUUID, ok := userId.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("Invalid user id format")
	}

	return userUUID, nil
}

func GetAdminID(ctx *gin.Context) (uuid.UUID, error) {
	adminId, exists := ctx.Get("adminId")
	if !exists {
		return uuid.Nil, fmt.Errorf("Unauthorized: admin ID not found in context")
	}

	userUUID, ok := adminId.(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("Invalid admin id format")
	}

	return userUUID, nil
}
