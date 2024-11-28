package models

type Song struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Group string `json:"group"`
	Title string `json:"song"`
	SongDetail
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
