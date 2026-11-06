package controller

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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

func NewLinkController(ls *service.LinkService) *LinkController {
	return &LinkController{
		ls: ls,
	}
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

	userId, exist := ctx.Get("id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "unauthorized request",
		})
		return
	}
	userIdInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "invalid user id type",
		})
		return
	}

	if err := l.ls.CreateLink(ctx.Request.Context(), body, userIdInt); err != nil {
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

// GetUserLinks
// @Summary      Get user links
// @Description  Retrieves the list of links owned by the currently logged-in user based on their JWT token
// @Tags         Links
// @Accept       json
// @Produce      json
// @Security     BearerToken
// @Param        search	query	string	false  "search"
// @Param        limit	query	int		false  "limit"
// @Param        page	query	int		false  "page"
// @Router       /links [get]
// @Success      200  {object}  dto.Response
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
func (l *LinkController) GetUserLinks(ctx *gin.Context) {
	var query dto.GetUserLinksRequest

	if err := ctx.ShouldBindWith(&query, binding.Query); err != nil {
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
			Code:    "INVALID_QUERY_PARAMS",
			Message: "invalid query params",
		})
		return
	}

	userId, exist := ctx.Get("id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "unauthorized request",
		})
		return
	}

	userIdInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "invalid user id type",
		})
		return
	}

	res, meta, err := l.ls.GetUserLinks(ctx.Request.Context(), userIdInt, ctx.Request.URL.Path, query)
	if err != nil {
		log.Println(err)

		if errors.Is(err, app_error.LinkPageNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Code:    "PAGE_NOT_FOUND",
				Message: "link page not found",
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

	data := dto.Response{
		Success: true,
		Result:  res,
		Message: "successfully get user links",
	}

	if len(res) > 0 {
		data.Metadata = meta
	}

	ctx.JSON(http.StatusOK, data)
}

// GetOriginalUrl
// @Summary      Redirect to Original URL
// @Description  Fetches the original URL based on the provided slug and redirects the client to it.
// @Tags         Links
// @Param        slug   path      string  true  "Unique Link Slug"
// @Router       /links/{slug}	[get]
// @Success      302    {string}  string  "Successfully redirected to the original URL"
// @Failure      404    {object}  dto.ErrorResponse
// @Failure      500    {object}  dto.ErrorResponse
func (l *LinkController) GetOriginalUrl(ctx *gin.Context) {
	slug := ctx.Param("slug")

	res, err := l.ls.GetOriginalUrl(ctx.Request.Context(), slug)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, app_error.LinkNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Code:    "LINK_NOT_FOUND",
				Message: "link not found",
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

	ctx.Redirect(http.StatusFound, res)
}

// DeleteSlug
// @Summary      Delete a link (Soft Delete)
// @Description  Soft deletes a specific short link by its ID, ensuring the link belongs to the authenticated user.
// @Tags         Links
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Link ID"
// @Security     BearerToken
// @Router       /links/{id} [delete]
// @Success      200  {object}  dto.Response
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
func (l *LinkController) DeleteSlug(ctx *gin.Context) {
	userId, exist := ctx.Get("id")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Code:    "UNAUTHORIZED",
			Message: "unauthorized request",
		})
		return
	}
	userIdInt, ok := userId.(int)
	if !ok {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Code:    "VALIDATION_ERROR",
			Message: "invalid user id type",
		})
		return
	}

	urlId := ctx.Param("id")
	urlIdInt, err := strconv.Atoi(urlId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Code:    "VALIDATION_ERROR",
			Message: "invalid url id type",
		})
		return
	}

	if err := l.ls.DeleteSlug(ctx.Request.Context(), urlIdInt, userIdInt); err != nil {
		log.Println(err.Error())

		if errors.Is(err, app_error.LinkNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Code:    "SLUG_NOT_FOUND",
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, app_error.InvalidSlugPayload) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Code:    "VALIDATION_ERROR",
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
		Message: "slug deleted successfully",
	})
}
