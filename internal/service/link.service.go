package service

import (
	"context"
	"net/url"
	"regexp"

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

func (l *LinkService) GetUserLink(ctx context.Context, userId int) ([]dto.Link, error) {
	res, err := l.lr.GetUserLink(ctx, userId)
	if err != nil {
		return nil, err
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

	return links, nil
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
