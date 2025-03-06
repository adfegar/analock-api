package models

type TokenKind int

const (
	Access TokenKind = iota + 1
	Refresh
)

type Token struct {
	Id         uint   `json:"id"`
	TokenValue string `json:"token"`
	UserRefer  uint   `json:"user_id"`
	Kind       TokenKind
}
