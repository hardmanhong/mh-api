package models

import "time"

type Pet struct {
	ID          uint32    `json:"id" gorm:"primaryKey"`
	CharacterID uint32    `json:"-" gorm:"column:character_id;not null"`
	Name        string    `json:"name"`
	Price       string    `json:"price"`
	Level       string    `json:"level"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
