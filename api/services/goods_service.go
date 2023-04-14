package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
)

type GoodsService interface {
	GetList(query *models.GoodsListQuery) (*models.PaginationResponse, error)
	GetItem(id uint64) (*models.Goods, error)
	Create(goods *models.Goods) error
	Exists(id uint64) (bool, error)
	Update(id uint64, goods *models.GoodsUpdate) error
	Delete(id uint64) error
}

type goodsService struct {
	dao *dao.GoodsDAO
}

func NewGoodsService(dao *dao.GoodsDAO) *goodsService {
	return &goodsService{dao}
}

func (s *goodsService) GetList(query *models.GoodsListQuery) (*models.PaginationResponse, error) {
	resp, err := s.dao.GetList(query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *goodsService) GetItem(id uint64) (*models.Goods, error) {
	goods, err := s.dao.GetItem(id)
	if err != nil {
		return nil, err
	}
	return goods, nil
}

func (s *goodsService) Create(goods *models.Goods) error {
	return s.dao.Create(goods)
}

func (s *goodsService) Exists(id uint64) (bool, error) {
	return s.dao.Exists(id)
}

func (s *goodsService) Update(id uint64, goods *models.GoodsUpdate) error {
	return s.dao.Update(id, goods)
}

func (s *goodsService) Delete(id uint64) error {
	return s.dao.Delete(id)
}
