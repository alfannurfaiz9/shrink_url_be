package dto

type Response struct {
	Success  bool           `json:"success" example:"true"`
	Result   any            `json:"result,omitempty"`
	Message  string         `json:"message" example:"successfully get data"`
	Metadata PaginationMeta `json:"metadata,omitzero"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"EMAIL_ALREADY_EXISTS"`
	Message string `json:"message" example:"internal server error"`
}

type PaginationMeta struct {
	CurrentPage int     `json:"current_page"`
	TotalPage   float64 `json:"total_page"`
	Next        string  `json:"next"`
	Prev        string  `json:"prev"`
}
