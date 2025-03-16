package models

type DiaryEntry struct {
	Id           uint                 `json:"id"`
	Title        string               `json:"title"`
	Content      string               `json:"content"`
	Registration ActivityRegistration `json:"activityRegistration"`
}
