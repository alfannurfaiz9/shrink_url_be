package dto

type Response struct {
	Success bool   `json:"success" example:"true"`
	Result  any    `json:"result,omitempty"`
	Message string `json:"message" example:"successfully get data"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Code    string `json:"code" example:"EMAIL_ALREADY_EXISTS"`
	Message string `json:"message" example:"internal server error"`
}
