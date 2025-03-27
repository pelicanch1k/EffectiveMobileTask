package model

type Song struct {
	ID          int    `json:"id" db:"id"`
	Genre       string `json:"genre" db:"genre"`
	Song        string `json:"song" db:"song"`
	ReleaseDate string `json:"releaseDate" db:"releaseDate"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
	GroupID     int    `json:"groupId" db:"group_id"`
	Group       string `json:"group" db:"group_name"`
}
