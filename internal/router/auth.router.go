package router

import (
	"github.com/alfannurfaiz9/shrink_url_be/internal/controller"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/alfannurfaiz9/shrink_url_be/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func initAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	r := router.Group("/auth")

	ar := repo.NewAuthRepo(db)
	as := service.NewAuthService(ar)
	ac := controller.NewAuthController(as)

	r.POST("register", ac.Register)
}
