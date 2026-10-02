package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type requestIDContext struct{}

var requestIDKey *requestIDContext

func SetRequestID(ctx *gin.Context, incRequestID uuid.UUID) {
	if incRequestID != uuid.Nil {
		ctx.Set(&requestIDKey, incRequestID)
	} else {
		ctx.Set(&requestIDKey, generateRequestID())
	}
}

func GetRequestID(ctx *gin.Context) uuid.UUID {
	id, exists := ctx.Get(&requestIDKey)
	if !exists {
		panic("request id not set")
	}
	return id.(uuid.UUID)
}

func generateRequestID() uuid.UUID {
	return uuid.New()
}
