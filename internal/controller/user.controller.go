package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	us *service.UserService
}

func NewUserController(us *service.UserService) *UserController {
	return &UserController{
		us: us,
	}
}

// GetUserProfile
// @Summary      Get User Profile
// @Description  Retrieves the user profile information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerToken
// @Router       /users/profile [get]
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
func (u *UserController) GetUserProfile(ctx *gin.Context) {
	id, exist := ctx.Get("id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "token does not exist",
		})
		return
	}

	idInt, ok := id.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "invalid payload format",
		})
		return
	}

	result, err := u.us.GetUserProfile(ctx, idInt)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, app_error.UserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Code:    "USER_NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Result:  result,
		Message: "successfully get user profile",
	})
}

// GetUserInfo
// @Summary      Get User Info
// @Description  Retrieves the user information
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerToken
// @Router       /users/info 	[get]
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
func (u *UserController) GetUserInfo(ctx *gin.Context) {
	id, exist := ctx.Get("id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "token does not exist",
		})
		return
	}

	idInt, ok := id.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "invalid payload format",
		})
		return
	}

	result, err := u.us.GetUserInfo(ctx, idInt)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, app_error.UserNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Code:    "USER_NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Result:  result,
		Message: "successfully get user profile",
	})
}
