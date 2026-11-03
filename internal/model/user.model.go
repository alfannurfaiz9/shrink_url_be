package model

import "time"

type User struct {
	Id           int       `db:"id"`
	Email        string    `db:"email"`
	FullName     string    `db:"full_name"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
}

type UserProfile struct {
	User
	ActiveLink int `db:"active_link"`
}
