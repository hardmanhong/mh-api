package services

import (
	"github.com/hardmanhong/api/dao"
	"github.com/hardmanhong/api/utils"
)

type StatisticsService interface {
	GetStatistics(dType string) utils.ApiResponse
	GetTotalProfit() utils.ApiResponse
}

type statisticsService struct {
	dao *dao.StatisticsDAO
}

func NewStatisticsService(dao *dao.StatisticsDAO) *statisticsService {
	return &statisticsService{dao}
}

func (s *statisticsService) GetStatistics(dType string) utils.ApiResponse {
	res, err := s.dao.GetStatistics(dType)
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}

func (s *statisticsService) GetTotalProfit() utils.ApiResponse {
	res, err := s.dao.GetTotalProfit()
	if err != nil {
		return utils.ApiErrorResponse(-1, err.Error())
	}
	return utils.ApiSuccessResponse(res)
}
