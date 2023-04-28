package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type UserService interface {
	SignUp(sell *models.Sell) utils.ApiResponse
	Login(sell *models.Sell) utils.ApiResponse
}

type userService struct {
	dao *dao.UserDAO
}

func NewUserService(dao *dao.UserDAO) *userService {
	return &userService{dao}
}

func (s *userService) SignUp(user *models.User) utils.ApiResponse {
	salt, err := utils.GenerateSalt()
	password, err := utils.GenerateHashedPassword(user.Password, salt)
	user.Password = password
	user.Salt = salt
	id, err := s.dao.SignUp(user)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(id)
}

func (s *userService) Login(user *models.User) utils.ApiResponse {
	find, err := s.dao.GetUser(user.Name)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	password, err := utils.GenerateHashedPassword(user.Password, find.Salt)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	token, err := s.dao.Login(&models.User{
		Name:     user.Name,
		Password: password,
	})
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}

	return utils.ApiSuccessResponse(token)
}
