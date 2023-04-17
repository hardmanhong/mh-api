package models

import "time"

type Buy struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2)"`
	Quantity    int       `json:"quantity" gorm:"column:quantity"`
	Inventory   int       `json:"inventory"`
	TotalAmount float64   `json:"totalAmount" gorm:"column:total_amount"`
	TotalProfit float64   `json:"totalProfit" gorm:"column:total_profit"`
	CreatedAt   time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Remark      string    `json:"remark"`
	HasSold     int       `json:"hasSold" gorm:"column:has_sold"`
	GoodsID     uint64    `json:"goodsId" gorm:"column:goods_id"`
	Goods       Goods     `json:"goods" gorm:"foreignKey:GoodsID;preload"`
	Sales       []Sell    `json:"sales" gorm:"foreignKey:BuyID"`
}

type BuyUpdate struct {
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity" gorm:"column:quantity"`
	Inventory int
	Remark    string    `json:"remark"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type BuyUpdateProfit struct {
	Inventory   int     `json:"inventory"`
	TotalProfit float64 `json:"totalProfit" gorm:"column:total_profit"`
}

type BuyListQuery struct {
	GoodsIDs        []uint64   `json:"goodsIds"`
	InventorySorter string     `json:"inventorySorter"`
	CreatedAtFrom   *time.Time `json:"createdAtFrom"`
	CreatedAtTo     *time.Time `json:"createdAtTo"`
	PaginationQuery
}

type BuyListResponse struct {
	TotalAmount float64 `json:"totalAmount"`
	TotalProfit float64 `json:"totalProfit"`
	PaginationResponse
}
