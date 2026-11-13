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

func initUserRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	r := router.Group("/users")
	r.Use(middleware.AuthMiddleware(rdb))

	ur := repo.NewUserRepo(db)
	us := service.NewUserService(ur, rdb)
	uc := controller.NewUserController(us)

	r.GET("profile", uc.GetUserProfile)
	r.GET("info", uc.GetUserInfo)
}
