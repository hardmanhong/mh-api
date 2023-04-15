package dao

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type BuyDAO struct {
	db *gorm.DB
}

func NewBuyDAO(db *gorm.DB) *BuyDAO {
	return &BuyDAO{db}
}
func (dao *BuyDAO) GetDB() *gorm.DB {
	return dao.db
}
func (dao *BuyDAO) GetList(query *models.BuyListQuery) (*models.BuyListResponse, error) {
	response := models.BuyListResponse{
		TotalProfit: 0,
		PaginationResponse: models.PaginationResponse{
			Total: 0,
			List:  make([]interface{}, 0),
		},
	}
	var total int64
	var buyList []models.Buy
	tx := dao.db.Model(&models.Buy{}).Preload("Goods").Preload("Sales")
	if query.CreatedAtFrom != nil {
		tx = tx.Where("created_at >= ?", query.CreatedAtFrom)
	}
	if query.CreatedAtTo != nil {
		tx = tx.Where("created_at <= ?", query.CreatedAtTo)
	}
	jsonData, _ := json.Marshal(query)
	println("query", string(jsonData))
	if len(query.GoodsIDs) > 0 {
		tx = tx.Where("goods_id IN (?)", query.GoodsIDs)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return nil, err
	}
	offset := (query.Page - 1) * query.PageSize
	err = tx.Order("created_at asc").Offset(offset).Limit(query.PageSize).Find(&buyList).Error
	if err != nil {
		return nil, err
	}
	response.Total = total
	for _, g := range buyList {
		response.TotalProfit += g.TotalProfit
		response.List = append(response.List, g)
	}
	return &response, nil
}

func (dao *BuyDAO) GetItem(id uint64) (*models.Buy, error) {
	buy := &models.Buy{}
	err := dao.db.Where("id = ?", id).Preload("Sales").First(buy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("buy id=%d not found", id)
	}
	return buy, err
}

func (dao *BuyDAO) Create(buy *models.Buy) (*models.Buy, error) {
	buy.Inventory = buy.Quantity
	err := dao.db.Create(buy).Error
	if err != nil {
		return nil, err
	}

	// 更新关联的 Goods 信息
	dao.db.Model(&buy).Association("Goods").Append(&models.Goods{ID: buy.GoodsID})

	return buy, nil
}
func (dao *BuyDAO) Exists(id uint64) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Buy{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *BuyDAO) Update(id uint64, buy *models.BuyUpdate) error {
	return dao.db.Table("buy").Where("id = ?", id).Updates(buy).Error
}

func (dao *BuyDAO) Delete(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&models.Buy{}).Error
}

func (dao *BuyDAO) UpdateBuyWhenSell(id uint64, buy *models.BuyUpdateProfit) error {
	return dao.db.Model(&models.Buy{}).Where("id = ?", id).Update("inventory", buy.Inventory).Update("total_profit", buy.TotalProfit).Error
}
