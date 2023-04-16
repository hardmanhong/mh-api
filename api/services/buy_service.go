package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
)

type BuyService interface {
	GetList(query *models.BuyListQuery) (*models.BuyListResponse, error)
	GetItem(id uint64) (*models.Buy, error)
	Create(buy *models.Buy) (*models.Buy, error)
	Exists(id uint64) (bool, error)
	Update(id uint64, buy *models.BuyUpdate) error
	Delete(id uint64) error
}

type buyService struct {
	dao *dao.BuyDAO
}

func NewBuyService(dao *dao.BuyDAO) *buyService {
	return &buyService{dao}
}

func (s *buyService) GetList(query *models.BuyListQuery) (*models.BuyListResponse, error) {
	resp, err := s.dao.GetList(query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *buyService) GetItem(id uint64) (*models.Buy, error) {
	buy, err := s.dao.GetItem(id)
	if err != nil {
		return nil, err
	}
	return buy, nil
}

func (s *buyService) Create(buy *models.Buy) (*models.Buy, error) {
	buy.TotalAmount = buy.Price * float64(buy.Quantity)
	buy, err := s.dao.Create(buy)
	if err != nil {
		return nil, err
	}
	return buy, nil
}
func (s *buyService) Exists(id uint64) (bool, error) {
	return s.dao.Exists(id)
}

func (s *buyService) Update(id uint64, buy *models.BuyUpdate) error {
	return s.dao.Update(id, buy)
}

func (s *buyService) Delete(id uint64) error {
	return s.dao.Delete(id)
}
func (s *buyService) UpdateInventory(buy *models.Buy, inventory int) error {
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
		return err
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
