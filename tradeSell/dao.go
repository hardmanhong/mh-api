package tradeSell

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TradeSellDAO interface {
	GetItem(id uint64) (*TradeSell, error)
	Create(trade *TradeSell) (*TradeSell, error)
	Update(id uint64, trade *TradeSellUpdate) error
	Delete(id uint64) error
}

type TradeDAOImpl struct {
	db *gorm.DB
}

func NewTradeDAO(db *gorm.DB) TradeSellDAO {
	return &TradeDAOImpl{db}
}

func (dao *TradeDAOImpl) GetItem(id uint64) (*TradeSell, error) {
	trade := &TradeSell{}
	err := dao.db.Where("id = ?", id).First(trade).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("trade id=%d not found", id)
	}
	return trade, err
}

func (dao *TradeDAOImpl) Create(trade *TradeSell) (*TradeSell, error) {
	return trade, dao.db.Create(&trade).Error
}

func (dao *TradeDAOImpl) Update(id uint64, trade *TradeSellUpdate) error {
	return dao.db.Table("trade_sell").Where("id = ?", id).Updates(trade).Error
}

func (dao *TradeDAOImpl) Delete(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&TradeSell{}).Error
}
