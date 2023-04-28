package dao

import (
	"errors"
	"fmt"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{db}
}
func (dao *UserDAO) GetDB() *gorm.DB {
	return dao.db
}
func (dao *UserDAO) SignUp(user *models.User) (int, error) {
	err := dao.db.Create(user).Error
	if err != nil {
		return 0, err
	}
	return int(user.ID), nil
}
func (dao *UserDAO) Login(user *models.User) (bool, error) {
	var count int64
	err := dao.db.Model(&models.User{}).Where("name = ? and password = ?", user.Name, user.Password).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *UserDAO) GetUser(name string) (*models.User, error) {
	user := &models.User{}
	err := dao.db.Where("name = ?", name).Preload("Sales").First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user name=%s not found", name)
	}
	return user, err
}
