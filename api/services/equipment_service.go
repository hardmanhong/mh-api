package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type EquipmentService interface {
	GetItem(id uint32) utils.ApiResponse
	Create(equipment *models.Equipment) utils.ApiResponse
	Exists(id uint32) (bool, error)
	Update(id uint32, equipment *models.Equipment) utils.ApiResponse
	Delete(id uint32) utils.ApiResponse
}

type equipmentService struct {
	dao *dao.EquipmentDAO
}

func NewEquipmentService(dao *dao.EquipmentDAO) *equipmentService {
	return &equipmentService{dao}
}

func (s *equipmentService) GetItem(id uint32) utils.ApiResponse {
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

func (s *equipmentService) Create(equipment *models.Equipment) utils.ApiResponse {
	db := s.dao.GetDB()
	tx := db.Begin()
	equipment, err := s.dao.Create(equipment, tx)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(equipment)
}

func (s *equipmentService) Exists(id uint32) (bool, error) {
	return s.dao.Exists(id)
}

func (s *equipmentService) Update(id uint32, equipment *models.Equipment) utils.ApiResponse {
	exists, err := s.Exists(id)

	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !exists {
		return utils.ApiErrorResponse(404, "记录不存在")
	}
	err = s.dao.Update(id, equipment)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}

func (s *equipmentService) Delete(id uint32) utils.ApiResponse {
	err := s.dao.Delete(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}
