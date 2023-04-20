package dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type BuyDAO struct {
	db *gorm.DB
}

func NewBuyDAO(db *gorm.DB) *BuyDAO {
	return &BuyDAO{db}
}
func (dao *BuyDAO) GetDB() *gorm.DB {
	return dao.db
}

// 计算时间范围
func CalculateDateRange(t time.Time, dType string) (time.Time, time.Time) {
	switch dType {
	case "day":
		return t.AddDate(0, -1, 0), t
	case "week":
		return t.AddDate(0, 0, -30), t
	case "month":
		return t.AddDate(0, -11, 0), t
	case "year":
		return t.AddDate(-3, 0, 0), t
	default:
		return t.AddDate(0, -1, 0), t
	}
}

// 按时间维度分组
func GroupByTimeDimension(dType string) string {
	switch dType {
	case "day":
		return "DATE(created_at)"
	case "week":
		return "YEARWEEK(created_at)"
	case "month":
		return "DATE(created_at)"
	case "year":
		return "DATE(created_at)"
	default:
		return "DATE(created_at)"
	}
}
func (dao *BuyDAO) GetProfit(dType string) ([]models.BuyProfit, error) {
	if dType == "day" {
		return dao.GetProfitByDay()
	} else if dType == "week" {
		return dao.GetProfitByWeek()
	} else if dType == "month" {
		return dao.GetProfitByMonth()
	} else if dType == "year" {
		return dao.GetProfitByYear()
	}
	return []models.BuyProfit{}, nil

}
func (dao *BuyDAO) GetProfitByDay() ([]models.BuyProfit, error) {
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := CalculateDateRange(now, "day")
	// 分组字段
	groupBy := GroupByTimeDimension("day")

	rows, err := dao.db.Table("buy").
		Select(groupBy+" AS date, sum(total_profit) AS profit").
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.BuyProfit
	for rows.Next() {
		var day time.Time
		var value float64
		if err := rows.Scan(&day, &value); err != nil {
			return nil, err
		}
		// // 将日期转换为指定格式的字符串
		label := day.Format("2006-01-02")
		list = append(list, models.BuyProfit{Value: value, Label: label})
	}

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.BuyProfit{}, nil
	}
	return list, nil
}
func (dao *BuyDAO) GetProfitByWeek() ([]models.BuyProfit, error) {
	rows, err := dao.db.Raw(`
		SELECT monday, SUM(total_profit) FROM (
			SELECT *,
				WEEKDAY(created_at) AS weekday,
				DATE_FORMAT(DATE_ADD(created_at, INTERVAL - (WEEKDAY(created_at) + 1) DAY),'%Y-%m-%d') AS monday
			FROM buy
		) AS x GROUP BY monday ORDER BY monday
	`).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.BuyProfit
	for rows.Next() {
		var label string
		var value float64
		if err := rows.Scan(&label, &value); err != nil {
			return nil, err
		}
		list = append(list, models.BuyProfit{Value: value, Label: label})
	}

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.BuyProfit{}, nil
	}
	return list, nil
}
func (dao *BuyDAO) GetProfitByMonth() ([]models.BuyProfit, error) {
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := CalculateDateRange(now, "month")
	// 分组字段
	// groupBy := GroupByTimeDimension("month")

	rows, err := dao.db.Table("buy").
		Select("DATE_FORMAT(created_at, '%Y-%m') AS month, SUM(total_profit) AS total_profit").
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("month").
		Order("month").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.BuyProfit
	for rows.Next() {
		var label string
		var value float64
		if err := rows.Scan(&label, &value); err != nil {
			return nil, err
		}
		// 2006-01-02 15:04:05 奇葩的格式
		list = append(list, models.BuyProfit{Value: value, Label: label})
	}

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.BuyProfit{}, nil
	}
	return list, nil
}
func (dao *BuyDAO) GetProfitByYear() ([]models.BuyProfit, error) {
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := CalculateDateRange(now, "year")
	// 分组字段

	rows, err := dao.db.Table("buy").
		Select("DATE_FORMAT(created_at, '%Y') AS year, SUM(total_profit) AS total_profit").
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("year").
		Order("year").
		Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.BuyProfit
	for rows.Next() {
		var label string
		var value float64
		if err := rows.Scan(&label, &value); err != nil {
			return nil, err
		}
		// 2006-01-02 15:04:05 奇葩的格式
		list = append(list, models.BuyProfit{Value: value, Label: label})
	}

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.BuyProfit{}, nil
	}
	return list, nil
}
func (dao *BuyDAO) GetTotalProfit() (float64, error) {
	var totalProfit float64
	err := dao.db.Model(&models.Buy{}).Select("SUM(total_profit)").Pluck("SUM(total_profit)", &totalProfit).Error
	if err != nil {
		return 0, err
	}
	return totalProfit, nil
}
func (dao *BuyDAO) GetList(query *models.BuyListQuery) (*models.BuyListResponse, error) {
	response := models.BuyListResponse{
		TotalProfit: 0,
		PaginationResponse: models.PaginationResponse{
			Total: 0,
			List:  make([]interface{}, 0),
		},
	}
	var total int64
	var buyList []models.Buy
	tx := dao.db.Model(&models.Buy{}).Preload("Goods").Preload("Sales")
	if query.CreatedAtFrom != nil {
		tx = tx.Where("created_at >= ?", query.CreatedAtFrom)
	}
	if query.CreatedAtTo != nil {
		tx = tx.Where("created_at <= ?", query.CreatedAtTo)
	}
	jsonData, _ := json.Marshal(query)
	println("query", string(jsonData))
	if len(query.GoodsIDs) > 0 {
		tx = tx.Where("goods_id IN (?)", query.GoodsIDs)
	}

	err := tx.Count(&total).Error
	if err != nil {
		return nil, err
	}
	offset := (query.Page - 1) * query.PageSize
	inventorySorter := query.InventorySorter
	if inventorySorter == "asc" {
		tx = tx.Order("inventory asc")
	} else if inventorySorter == "desc" {
		tx = tx.Order("inventory desc")
	}
	hasSoldSorter := query.HasSoldSorter
	if hasSoldSorter == "asc" {
		tx = tx.Order("has_sold asc")
	} else if hasSoldSorter == "desc" {
		tx = tx.Order("has_sold desc")
	}

	err = tx.Order("created_at desc").Offset(offset).Limit(query.PageSize).Find(&buyList).Error
	if err != nil {
		return nil, err
	}
	response.Total = total
	for _, g := range buyList {
		response.TotalAmount += g.TotalAmount
		response.TotalProfit += g.TotalProfit
		response.List = append(response.List, g)
	}
	response.TotalAmount, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", response.TotalAmount), 64)
	response.TotalProfit, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", response.TotalProfit), 64)
	return &response, nil
}

func (dao *BuyDAO) GetItem(id uint64) (*models.Buy, error) {
	buy := &models.Buy{}
	err := dao.db.Where("id = ?", id).Preload("Sales").First(buy).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("buy id=%d not found", id)
	}
	return buy, err
}

func (dao *BuyDAO) Create(buy *models.Buy) (*models.Buy, error) {
	buy.Inventory = buy.Quantity
	err := dao.db.Create(buy).Error
	if err != nil {
		return nil, err
	}

	// 更新关联的 Goods 信息
	dao.db.Model(&buy).Association("Goods").Append(&models.Goods{ID: buy.GoodsID})

	return buy, nil
}
func (dao *BuyDAO) Exists(id uint64) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Buy{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *BuyDAO) Update(id uint64, buy *models.BuyUpdate) error {
	return dao.db.Table("buy").Where("id = ?", id).Updates(buy).Error
}

func (dao *BuyDAO) Delete(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&models.Buy{}).Error
}

func (dao *BuyDAO) UpdateBuyWhenSell(id uint64, buy *models.BuyUpdateProfit) error {
	return dao.db.Model(&models.Buy{}).Where("id = ?", id).Update("has_sold", 1).Update("inventory", buy.Inventory).Update("total_profit", buy.TotalProfit).Error
}
