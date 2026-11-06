package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/alfannurfaiz9/shrink_url_be/internal/utils"
	"github.com/alfannurfaiz9/shrink_url_be/pkg"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	ar  *repo.AuthRepo
	rdb *redis.Client
}

func NewAuthService(ar *repo.AuthRepo, rdb *redis.Client) *AuthService {
	return &AuthService{
		ar:  ar,
		rdb: rdb,
	}
}

func (a *AuthService) Register(ctx context.Context, body dto.RegisterRequest) error {
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

func (a *AuthService) Login(ctx context.Context, body dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := a.ar.FindUserByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.LoginResponse{}, app_error.InvalidEmailOrPassword
		}
		return dto.LoginResponse{}, err
	}

	if err := pkg.CompareHash(body.Password, user.PasswordHash); err != nil {
		return dto.LoginResponse{}, app_error.InvalidEmailOrPassword
	}

	claims := pkg.NewJWTClaims(user.Id)
	token, err := claims.GenerateToken()
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{Token: token}, nil
}

func (a *AuthService) Logout(ctx context.Context, jwtID string, expToken time.Time) error {
	exp := time.Until(expToken)
	key := fmt.Sprintf("shrinkurl:token:%s", jwtID)

	if err := utils.SetRedisStr(ctx, a.rdb, key, "revoked", exp); err != nil {
		return err
	}

	return nil
}
