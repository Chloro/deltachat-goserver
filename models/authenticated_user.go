package models

type AuthenticatedUser struct {
	Username string `json:"username" bson:"username"`
	Token string `json:"token" bson:"token"`
	Email string `json:"email" bson:"email"`
}
