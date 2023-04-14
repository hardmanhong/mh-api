package models

import (
	"time"
)

type Sell struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	BuyID       uint64    `json:"buyId" gorm:"column:buy_id;not null"`
	GoodsID     uint64    `json:"goodsId" gorm:"column:goods_id;not null"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
	Quantity    int       `json:"quantity" gorm:"column:quantity"`
	Remark      string    `json:"remark"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Profit      float64   `json:"profit"`
	TotalProfit float64   `json:"totalProfit"`
}

type SellUpdate struct {
	Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
	Quantity    int       `json:"quantity" gorm:"column:quantity"`
	Remark      string    `json:"remark"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Profit      float64   `json:"profit"`
	TotalProfit float64   `json:"totalProfit"`
}
