package services

import (
	"errors"

	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
)

type SellService interface {
	GetItem(id uint64) (*models.Sell, error)
	Create(sell *models.Sell) (*models.Sell, error)
	Exists(id uint64) (bool, error)
	Update(id uint64, sell *models.SellUpdate) error
	Delete(id uint64) error
}

type sellService struct {
	dao    *dao.SellDAO
	buyDao *dao.BuyDAO
}

func NewSellService(dao *dao.SellDAO, buyDao *dao.BuyDAO) *sellService {
	return &sellService{dao, buyDao}
}
func (s *sellService) GetItem(id uint64) (*models.Sell, error) {
	sell, err := s.dao.GetItem(id)
	if err != nil {
		return nil, err
	}
	return sell, nil
}

func (s *sellService) Create(sell *models.Sell) (*models.Sell, error) {
	db := s.dao.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Get the corresponding buy record
	buy, err := s.buyDao.GetItem(sell.BuyID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// Check if there is enough stock to sell
	if sell.Price == 0 {
		return nil, errors.New("请填写卖出价")
	}
	if sell.Quantity == 0 {
		return nil, errors.New("请填写卖出数量")
	}
	if buy.Inventory < sell.Quantity {
		tx.Rollback()
		return nil, errors.New("卖出数量超过买入数量")
	}
	// Update the stock of the corresponding buy record
	sellTotalProfit := (sell.Price - buy.Price) * float64(sell.Quantity) // 此次卖出利润
	inventory := buy.Inventory - sell.Quantity
	println("inventory", inventory, buy.Inventory, sell.Quantity)

	totalProfit := buy.TotalProfit + sellTotalProfit

	err = s.buyDao.UpdateBuyWhenSell(buy.ID, &models.BuyUpdateProfit{
		Inventory:   inventory,
		TotalProfit: totalProfit,
	})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	sell.Profit = sell.Price - buy.Price
	sell.TotalProfit = sellTotalProfit
	result, err := s.dao.Create(sell)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return result, nil
}

func (s *sellService) Exists(id uint64) (bool, error) {
	return s.dao.Exists(id)
}

func (s *sellService) Update(id uint64, sell *models.SellUpdate) error {
	if sell.Price == 0 {
		return errors.New("请填写卖出价")
	}
	if sell.Quantity == 0 {
		return errors.New("请填写卖出数量")
	}
	db := s.dao.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Get the corresponding buy record
	lastSell, err := s.dao.GetItem(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	buy, err := s.buyDao.GetItem(lastSell.BuyID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if buy.Inventory < sell.Quantity {
		tx.Rollback()
		return errors.New("卖出数量超过买入数量")
	}
	// Update the stock of the corresponding buy record
	inventory := buy.Inventory - lastSell.Quantity - sell.Quantity
	totalProfit := buy.TotalProfit - lastSell.Price*float64(lastSell.Quantity) + sell.Price*float64(sell.Quantity)

	err = s.buyDao.UpdateBuyWhenSell(buy.ID, &models.BuyUpdateProfit{
		Inventory:   inventory,
		TotalProfit: totalProfit,
	})
	if err != nil {
		tx.Rollback()
		return err
	}
	sell.Profit = sell.Price - buy.Price
	sell.TotalProfit = (sell.Price - buy.Price) * float64(sell.Quantity)
	return s.dao.Update(id, sell)
}

func (s *sellService) Delete(id uint64) error {
	return s.dao.Delete(id)
}
