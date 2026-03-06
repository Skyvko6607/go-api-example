package models

import "go.mongodb.org/mongo-driver/v2/bson"

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty" json:"id"`

	UserName      string `bson:"userName,omitempty" json:"userName"`
	UserNameLower string `bson:"userNameLower,omitempty"` // Index

	Email      string `bson:"email,omitempty" json:"email"`
	EmailLower string `bson:"emailLower,omitempty"` // Index
}

type UserDTO struct {
	UserName string `json:"userName"`
	Email    string `json:"email"`
}

func (u *User) AsDTO() UserDTO {
	return UserDTO{
		UserName: u.UserName,
		Email:    u.Email,
	}
}
