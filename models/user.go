package models

type UserRole int

const (
	Admin UserRole = iota + 1
	Standard
)

type User struct {
	Id       uint     `json:"id"`
	UserName string   `json:"userName"`
	Role     UserRole `json:"role"`
}
