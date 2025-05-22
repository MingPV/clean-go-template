package entities

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
