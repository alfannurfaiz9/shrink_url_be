package dto

import "time"

type UserInfo struct {
	Email    string `json:"email"`
	FullName string `json:"FullName" json:"full_name"`
}

type UserProfile struct {
	Id         int       `json:"id"`
	Email      string    `json:"email"`
	FullName   string    `json:"FullName" json:"full_name"`
	CreatedAt  time.Time `json:"CreatedAt" json:"created_at"`
	ActiveLink int       `json:"ActiveLink" json:"active_link"`
}
