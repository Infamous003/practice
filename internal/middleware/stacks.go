package middleware

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type User struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// 	Role     string `json:"role"`
// }

// var users = []User{
// 	{ID: 1, Username: "infamous", Role: "admin"},
// 	{ID: 2, Username: "syed", Role: "user"},
// }

// func checkUserExists(username string) *User {
// 	var user *User
// 	for _, u := range users {
// 		if u.Username == username {
// 			user = &u
// 			return user
// 		}
// 	}
// 	return nil
// }

// func AuthenticateMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		h := ctx.GetHeader("Authorization")
// 		if h == "" {
// 			// response.ErrorResponse(ctx, http.StatusUnauthorized, "")
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
// 			return
// 		}

// 		user := checkUserExists(h)
// 		if user == nil {
// 			// response.ErrorResponse(ctx, http.StatusUnauthorized, "")
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set(UserContextKey, user)
// 		ctx.Next()
// 	}
// }
