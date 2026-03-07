package models

import "go.mongodb.org/mongo-driver/v2/bson"

// UserNameLower is indexed for case insensitive lookup
type User struct {
	ID            bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserName      string        `bson:"username,omitempty" json:"userName"`
	UserNameLower string        `bson:"username_lower,omitempty"`
	Email         string        `bson:"email,omitempty" json:"email"`
	PasswordHash  string        `bson:"password_hash,omitempty" json:"passwordHash"`
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
