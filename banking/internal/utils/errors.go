package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
	Error   string `json:"error,omitempty" example:"Invalid input"`
}

func RespondWithError(c *gin.Context, code int, message string, err interface{}) {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}
	if err != nil {
		var errStr string
		switch v := err.(type) {
		case error:
			errStr = v.Error()
		case string:
			errStr = v
		default:
			errStr = fmt.Sprintf("%v", err)
		}
		response.Error = errStr
	}
	c.JSON(code, response)
}


func RespondWithBadRequest(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusBadRequest, message, err)
}

func RespondWithNotFound(c *gin.Context, message string) {
	RespondWithError(c, http.StatusNotFound, message, nil)
}

func RespondWithInternalError(c *gin.Context, message string, err interface{}) {
	RespondWithError(c, http.StatusInternalServerError, message, err)
}

func RespondWithUnauthorized(c *gin.Context, message string) {
	RespondWithError(c, http.StatusUnauthorized, message, nil)
}
