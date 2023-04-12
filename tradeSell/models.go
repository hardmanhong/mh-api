package tradeSell

import (
	"time"
)

type TradeSell struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	TradeBuyID   uint64    `json:"tradeBuyId" gorm:"column:trade_buy_id;not null"`
	GoodsID      uint64    `json:"goodsId" gorm:"column:goods_id;not null"`
	SellPrice    float64   `json:"sellPrice" gorm:"type:decimal(10,2)"`
	SellQuantity uint64    `json:"sellQuantity" gorm:"column:sell_quantity"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type TradeSellUpdate struct {
	SellPrice    *float64  `json:"sellPrice" gorm:"type:decimal(10,2)"`
	SellQuantity *uint64   `json:"sellQuantity" gorm:"column:sell_quantity"`
	Remark       *string   `json:"remark"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
