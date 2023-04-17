package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type GoodsService interface {
	GetList(query *models.GoodsListQuery) utils.ApiResponse
	GetItem(id uint64) utils.ApiResponse
	Create(goods *models.Goods) utils.ApiResponse
	Exists(id uint64) (bool, error)
	Update(id uint64, goods *models.GoodsUpdate) utils.ApiResponse
	Delete(id uint64) utils.ApiResponse
}

type goodsService struct {
	dao *dao.GoodsDAO
}

func NewGoodsService(dao *dao.GoodsDAO) *goodsService {
	return &goodsService{dao}
}

func (s *goodsService) GetList(query *models.GoodsListQuery) utils.ApiResponse {
	res, err := s.dao.GetList(query)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}

func (s *goodsService) GetItem(id uint64) utils.ApiResponse {
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

func (s *goodsService) Create(goods *models.Goods) utils.ApiResponse {
	exists, err := s.dao.ExistsByName(goods.Name)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if exists {
		return utils.ApiErrorResponse(-1, "商品名称已存在")
	}
	goods, err = s.dao.Create(goods)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(goods)
}

func (s *goodsService) Exists(id uint64) (bool, error) {
	return s.dao.Exists(id)
}

func (s *goodsService) Update(id uint64, goods *models.GoodsUpdate) utils.ApiResponse {
	exists, err := s.Exists(id)

	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !exists {
		return utils.ApiErrorResponse(404, "记录不存在")
	}
	err = s.dao.Update(id, goods)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}

func (s *goodsService) Delete(id uint64) utils.ApiResponse {
	err := s.dao.Delete(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}
