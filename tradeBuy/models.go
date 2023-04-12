package tradeBuy

import (
	"time"

	"github.com/hardmanhong/goods"
	"github.com/hardmanhong/tradeSell"
)

type TradeBuy struct {
	ID          uint64                `json:"id" gorm:"primaryKey"`
	BuyPrice    float64               `json:"buyPrice" gorm:"type:decimal(10,2)"`
	BuyQuantity uint64                `json:"buyQuantity" gorm:"column:buy_quantity"`
	Stock       uint64                `json:"stock"`
	CreatedAt   time.Time             `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt   time.Time             `json:"updatedAt" gorm:"autoUpdateTime"`
	Remark      string                `json:"remark"`
	GoodsID     uint64                `json:"goodsId" gorm:"column:goods_id"`
	Goods       goods.Goods           `json:"goods" gorm:"foreignKey:GoodsID;preload"`
	Sales       []tradeSell.TradeSell `json:"sales" gorm:"foreignKey:TradeBuyID"`
}
type TradeUpdate struct {
	BuyPrice    *float64  `json:"buyPrice"`
	BuyQuantity *uint64   `json:"buyQuantity" gorm:"column:buy_quantity"`
	Remark      *string   `json:"remark"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

type ListResult struct {
	List  []TradeBuy `json:"list"`
	Total int64      `json:"total"`
}

type TradeListFilter struct {
	GoodsIDs      []uint64   `json:"goodsIds"`
	CreatedAtFrom *time.Time `json:"createdAtFrom"`
	CreatedAtTo   *time.Time `json:"createdAtTo"`
}
