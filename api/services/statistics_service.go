package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/utils"
)

type StatisticsService interface {
	GetStatistics(userId uint64, dType string) utils.ApiResponse
	GetTotalProfit(userId uint64) utils.ApiResponse
}

type statisticsService struct {
	dao *dao.StatisticsDAO
}

func NewStatisticsService(dao *dao.StatisticsDAO) *statisticsService {
	return &statisticsService{dao}
}

func (s *statisticsService) GetStatistics(userId uint64, dType string) utils.ApiResponse {
	res, err := s.dao.GetStatistics(userId, dType)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}

func (s *statisticsService) GetTotalProfit(userId uint64) utils.ApiResponse {
	res, err := s.dao.GetTotalProfit(userId)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}
