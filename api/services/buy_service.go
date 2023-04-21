package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type BuyService interface {
	GetList(query *models.BuyListQuery) utils.ApiResponse
	GetItem(id uint64) utils.ApiResponse
	Create(buy *models.Buy) utils.ApiResponse
	Exists(id uint64) (bool, error)
	Update(id uint64, buy *models.BuyUpdate) utils.ApiResponse
	Delete(id uint64) utils.ApiResponse
}

type buyService struct {
	dao *dao.BuyDAO
}

func NewBuyService(dao *dao.BuyDAO) *buyService {
	return &buyService{dao}
}

func (s *buyService) GetList(query *models.BuyListQuery) utils.ApiResponse {
	res, err := s.dao.GetList(query)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}

func (s *buyService) GetItem(id uint64) utils.ApiResponse {
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

func (s *buyService) Create(buy *models.Buy) utils.ApiResponse {
	buy.Inventory = buy.Quantity
	buy.TotalAmount = buy.Price * float64(buy.Quantity)
	buy, err := s.dao.Create(buy)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(buy)
}
func (s *buyService) Exists(id uint64) (bool, error) {
	return s.dao.Exists(id)
}

func (s *buyService) Update(id uint64, buy *models.BuyUpdate) utils.ApiResponse {
	exists, err := s.Exists(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !exists {
		return utils.ApiErrorResponse(404, "记录不存在")
	}
	buyItem, err := s.dao.GetItem(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if buyItem.HasSold == 0 {
		buy.Inventory = buy.Quantity
		buy.TotalAmount = buy.Price * float64(buy.Quantity)
	} else {
		// 已有卖出记录时
		if buy.Quantity < buyItem.Inventory {
			return utils.ApiErrorResponse(-1, "已有卖出记录，买入数量不能小于库存")
		}
		// 卖出数量 = 上次买入 - 上次库存
		soldQuantity := buyItem.Quantity - buyItem.Inventory
		// 当前库存 = 当前数量 - 卖出数量
		buy.Inventory = buy.Quantity - soldQuantity
		// 当前总计 = 上次总计 + 新增的金额
		buy.TotalAmount = buyItem.TotalAmount + float64(buy.Quantity-buyItem.Quantity)*buy.Price
	}
	err = s.dao.Update(id, buy)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}

func (s *buyService) Delete(id uint64) utils.ApiResponse {
	err := s.dao.Delete(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}
func (s *buyService) UpdateInventory(buy *models.Buy, inventory int) utils.ApiResponse {
	db := s.dao.GetDB()
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	buy.Inventory = inventory
	err := db.Save(buy).Error
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return utils.ApiErrorResponse(-1, err.Error())
	}

	return utils.ApiSuccessResponse(nil)
}
