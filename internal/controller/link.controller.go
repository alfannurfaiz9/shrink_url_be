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

type LinkController struct {
	ls *service.LinkService
}

// CreateLink
// @Summary      Create a new link
// @Description  Create a new short link with an option for a custom slug or auto-generated slug.
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateLinkRequest true "Create Link Payload"
// @Security     BearerToken
// @Router       /links [post]
// @Success      201  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
func NewLinkController(ls *service.LinkService) *LinkController {
	return &LinkController{
		ls: ls,
	}
}

func (l *LinkController) CreateLink(ctx *gin.Context) {
	var body dto.CreateLinkRequest

	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
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

	userId, _ := ctx.Get("id")
	userIdInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "invalid payload format",
		})
		return
	}

	if err := l.ls.CreateLink(ctx, body, userIdInt); err != nil {
		log.Println(err.Error())

		if errors.Is(err, app_error.InvalidURL) || errors.Is(err, app_error.InvalidSlug) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, app_error.SlugAlreadyExist) {
			ctx.JSON(http.StatusConflict, dto.ErrorResponse{
				Success: false,
				Code:    "SLUG_ALREADY_EXIST",
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
		Message: "link created successfully",
	})
}
