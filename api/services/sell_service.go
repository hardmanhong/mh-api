package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type SellService interface {
	GetItem(id uint64) utils.ApiResponse
	Create(sell *models.Sell) utils.ApiResponse
	Exists(id uint64) (bool, error)
	Update(id uint64, sell *models.SellUpdate) utils.ApiResponse
	Delete(id uint64) utils.ApiResponse
}

type sellService struct {
	dao    *dao.SellDAO
	buyDao *dao.BuyDAO
}

func NewSellService(dao *dao.SellDAO, buyDao *dao.BuyDAO) *sellService {
	return &sellService{dao, buyDao}
}
func (s *sellService) GetItem(id uint64) utils.ApiResponse {
	item, err := s.dao.GetItem(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, "Failed to get item")
	}
	// 返回结果
	if item == nil {
		return utils.ApiErrorResponse(-1, "Item not found")
	}
	return utils.ApiSuccessResponse(&item)
}

func (s *sellService) Create(sell *models.Sell) utils.ApiResponse {
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
		return utils.ApiErrorResponse(-1, err.Error())
	}
	// Check if there is enough stock to sell
	if sell.Price == 0 {
		return utils.ApiErrorResponse(-1, "请填写卖出价")
	}
	if sell.Quantity == 0 {
		return utils.ApiErrorResponse(-1, "请填写卖出数量")
	}
	if buy.Inventory < sell.Quantity {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, "卖出数量超过买入数量")
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
		return utils.ApiErrorResponse(-1, err.Error())
	}
	sell.Profit = sell.Price - buy.Price
	sell.TotalProfit = sellTotalProfit
	result, err := s.dao.Create(sell)
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	// Commit the transaction
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(&result)
}

func (s *sellService) Exists(id uint64) (bool, error) {
	return s.dao.Exists(id)
}

func (s *sellService) Update(id uint64, sell *models.SellUpdate) utils.ApiResponse {
	exists, err := s.Exists(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !exists {
		return utils.ApiErrorResponse(404, "记录不存在")
	}

	if sell.Price == 0 {
		return utils.ApiErrorResponse(-1, "请填写卖出价")
	}
	if sell.Quantity == 0 {
		return utils.ApiErrorResponse(-1, "请填写卖出数量")
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
		return utils.ApiErrorResponse(-1, err.Error())
	}
	buy, err := s.buyDao.GetItem(lastSell.BuyID)
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if buy.Inventory < sell.Quantity {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, "卖出数量超过买入数量")
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
		return utils.ApiErrorResponse(-1, err.Error())
	}
	sell.Profit = sell.Price - buy.Price
	sell.TotalProfit = (sell.Price - buy.Price) * float64(sell.Quantity)
	err = s.dao.Update(id, sell)
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}

func (s *sellService) Delete(id uint64) utils.ApiResponse {
	err := s.dao.Delete(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}
