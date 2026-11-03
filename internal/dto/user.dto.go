package dto

type User struct {
	Id           int    `json:"id"`
	Email        string `json:"email"`
	FullName     string `json:"full_name"`
	PasswordHash string `json:"password_hash"`
}
