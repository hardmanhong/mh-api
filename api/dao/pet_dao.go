package dao

import (
	"errors"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type PetDAO struct {
	db *gorm.DB
}

func NewPetDAO(db *gorm.DB) *PetDAO {
	return &PetDAO{db}
}

func (dao *PetDAO) GetDB() *gorm.DB {
	return dao.db
}

func (dao *PetDAO) GetItem(id uint32) (*models.Pet, error) {
	var pet models.Pet
	err := dao.db.Where("id = ?", id).First(&pet).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &pet, nil
}
func (dao *PetDAO) Create(pet *models.Pet, tx *gorm.DB) (*models.Pet, error) {
	err := tx.Create(pet).Error
	if err != nil {
		return nil, err
	}
	return pet, err
}
func (dao *PetDAO) Exists(id uint32) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Pet{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *PetDAO) Update(id uint32, pet *models.Pet) error {
	return dao.db.Model(&models.Pet{}).Where("id = ?", id).Updates(pet).Error
}

func (dao *PetDAO) Delete(id uint32) error {
	return dao.db.Where("id = ?", id).Delete(&models.Pet{}).Error
}
