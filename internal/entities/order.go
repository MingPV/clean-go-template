package entities

type Order struct {
	ID    uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	Total float64 `json:"total"`
}
