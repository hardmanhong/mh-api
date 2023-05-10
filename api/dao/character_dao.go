package dao

import (
	"errors"
	"fmt"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type CharacterDAO struct {
	db *gorm.DB
}

func NewCharacterDAO(db *gorm.DB) *CharacterDAO {
	return &CharacterDAO{db}
}
func (dao *CharacterDAO) GetDB() *gorm.DB {
	return dao.db
}

func (dao *CharacterDAO) GetList() (*models.ListResponse, error) {
	response := models.ListResponse{
		List: make([]interface{}, 0),
	}
	db := dao.db

	var list []*models.Character
	err := db.Model(&models.Character{}).Preload("Account").Preload("Equipment").Preload("Pet").Find(&list).Error
	if err != nil {
		return nil, err
	}

	for _, g := range list {
		response.List = append(response.List, g)
	}

	return &response, nil
}
func (dao *CharacterDAO) GetItem(id uint32, tx *gorm.DB) (*models.Character, error) {
	sell := &models.Character{}
	err := tx.Where("id = ?", id).Preload("Account").Preload("Equipment").Preload("Pet").First(sell).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("找不到该记录")
	}
	return sell, err
}

func (dao *CharacterDAO) Create(sell *models.Character, tx *gorm.DB) (*models.Character, error) {
	return sell, tx.Create(&sell).Error
}

func (dao *CharacterDAO) Exists(id uint32) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Character{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *CharacterDAO) ExistsByName(name string) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Character{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *CharacterDAO) Update(id uint32, sell *models.Character) error {
	return dao.db.Model(&models.Character{}).Where("id = ?", id).Updates(sell).Error
}

func (dao *CharacterDAO) Delete(id uint32) error {
	return dao.db.Where("id = ?", id).Delete(&models.Character{}).Error
}
