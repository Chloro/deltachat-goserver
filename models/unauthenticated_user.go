package models

type UnauthenticatedUser struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	HashedPassword []byte `json:"hashedPassword,omitempty"`
	Email string `json:"email,omitempty"`
}