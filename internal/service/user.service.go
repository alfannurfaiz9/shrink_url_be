package service

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/jackc/pgx/v5"
)

type UserService struct {
	ur *repo.UserRepo
}

func NewUserService(ur *repo.UserRepo) *UserService {
	return &UserService{
		ur: ur,
	}
}

func (u *UserService) GetUserProfile(ctx context.Context, id int) (dto.UserProfile, error) {
	result, err := u.ur.GetUserProfile(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserProfile{}, app_error.UserNotFound
		}
		return dto.UserProfile{}, err
	}

	return dto.UserProfile{
		Id:         result.Id,
		Email:      result.Email,
		FullName:   result.FullName,
		CreatedAt:  result.CreatedAt,
		ActiveLink: result.ActiveLink,
	}, nil
}
