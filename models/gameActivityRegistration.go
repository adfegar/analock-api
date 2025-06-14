package models

type GameActivityRegistration struct {
	Id           uint                 `json:"id"`
	Registration ActivityRegistration `json:"registration"`
	GameName     string               `json:"gameName"`
}
