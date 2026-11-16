package repo

import (
	"context"
	"errors"

	app_error "github.com/alfannurfaiz9/shrink_url_be/internal/error"
	"github.com/alfannurfaiz9/shrink_url_be/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LinkRepo struct {
	db *pgxpool.Pool
}

func NewLinksRepo(db *pgxpool.Pool) *LinkRepo {
	return &LinkRepo{
		db: db,
	}
}

func (l *LinkRepo) CheckExistingLinks(ctx context.Context, slug string) (bool, error) {
	sql := `
	SELECT EXISTS (
		SELECT 1 from links 
		WHERE slug = $1
	)
	`
	var exist bool
	if err := l.db.QueryRow(ctx, sql, slug).Scan(
		&exist,
	); err != nil {
		return false, err
	}

	return exist, nil
}

func (l *LinkRepo) CreateLink(ctx context.Context, body model.Link) error {
	sql := `
	INSERT INTO links(user_id, original_url, slug)
	VALUES($1, $2, $3)
	`
	args := []any{body.UserId, body.OriginalUrl, body.Slug}

	_, err := l.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (l *LinkRepo) GetAllUserLinks(ctx context.Context, userId int, search string) (int, error) {
	sql := `
	SELECT COUNT(slug)
	FROM links
	WHERE user_id = $1
	AND slug LIKE '%' || $2 || '%'
	AND deleted_at IS NULL
	`
	var total int
	if err := l.db.QueryRow(ctx, sql, userId, search).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (l *LinkRepo) GetUserLinks(ctx context.Context, userId int, body model.GetUserLinksParams) ([]model.Link, error) {
	sql := `
	SELECT id, slug, original_url, created_at
	FROM links
	WHERE user_id = $1 AND deleted_at IS NULL
	AND slug LIKE '%' || $2 || '%'
	ORDER BY created_at DESC
	LIMIT $3 OFFSET $4`
	args := []any{userId, body.Search, body.Limit, body.Offset}
	rows, err := l.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	links := make([]model.Link, 0)
	for rows.Next() {
		var link model.Link
		if err := rows.Scan(
			&link.Id,
			&link.Slug,
			&link.OriginalUrl,
			&link.CreatedAt,
		); err != nil {
			return nil, err
		}
		links = append(links, link)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}

func (l *LinkRepo) GetOriginalUrl(ctx context.Context, slug string) (string, error) {
	sql := `
	SELECT original_url
	FROM links
	WHERE slug = $1
	AND deleted_at IS NULL`

	var originalUrl string
	if err := l.db.QueryRow(ctx, sql, slug).Scan(
		&originalUrl,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", app_error.LinkNotFound
		}
		return "", err
	}

	return originalUrl, nil
}

func (l *LinkRepo) DeleteSlug(ctx context.Context, urlId, userId int) error {
	sql := `
	UPDATE links
	SET deleted_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND deleted_at IS NULL`
	cmd, err := l.db.Exec(ctx, sql, urlId, userId)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return app_error.LinkNotFound
	}

	return nil
}

func (l *LinkRepo) GetSlugById(ctx context.Context, id int) (string, error) {
	sql := `
	SELECT slug
	FROM links
	WHERE id = $1
	AND deleted_at IS NULL
	`
	var slug string
	if err := l.db.QueryRow(ctx, sql, id).Scan(&slug); err != nil {
		return "", err
	}

	return slug, nil
}
