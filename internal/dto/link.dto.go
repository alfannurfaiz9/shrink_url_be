package dto

type CreateLinkRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	Slug        string `json:"slug" binding:"omitempty,min=3,max=50"`
}
