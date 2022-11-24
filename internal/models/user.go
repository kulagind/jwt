package models

import "time"

type User struct {
	Id        string    `json:"id" validate:"required" sql:"id"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	Username  string    `json:"username" validate:"required" sql:"name"`
	TokenHash string    `json:"tokenhash" sql:"tokenhash"`
	CreatedAt time.Time `json:"createdat" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedat" sql:"updated_at"`
}
