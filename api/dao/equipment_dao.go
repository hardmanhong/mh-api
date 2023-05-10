package dao

import (
	"errors"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type EquipmentDAO struct {
	db *gorm.DB
}

func (dao *EquipmentDAO) GetDB() *gorm.DB {
	return dao.db
}

func NewEquipmentDAO(db *gorm.DB) *EquipmentDAO {
	return &EquipmentDAO{db}
}

func (dao *EquipmentDAO) GetItem(id uint32) (*models.Equipment, error) {
	var equipment models.Equipment
	err := dao.db.Where("id = ?", id).First(&equipment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &equipment, nil
}
func (dao *EquipmentDAO) Create(equipment *models.Equipment, tx *gorm.DB) (*models.Equipment, error) {
	err := tx.Create(equipment).Error
	if err != nil {
		return nil, err
	}
	return equipment, err
}
func (dao *EquipmentDAO) Exists(id uint32) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Equipment{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *EquipmentDAO) Update(id uint32, equipment *models.Equipment) error {
	return dao.db.Model(&models.Equipment{}).Where("id = ?", id).Updates(equipment).Error
}

func (dao *EquipmentDAO) Delete(id uint32) error {
	return dao.db.Where("id = ?", id).Delete(&models.Equipment{}).Error

}
