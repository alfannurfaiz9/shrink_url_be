package router

import (
	"github.com/alfannurfaiz9/shrink_url_be/internal/controller"
	"github.com/alfannurfaiz9/shrink_url_be/internal/middleware"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/alfannurfaiz9/shrink_url_be/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func initLinkRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/links")
	r.Use(middleware.AuthMiddleware(rdb))

	lr := repo.NewLinksRepo(db)
	ls := service.NewLinkService(lr, rdb)
	lc := controller.NewLinkController(ls)

	r.POST("", lc.CreateLink)
	r.GET("", lc.GetUserLinks)
	r.DELETE("/:id", lc.DeleteSlug)
	r.GET("/check", lc.CheckExistingLinks)

	router.GET("/:slug", lc.GetOriginalUrl)
}
