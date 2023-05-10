package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type PetService interface {
	GetItem(id uint32) utils.ApiResponse
	Create(pet *models.Pet) utils.ApiResponse
	Exists(id uint32) (bool, error)
	Update(id uint32, pet *models.Pet) utils.ApiResponse
	Delete(id uint32) utils.ApiResponse
}

type petService struct {
	dao *dao.PetDAO
}

func NewPetService(dao *dao.PetDAO) *petService {
	return &petService{dao}
}

func (s *petService) GetItem(id uint32) utils.ApiResponse {
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

func (s *petService) Create(pet *models.Pet) utils.ApiResponse {
	db := s.dao.GetDB()
	tx := db.Begin()
	pet, err := s.dao.Create(pet, tx)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(pet)
}

func (s *petService) Exists(id uint32) (bool, error) {
	return s.dao.Exists(id)
}

func (s *petService) Update(id uint32, pet *models.Pet) utils.ApiResponse {
	exists, err := s.Exists(id)

	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !exists {
		return utils.ApiErrorResponse(404, "记录不存在")
	}
	err = s.dao.Update(id, pet)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}

func (s *petService) Delete(id uint32) utils.ApiResponse {
	err := s.dao.Delete(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}
