package models

import "time"

type Character struct {
	ID        uint32 `json:"id" gorm:"primaryKey"`
	AccountID uint32 `json:"accountId" gorm:"column:account_id;not null"`
	Name      string `json:"name"`
	Molding   string `json:"molding"` // 造型
	Sect      string `json:"sect"`    // 门派
	Level     string `json:"level"`
	Remark    string `json:"remark"`

	Account   Account   `json:"account" gorm:"foreignKey:AccountID;"`
	Equipment Equipment `json:"equipment" gorm:"foreignKey:CharacterID;"`
	Pets      []Pet     `json:"pets" gorm:"foreignKey:CharacterID;"`

	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
