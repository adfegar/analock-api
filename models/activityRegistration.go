package models

type ActivityRegistration struct {
	Id               uint  `json:"id"`
	RegistrationDate int64 `json:"registrationDate"`
	UserRefer        uint  `json:"userId"`
}
