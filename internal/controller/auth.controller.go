package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

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
// @Param request body dto.RegisterRequest true "Register Payload"
// @Router /auth/register [post]
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
func (a *AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterRequest
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

// Login
// @Summary      User Login
// @Description  Authenticates user credentials and returns a JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login Payload"
// @Router       /auth/login [post]
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
func (a *AuthController) Login(ctx *gin.Context) {
	var body dto.LoginRequest
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

	result, err := a.as.Login(ctx.Request.Context(), body)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, app_error.InvalidEmailOrPassword) {
			ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Success: false,
				Code:    "INVALID_CREDENTIALS",
				Message: "email or password is invalid",
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
		Message: "successfully log in",
	})
}

// Logout
// @Summary      User Logout
// @Description  Revokes the current JWT token by blacklisting its JTI in Redis.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerToken
// @Router       /auth/logout [post]
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
func (a *AuthController) Logout(ctx *gin.Context) {
	jwtId, exist := ctx.Get("ID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "unauthorized request",
		})
		return
	}
	jwtIdStr, ok := jwtId.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Code:    "INVALID_PAYLOAD",
			Message: "invalid jwt id type",
		})
		return
	}

	expToken, exist := ctx.Get("exp")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "unauthorized request",
		})
		return
	}
	expTokenVal, ok := expToken.(time.Time)
	if !ok {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Code:    "INVALID_PAYLOAD",
			Message: "invalid exp token type",
		})
		return
	}

	if err := a.as.Logout(ctx.Request.Context(), jwtIdStr, expTokenVal); err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "successfully logged out",
	})
}
