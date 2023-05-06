package services

import (
	"fmt"
	"time"

	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/models"
	"github.com/hardmanhong/api/utils"
)

type UserService interface {
	SignUp(sell *models.User) utils.ApiResponse
	Login(sell *models.User) utils.ApiResponse
}

type userService struct {
	dao *dao.UserDAO
}

func NewUserService(dao *dao.UserDAO) *userService {
	return &userService{dao}
}

func (s *userService) SignUp(user *models.User) utils.ApiResponse {
	isExist, err := s.dao.ExistName(user.Name)
	fmt.Println("SignUp", user.Name, user.Password, user.Salt)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if isExist {
		return utils.ApiErrorResponse(-1, "用户名已存在")
	}
	salt, _ := utils.GenerateString(16)
	password, _ := utils.GenerateHashedPassword(user.Password, salt)
	user.Password = password
	user.Salt = salt
	id, err := s.dao.SignUp(user)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(id)
}

func (s *userService) Login(user *models.User) utils.ApiResponse {
	isExist, err := s.dao.ExistName(user.Name)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !isExist {
		return utils.ApiErrorResponse(-1, "用户名不存在")
	}
	find, err := s.dao.GetUser(user.Name)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	password, _ := utils.GenerateHashedPassword(user.Password, find.Salt)
	isOk, err := s.dao.Login(&models.User{
		Name:     user.Name,
		Password: password,
	})
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	if !isOk {
		return utils.ApiErrorResponse(-1, "用户名或密码错误")
	}
	token, err := utils.GenerateString(32)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	expireAt := time.Now().AddDate(0, 1, 0)

	err = s.dao.CreateUserToken(&models.Token{
		UserId:   find.ID,
		Token:    token,
		ExpireAt: expireAt,
	})
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(token)
}
