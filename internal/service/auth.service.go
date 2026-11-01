package service

import (
	"context"
	"strings"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/alfannurfaiz9/shrink_url_be/internal/utils"
	"github.com/alfannurfaiz9/shrink_url_be/pkg"
)

type AuthService struct {
	ar *repo.AuthRepo
}

func NewAuthService(ar *repo.AuthRepo) *AuthService {
	return &AuthService{
		ar: ar,
	}
}

func (a *AuthService) Register(ctx context.Context, body dto.Register) error {
	exist, err := a.ar.IsUserExist(ctx, body.Email)
	if err != nil {
		return app_error.ErrConnectDatabase
	}
	if exist {
		return app_error.EmailAlreadyExist
	}

	hash := pkg.NewHashConfig()
	hashedPass := hash.GenerateHash(body.Password)
	email := strings.ToLower(strings.TrimSpace(body.Email))

	if err := a.ar.Register(ctx, model.User{
		Email:        email,
		FullName:     body.FullName,
		PasswordHash: hashedPass,
	}); err != nil {
		if utils.IsUniqueViolation(err) {
			return app_error.EmailAlreadyExist
		}

		return err
	}

	return nil
}
