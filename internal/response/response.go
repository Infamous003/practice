package response

import (
	"net/http"

	"practice/internal/context"

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
	requestID := context.GetRequestID(c)

	c.JSON(http.StatusOK, response{
		Success:   true,
		Data:      data,
		RequestID: requestID,
	})
}

func Created(c *gin.Context, data any) {
	requestID := context.GetRequestID(c)

	c.JSON(http.StatusCreated, response{
		Success:   true,
		Data:      data,
		RequestID: requestID,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	requestID := context.GetRequestID(c)
	c.JSON(code, response{
		Success:   false,
		Error:     message,
		RequestID: requestID,
	})
}
