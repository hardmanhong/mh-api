package tradeBuy

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hardmanhong/goods"
	"gorm.io/gorm"
)

type TradeDAO interface {
	GetList(page, pageSize int, filter *TradeListFilter) (*ListResult, error)
	GetItem(id uint64) (*TradeBuy, error)
	Create(trade *TradeBuy) (*TradeBuy, error)
	Update(id uint64, trade *TradeUpdate) error
	Delete(id uint64) error
}

type TradeDAOImpl struct {
	db *gorm.DB
}

func NewTradeDAO(db *gorm.DB) TradeDAO {
	return &TradeDAOImpl{db}
}

func (dao *TradeDAOImpl) GetList(page, pageSize int, filter *TradeListFilter) (*ListResult, error) {
	var total int64
	var trades []TradeBuy
	tx := dao.db.Model(&TradeBuy{}).Preload("Goods").Preload("Sales")

	if filter != nil {
		if filter.CreatedAtFrom != nil {
			tx = tx.Where("created_at >= ?", filter.CreatedAtFrom)
		}
		if filter.CreatedAtTo != nil {
			tx = tx.Where("created_at <= ?", filter.CreatedAtTo)
		}
		jsonData, _ := json.Marshal(filter)
		println("filter", string(jsonData))
		if len(filter.GoodsIDs) > 0 {
			tx = tx.Where("goods_id IN (?)", filter.GoodsIDs)
		}
	}

	err := tx.Count(&total).Error
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	err = tx.Order("id desc").Offset(offset).Limit(pageSize).Find(&trades).Error
	if err != nil {
		return nil, err
	}

	result := &ListResult{
		Total: total,
		List:  trades,
	}

	return result, nil
}

func (dao *TradeDAOImpl) GetItem(id uint64) (*TradeBuy, error) {
	trade := &TradeBuy{}
	err := dao.db.Where("id = ?", id).Preload("Sales").First(trade).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("trade id=%d not found", id)
	}
	return trade, err
}

func (dao *TradeDAOImpl) Create(trade *TradeBuy) (*TradeBuy, error) {
	// 设置 CreatedAt 和 UpdatedAt 字段的值
	now := time.Now()
	trade.CreatedAt = now
	trade.UpdatedAt = now
	trade.Stock = trade.BuyQuantity
	// trade.ProfitRate = 0

	err := dao.db.Create(trade).Error
	if err != nil {
		return nil, err
	}

	// 更新关联的 Goods 信息
	dao.db.Model(&trade).Association("Goods").Append(&goods.Goods{ID: trade.GoodsID})

	return trade, nil
}

func (dao *TradeDAOImpl) Update(id uint64, trade *TradeUpdate) error {
	return dao.db.Table("trade_buy").Where("id = ?", id).Updates(trade).Error
}

func (dao *TradeDAOImpl) Delete(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&TradeBuy{}).Error
}
