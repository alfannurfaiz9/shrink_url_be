package router

import (
	_ "github.com/alfannurfaiz9/shrink_url_be/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitMainRouter(r *gin.Engine) {
	r.GET("/documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	initAuthRouter(r)
}
