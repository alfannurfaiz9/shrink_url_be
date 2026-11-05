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
