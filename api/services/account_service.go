package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type AccountService interface {
	GetList(userId uint64, query *models.AccountListQuery) utils.ApiResponse
	GetItem(id uint32) utils.ApiResponse
	Create(userId uint64, account *models.Account) utils.ApiResponse
	Exists(id uint32) (bool, error)
	Update(id uint32, account *models.AccountUpdate) utils.ApiResponse
	Delete(id uint32) utils.ApiResponse
}

type accountService struct {
	dao *dao.AccountDAO
}

func NewAccountService(dao *dao.AccountDAO) *accountService {
	return &accountService{dao}
}

func (s *accountService) GetList(userId uint64, query *models.AccountListQuery) utils.ApiResponse {
	res, err := s.dao.GetList(userId, query)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}

func (s *accountService) GetItem(id uint32) utils.ApiResponse {
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

func (s *accountService) Create(userId uint64, account *models.Account) utils.ApiResponse {
	account.UserId = userId
	exists, err := s.dao.ExistsByName(account.Name)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if exists {
		return utils.ApiErrorResponse(-1, "账号已存在")
	}
	account, err = s.dao.Create(account)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(account)
}

func (s *accountService) Exists(id uint32) (bool, error) {
	return s.dao.Exists(id)
}

func (s *accountService) Update(id uint32, account *models.AccountUpdate) utils.ApiResponse {
	exists, err := s.Exists(id)

	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !exists {
		return utils.ApiErrorResponse(404, "记录不存在")
	}
	err = s.dao.Update(id, account)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}

func (s *accountService) Delete(id uint32) utils.ApiResponse {
	err := s.dao.Delete(id)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(nil)
}
