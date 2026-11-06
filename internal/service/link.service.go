package service

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/alfannurfaiz9/shrink_url_be/internal/utils"
)

type LinkService struct {
	lr *repo.LinkRepo
}

func NewLinkService(lr *repo.LinkRepo) *LinkService {
	return &LinkService{
		lr: lr,
	}
}

func (l *LinkService) CreateLink(ctx context.Context, req dto.CreateLinkRequest, userId int) error {
	parsedURL, err := url.ParseRequestURI(req.OriginalUrl)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return app_error.InvalidURL
	}

	isCustomSlug := req.Slug != ""
	slugRegex := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !isCustomSlug {
		req.Slug = utils.GenerateSlug()
	} else {
		if !slugRegex.MatchString(req.Slug) {
			return app_error.InvalidSlug
		}
	}

	if isCustomSlug {
		exist, err := l.lr.CheckExistingLinks(ctx, req.Slug)
		if err != nil {
			return err
		}
		if exist {
			return app_error.SlugAlreadyExist
		}
	}

	return l.lr.CreateLink(ctx, model.Link{
		UserId:      userId,
		OriginalUrl: req.OriginalUrl,
		Slug:        req.Slug,
	})
}

func (l *LinkService) GetUserLinks(ctx context.Context, userId int, endpoint string, query dto.GetUserLinksRequest) ([]dto.Link, dto.PaginationMeta, error) {
	if query.Limit <= 0 {
		query.Limit = 6
	}
	if query.Page <= 0 {
		query.Page = 1
	}

	res, err := l.lr.GetUserLinks(ctx, userId, model.GetUserLinksParams{
		Search: query.Search,
		Limit:  query.Limit,
		Offset: (query.Page - 1) * query.Limit,
	})
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	totalData, err := l.lr.GetAllUserLinks(ctx, userId, query.Search)
	if err != nil {
		return nil, dto.PaginationMeta{}, err
	}

	var totalPage float64
	if totalData > 0 {
		totalPage = math.Ceil(float64(totalData) / float64(query.Limit))
		if query.Page > int(totalPage) {
			return nil, dto.PaginationMeta{}, app_error.LinkPageNotFound
		}
	}

	links := make([]dto.Link, 0, len(res))
	for _, v := range res {
		links = append(links, dto.Link{
			Id:          v.Id,
			Slug:        v.Slug,
			OriginalUrl: v.OriginalUrl,
			CreatedAt:   v.CreatedAt,
		})
	}

	var next string
	if query.Page < int(totalPage) {
		q := url.Values{}
		q.Set("limit", strconv.Itoa(query.Limit))
		q.Set("page", strconv.Itoa(query.Page+1))
		next = fmt.Sprintf("%s?%s", endpoint, q.Encode())
	}

	var prev string
	if query.Page != 1 {
		q := url.Values{}
		q.Set("limit", strconv.Itoa(query.Limit))
		q.Set("page", strconv.Itoa(query.Page-1))
		prev = fmt.Sprintf("%s?%s", endpoint, q.Encode())
	}

	metadata := dto.PaginationMeta{
		CurrentPage: query.Page,
		TotalPage:   totalPage,
		Prev:        prev,
		Next:        next,
	}

	return links, metadata, nil
}

func (l *LinkService) GetOriginalUrl(ctx context.Context, slug string) (string, error) {
	return l.lr.GetOriginalUrl(ctx, slug)
}

func (l *LinkService) DeleteSlug(ctx context.Context, urlId, userId int) error {
	if urlId <= 0 || userId <= 0 {
		return app_error.InvalidSlugPayload
	}

	return l.lr.DeleteSlug(ctx, urlId, userId)
}
