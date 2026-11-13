package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/alfannurfaiz9/shrink_url_be/internal/dto"
	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/repo"
	"github.com/alfannurfaiz9/shrink_url_be/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

type UserService struct {
	ur  *repo.UserRepo
	rdb *redis.Client
}

func NewUserService(ur *repo.UserRepo, rdb *redis.Client) *UserService {
	return &UserService{
		ur:  ur,
		rdb: rdb,
	}
}

func (u *UserService) GetUserInfo(ctx context.Context, id int) (dto.UserInfo, error) {
	key := fmt.Sprintf("shrinkurl:user:info:%v", id)
	cache, err := utils.GetRedisJson[dto.UserInfo](ctx, u.rdb, key)
	if err == nil {
		return cache, nil
	}
	if !errors.Is(err, redis.Nil) {
		log.Printf("failed to get %s cache %s", key, err)
	}

	result, err := u.ur.GetUserInfo(ctx, id)
	if err != nil {
		return dto.UserInfo{}, err
	}

	userInfo := dto.UserInfo{
		FullName: result.FullName,
		Email:    result.Email,
	}

	if err := utils.SetRedisJson(ctx, u.rdb, key, userInfo, 24*time.Hour); err != nil {
		log.Printf("failed to set %s cache %s", key, err)
	}

	return userInfo, nil
}

func (u *UserService) GetUserProfile(ctx context.Context, id int) (dto.UserProfile, error) {
	key := fmt.Sprintf("shrinkurl:user:profile:%v", id)
	cache, err := utils.GetRedisJson[dto.UserProfile](ctx, u.rdb, key)
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			log.Printf("failed to get %s cache %s", key, err)
		}
	}
	if err == nil {
		return cache, nil
	}

	result, err := u.ur.GetUserProfile(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.UserProfile{}, app_error.UserNotFound
		}
		return dto.UserProfile{}, err
	}

	userProfile := dto.UserProfile{
		Id:         result.Id,
		Email:      result.Email,
		FullName:   result.FullName,
		CreatedAt:  result.CreatedAt,
		ActiveLink: result.ActiveLink,
	}

	if err := utils.SetRedisJson(ctx, u.rdb, key, userProfile, 24*time.Hour); err != nil {
		log.Printf("failed to set %s cache %s", key, err)
	}

	return userProfile, nil
}
