package model

type Movie struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Title    string `json:"title" binding:"required"`
	Director string `json:"director"`
	Year     int    `json:"year"`
	Plot     string `json:"plot"`
}
