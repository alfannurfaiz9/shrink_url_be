package router

import (
	_ "github.com/alfannurfaiz9/shrink_url_be/docs"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitMainRouter(r *gin.Engine, db *pgxpool.Pool) {
	r.GET("/documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	initAuthRouter(r, db)
}
