package dto

import "time"

type CreateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	Slug        string `json:"slug" binding:"omitempty,min=3,max=50"`
}

type Link struct {
	Id          int       `json:"id"`
	Slug        string    `json:"slug"`
	OriginalUrl string    `json:"original_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetUserLinksRequest struct {
	Search string `form:"search"`
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
}
