package dto

type Register struct {
	FullName string `json:"full_name" binding:"required,min=3,max=100" example:"Alex Walker"`
	Email    string `json:"email" binding:"required,email" example:"alexwalker@mail.com"`
	Password string `json:"password" binding:"required,min=8,max=72" example:"alex123"`
}
