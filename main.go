package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Cors(ctx *gin.Context) {
	origin := os.Getenv("ALLOWED_ORIGIN")

	ctx.Header("Access-Control-Allow-Origin", origin)
	ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	ctx.Header("Access-Control-Allow-Methods", "OPTIONS, POST, DELETE, PATCH, PUT, GET")
	ctx.Header("Access-Control-Allow-Credentials", "true")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}
	ctx.Next()
}
