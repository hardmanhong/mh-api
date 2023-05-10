package dao

import (
	"errors"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type AccountDAO struct {
	db *gorm.DB
}

func NewAccountDAO(db *gorm.DB) *AccountDAO {
	return &AccountDAO{db}
}

func (dao *AccountDAO) GetList(userId uint64, query *models.AccountListQuery) (*models.ListResponse, error) {
	response := models.ListResponse{
		List: make([]interface{}, 0),
	}
	db := dao.db
	if query.Name != "" {
		db = db.Where("user_name LIKE ?", "%"+query.Name+"%")
	}

	var accountList []*models.Account
	err := db.Model(&models.Account{}).Where("user_id = ?", userId).Find(&accountList).Error
	if err != nil {
		return nil, err
	}

	for _, g := range accountList {
		response.List = append(response.List, g)
	}

	return &response, nil
}
func (dao *AccountDAO) GetItem(id uint32) (*models.Account, error) {
	var account models.Account
	err := dao.db.Where("id = ?", id).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &account, nil
}
func (dao *AccountDAO) Create(account *models.Account) (*models.Account, error) {
	err := dao.db.Create(account).Error
	if err != nil {
		return nil, err
	}
	return account, err
}
func (dao *AccountDAO) ExistsByName(name string) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Account{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *AccountDAO) Exists(id uint32) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Account{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *AccountDAO) Update(id uint32, account *models.AccountUpdate) error {
	return dao.db.Model(&models.Account{}).Where("id = ?", id).Updates(&models.Account{
		Name:   account.Name,
		Server: account.Server,
	}).Error
}

func (dao *AccountDAO) Delete(id uint32) error {
	return dao.db.Where("id = ?", id).Delete(&models.Account{}).Error
}
