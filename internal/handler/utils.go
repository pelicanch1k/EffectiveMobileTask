package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/pelicanch1k/EffectiveMobileTestTask/pkg/logging"
	"net/http"
)

type errorResponse struct {
	Message string `json:"message"`
}

type statusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logging.GetLogger().Error(message)
	if errorHandler(c, message) {
		c.AbortWithStatusJSON(statusCode, errorResponse{message})
	}
}

func errorHandler(c *gin.Context, message string) bool {
	if message == "pq: duplicate key value violates unique constraint \"users_username_key\"" {
		c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{"username is already taken"})
		return false
	}

	return true
}
