package models

type User struct {
	Username string `json:"username,omitempty" bson:"username"`
	Email string `json:"email,omitempty" bson:"email"`
}