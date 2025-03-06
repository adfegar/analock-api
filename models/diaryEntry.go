package models

type DiaryEntry struct {
	Id        uint   `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Date      int64  `json:"publishDate"`
	UserRefer uint   `json:"user_id"`
}
