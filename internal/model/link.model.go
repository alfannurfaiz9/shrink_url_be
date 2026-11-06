package model

import "time"

type Link struct {
	Id          int        `db:"id"`
	UserId      int        `db:"user_id"`
	OriginalUrl string     `db:"original_url"`
	Slug        string     `db:"slug"`
	CreatedAt   time.Time  `db:"created_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

type GetUserLinksParams struct {
	UserId int    `db:"user_id"`
	Search string `db:"search"`
	Limit  int    `db:"limit"`
	Offset int    `db:"offset"`
}
