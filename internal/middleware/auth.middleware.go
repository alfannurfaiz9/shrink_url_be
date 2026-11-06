package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	"github.com/alfannurfaiz9/shrink_url_be/internal/utils"
	"github.com/alfannurfaiz9/shrink_url_be/pkg"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearer := ctx.GetHeader("Authorization")
		if bearer == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Code:    "UNAUTHORIZED",
				Message: "please login first",
			})
			return
		}

		result := strings.Split(bearer, " ")
		if len(result) != 2 || result[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Code:    "UNAUTHORIZED",
				Message: "invalid bearer token",
			})
			return
		}

		token := &pkg.JWTClaims{}
		err := token.DecodeToken(result[1])

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenInvalidIssuer) {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Success: false,
					Code:    "UNAUTHORIZED",
					Message: "invalid token",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{
				Success: false,
				Code:    "INTERNAL_SERVER_ERROR",
				Message: "internal server error",
			})
			return
		}

		key := fmt.Sprintf("shrinkurl:token:%s", token.ID)
		_, err = utils.GetRedisStr(ctx, rdb, key)
		fmt.Println(key)
		if err == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Code:    "UNAUTHORIZED",
				Message: "please login first",
			})
			return
		}

		ctx.Set("id", token.Id)
		ctx.Set("ID", token.ID)
		ctx.Set("exp", token.ExpiresAt.Time)
		ctx.Next()
	}
}
