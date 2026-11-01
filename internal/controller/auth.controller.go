package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/service"
	"github.com/alfannurfaiz9/shrink_url_be/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	as *service.AuthService
}

func NewAuthController(as *service.AuthService) *AuthController {
	return &AuthController{
		as: as,
	}
}

// Register
// @Summary Register a new user
// @Description Register a new user account with email, full name, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.Register true "Register Payload"
// @Router /auth/register [post]
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.Register
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		log.Println(err.Error())
		errMsg := utils.ValidatorError(err)
		if errMsg != "" {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Code:    "VALIDATION_ERROR",
				Message: errMsg,
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Code:    "INVALID_PAYLOAD",
			Message: "invalid payload format",
		})
		return
	}

	if err := a.as.Register(ctx, body); err != nil {
		log.Println(err.Error())
		if errors.Is(err, app_error.EmailAlreadyExist) {
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{
				Success: false,
				Code:    "EMAIL_ALREADY_EXIST",
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

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "account successfully created",
	})
}
