package context

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type requestIDContext struct{}
type requestContext struct{}

var requestIDKey *requestIDContext
var userContextKey *requestContext

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

func SetUser(ctx *gin.Context, user User) {
	ctx.Set(&userContextKey, user)
}

func GetUser(ctx *gin.Context) *User {
	u, exists := ctx.Get(&userContextKey)
	if !exists {
		return nil
	}
	user := u.(User)
	return &user
}

func generateRequestID() uuid.UUID {
	return uuid.New()
}
