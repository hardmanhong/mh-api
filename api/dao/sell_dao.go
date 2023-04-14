package dao

import (
	"errors"
	"fmt"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type SellDAO struct {
	db *gorm.DB
}

func NewSellDAO(db *gorm.DB) *SellDAO {
	return &SellDAO{db}
}
func (dao *SellDAO) GetDB() *gorm.DB {
	return dao.db
}
func (dao *SellDAO) GetItem(id uint64) (*models.Sell, error) {
	sell := &models.Sell{}
	err := dao.db.Where("id = ?", id).First(sell).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("sell id=%d not found", id)
	}
	return sell, err
}

func (dao *SellDAO) Create(sell *models.Sell) (*models.Sell, error) {
	return sell, dao.db.Create(&sell).Error
}

func (dao *SellDAO) Exists(id uint64) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Sell{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (dao *SellDAO) Update(id uint64, sell *models.SellUpdate) error {
	return dao.db.Table("sell").Where("id = ?", id).Updates(sell).Error
}

func (dao *SellDAO) Delete(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&models.Sell{}).Error
}
