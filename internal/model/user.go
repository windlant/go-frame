package model

type User struct {
	ID        int    `json:"id"         orm:"id,primary"`
	Name      string `json:"name"       v:"required"       orm:"name"`
	Email     string `json:"email"      v:"required|email" orm:"email,unique"`
	CreatedAt string `json:"created_at" orm:"created_at"`
}
