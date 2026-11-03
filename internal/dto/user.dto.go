package dto

import "time"

type UserProfile struct {
	Id         int       `json:"id"`
	Email      string    `json:"email"`
	FullName   string    `json:"full_name"`
	CreatedAt  time.Time `json:"created_at"`
	ActiveLink int       `json:"active_link"`
}
