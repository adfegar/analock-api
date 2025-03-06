package models

type LoginProvider int

const (
	Google LoginProvider = iota + 1
)

type ExternalLogin struct {
	Id          uint          `json:"id"`
	Provider    LoginProvider `json:"provider"`
	ClientId    string        `json:"provider_client_id"`
	ClientToken string        `json:"provider_client_token"`
	UserRefer   uint          `json:"user_id"`
}
