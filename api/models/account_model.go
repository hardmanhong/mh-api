package models

import "time"

type Account struct {
	UserId    uint64    `json:"-" gorm:"column:user_id"`
	ID        uint32    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"column:name"`
	Server    string    `json:"server" gorm:"column:server"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type AccountListQuery struct {
	Name string `json:"name"`
}

type AccountUpdate struct {
	Name   string `json:"name"`
	Server string `json:"server" gorm:"column:server"`
}
