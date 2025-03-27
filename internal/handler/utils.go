package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) newErrorResponse(c *gin.Context, statusCode int, message string) {
	h.logger.Error(message)
	
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

func (h *Handler) newSuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	h.logger.Info(data)

	c.JSON(statusCode, gin.H{
		"data": data,
	})
}

func (h *Handler) initPagination(c *gin.Context) (*pagination, error) {
	limit, err := strconv.Atoi(c.GetHeader("limit"))
	if err != nil {
		if c.GetHeader("limit") != "" {
			h.newErrorResponse(c, http.StatusBadRequest, "Invalid limit parameter")
			return nil, err
		}
	}

	offset, err := strconv.Atoi(c.GetHeader("offset"))
	if err != nil {
		if c.GetHeader("offset") != "" {
			h.newErrorResponse(c, http.StatusBadRequest, "Invalid offset parameter")
			return nil, err
		}
	}

	return &pagination{limit: limit, offset: offset}, nil
}