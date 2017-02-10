package models

type UnregisteredUser struct {
	Username string `json:"username,omitempty" bson:"username"`
	Password string `json:"password,omitempty" bson:"password"`
	HashedPassword []byte `json:"hashedPassword,omitempty" bson:"hashedPassword"`
	Email string `json:"email,omitempty" bson:"email"`
}