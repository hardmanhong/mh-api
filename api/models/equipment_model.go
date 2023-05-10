package models

import "time"

// 游戏装备
type Equipment struct {
	ID          uint32    `json:"id" gorm:"primaryKey"`
	CharacterID uint32    `json:"-" gorm:"column:character_id;not null"`
	Arms        string    `json:"arms"`     // 武器
	Helmet      string    `json:"helmet"`   // 头盔
	Necklace    string    `json:"necklace"` // 项链
	Clothes     string    `json:"clothes"`  // 衣服
	Belt        string    `json:"belt"`     // 腰带
	Shoe        string    `json:"shoe"`     // 鞋子
	Ring        string    `json:"ring"`     // 戒指
	Bracelet    string    `json:"bracelet"` // 手镯
	Earring     string    `json:"earring"`  // 耳饰
	Trimming    string    `json:"trimming"` // 佩饰
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
