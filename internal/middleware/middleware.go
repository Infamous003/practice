package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestIDKey string

const (
	RequestID RequestIDKey = "request_id"
)

func AttachRequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var id uuid.UUID
		reqID := ctx.GetHeader("X-Request-ID")

		if reqID == "" {
			id = uuid.New()
		} else {
			parsedRequestID, err := uuid.Parse(reqID)
			if err != nil {
				panic("failed to parse uuid")
			}
			id = parsedRequestID
		}

		ctx.Set(RequestID, id)
		ctx.Next()
	}
}
