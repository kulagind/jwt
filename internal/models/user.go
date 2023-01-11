package models

import "time"

type User struct {
	Id        string    `json:"id" sql:"id"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	Name      string    `json:"name" sql:"name"`
	TokenHash string    `json:"tokenhash" sql:"tokenhash"`
	CreatedAt time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" sql:"updated_at"`
}

type UserContextToken struct{}

type UserResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (u *User) Valid() bool {
	if u.Email == "" || u.Password == "" {
		return false
	}
	return true
}
