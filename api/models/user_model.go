package models

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string
	Salt     string
}
