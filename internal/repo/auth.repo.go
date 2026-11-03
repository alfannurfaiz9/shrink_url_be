package repo

import (
	"context"

	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	db *pgxpool.Pool
}

func NewAuthRepo(db *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{
		db: db,
	}
}

func (a *AuthRepo) IsUserExist(ctx context.Context, email string) (bool, error) {
	sql := `
	SELECT EXISTS (
		SELECT 1 FROM users 
		WHERE lower(email) = LOWER($1)
	)`
	args := []any{email}

	var exist bool
	if err := a.db.QueryRow(ctx, sql, args...).Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (a *AuthRepo) Register(ctx context.Context, body model.User) error {
	sql := `
	INSERT INTO users(email, full_name, password_hash)
	VALUES($1, $2, $3)
	`
	args := []any{body.Email, body.FullName, body.PasswordHash}

	_, err := a.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthRepo) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `
	SELECT id, email, full_name, password_hash
	FROM users
	WHERE lower(email) = lower($1)`
	args := []any{email}

	var user model.User
	if err := a.db.QueryRow(ctx, sql, args...).Scan(
		&user.Id,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
	); err != nil {
		return model.User{}, err
	}

	return user, nil
}
