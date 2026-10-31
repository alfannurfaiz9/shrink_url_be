package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func initAuthRouter(r *gin.Engine) {
	r.Group("/auth")

	r.GET("/", func(ctx *gin.Context) {
		fmt.Println("ok")
	})
}
