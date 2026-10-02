package middleware

import (
	"practice/internal/context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AttachRequestIDMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := ctx.GetHeader("X-Request-ID")

		if reqID == "" {
			context.SetRequestID(ctx, uuid.Nil)
		} else {
			parsedRequestID, err := uuid.Parse(reqID)
			if err != nil {
				panic("failed to parse uuid")
			}
			context.SetRequestID(ctx, parsedRequestID)
		}
		ctx.Next()
	}
}
