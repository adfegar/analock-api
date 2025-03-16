package models

type BookActivityRegistration struct {
	Id                        uint                 `json:"id"`
	Registration              ActivityRegistration `json:"registration"`
	InternetArchiveIdentifier string               `json:"internetArchiveId"`
}
