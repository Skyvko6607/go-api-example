package models

type LoginDTO struct {
	UserNameOrEmail string `json="userNameOrEmail"`
	Password        string `json="password"`
}
