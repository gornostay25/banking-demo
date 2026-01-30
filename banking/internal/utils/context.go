package utils

import (
	"errors"

	"banking/internal/models"

	"github.com/gin-gonic/gin"
)

// IdentityKey is the key used to store user in gin context
const IdentityKey = "user_id"

var (
	ErrUserNotInContext = errors.New("user not found in context")
	ErrInvalidUserType  = errors.New("invalid user type in context")
)

// GetUserFromContext extracts the authenticated user from gin context.
// This function assumes the user is already authenticated by middleware.
// It should only be called in handlers protected by authentication middleware.
func GetUserFromContext(c *gin.Context) (*models.User, error) {
	user, exists := c.Get(IdentityKey)
	if !exists {
		return nil, ErrUserNotInContext
	}

	userModel, ok := user.(*models.User)
	if !ok {
		return nil, ErrInvalidUserType
	}

	return userModel, nil
}

// MustGetUserFromContext extracts the authenticated user from gin context.
// Panics if user is not found or has invalid type.
// Use this when you're certain the middleware has set the user.
func MustGetUserFromContext(c *gin.Context) *models.User {
	user, err := GetUserFromContext(c)
	if err != nil {
		panic(err)
	}
	return user
}
