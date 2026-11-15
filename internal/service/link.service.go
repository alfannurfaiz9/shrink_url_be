package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/alfannurfaiz9/shrink_url_be/internal/utils"
	"github.com/redis/go-redis/v9"
)

type LinkService struct {
	lr  *repo.LinkRepo
	rdb *redis.Client
}

func NewLinkService(lr *repo.LinkRepo, rdb *redis.Client) *LinkService {
	return &LinkService{
		lr:  lr,
		rdb: rdb,
	}
}

func (l *LinkService) CheckExistingLinks(ctx context.Context, slug string) (bool, error) {
	if slug == "" {
		return false, app_error.LinkNotFound
	}

	slugRegex := regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)
	if !slugRegex.MatchString(slug) {
		return false, app_error.InvalidSlug
	}

	exist, err := l.lr.CheckExistingLinks(ctx, slug)
	if err != nil {
		return false, err
	}

	return exist, nil
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
	key := fmt.Sprintf("shrinkurl:%s", slug)
	if cache, err := utils.GetRedisStr(ctx, l.rdb, key); err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Printf("failed to get %s cache %s", key, err)
		}
	} else {
		return cache, nil
	}

	url, err := l.lr.GetOriginalUrl(ctx, slug)
	if err != nil {
		return "", err
	}

	if err := utils.SetRedisStr(ctx, l.rdb, key, url, 24*time.Hour); err != nil {
		log.Printf("failed to set %s cache %s", key, err)
	}

	return url, nil
}

func (l *LinkService) DeleteSlug(ctx context.Context, urlId, userId int) error {
	if urlId <= 0 || userId <= 0 {
		return app_error.InvalidSlugPayload
	}

	if slug, err := l.lr.GetSlugById(ctx, urlId); err != nil {
		log.Println("failed to delete cache, slug not found")
	} else {
		key := fmt.Sprintf("shrinkurl:%s", slug)
		if err := utils.DelRedis(ctx, l.rdb, key); err != nil {
			log.Printf("failed to delete %s cache %s", key, err)
		}
	}

	if err := l.lr.DeleteSlug(ctx, urlId, userId); err != nil {
		return err
	}

	return nil
}
