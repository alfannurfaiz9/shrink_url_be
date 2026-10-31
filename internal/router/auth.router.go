package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initAuthRouter(r *gin.Engine, db *pgxpool.Pool) {
	r.Group("/auth")

	r.GET("/", func(ctx *gin.Context) {
		fmt.Println("ok")
	})
}
