package model

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" v:"required"`
	Email string `json:"email" v:"required|email"`
}
