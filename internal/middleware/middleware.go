package middleware

import (
	"net/http"
	"practice/internal/context"
	"practice/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AttachRequestID() gin.HandlerFunc {
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

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := ctx.GetHeader("Authorization")
		if h == "" {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		user := context.CheckUserExists(h)
		if user == nil {
			response.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		context.SetUser(ctx, *user)
		ctx.Next()
	}
}
