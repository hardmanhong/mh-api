package models

import "time"

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Password string
	Salt     string
}

type Token struct {
	ID        uint64 `gorm:"primaryKey"`
	UserId    uint64 `gorm:"column:user_id;not null"`
	Token     string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	ExpireAt  time.Time
}
