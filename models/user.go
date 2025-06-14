package models

type UserRole int

const (
	Admin UserRole = iota + 1
	Standard
)

type User struct {
	Id       uint     `json:"id"`
	Email    string   `json:"email"`
	UserName string   `json:"userName"`
	Role     UserRole `json:"role"`
}
