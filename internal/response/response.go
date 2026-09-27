package response

import (
	"net/http"

	"practice/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type response struct {
	Success   bool      `json:"success"`
	Data      any       `json:"data,omitempty"`
	Error     string    `json:"error,omitempty"`
	RequestID uuid.UUID `json:"request_id"`
}

func OK(c *gin.Context, data any) {
	id, exists := c.Get(middleware.RequestID)
	if !exists {
		panic("request id not set")
	}

	requestUUID, ok := id.(uuid.UUID)
	if !ok {
		panic("type assertion failed: request id must be a valid uuid")
	}

	c.JSON(http.StatusOK, response{
		Success:   true,
		Data:      data,
		RequestID: requestUUID,
	})
}

func Created(c *gin.Context, data any) {
	id, exists := c.Get(middleware.RequestID)
	if !exists {
		panic("request id not set")
	}

	requestUUID, ok := id.(uuid.UUID)
	if !ok {
		panic("type assertion failed: request id must be a valid uuid")
	}

	c.JSON(http.StatusCreated, response{
		Success:   true,
		Data:      data,
		RequestID: requestUUID,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	id, exists := c.Get(middleware.RequestID)
	if !exists {
		panic("request id not set")
	}

	requestUUID, ok := id.(uuid.UUID)
	if !ok {
		panic("type assertion failed: request id must be a valid uuid")
	}
	c.JSON(code, response{
		Success:   false,
		Error:     message,
		RequestID: requestUUID,
	})
}
