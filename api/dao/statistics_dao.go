package dao

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type StatisticsDAO struct {
	db *gorm.DB
}

func NewStatisticsDAO(db *gorm.DB) *StatisticsDAO {
	return &StatisticsDAO{db}
}
func (dao *StatisticsDAO) GetDB() *gorm.DB {
	return dao.db
}

// 计算时间范围
func calculateDateRange(t time.Time, dType string) (time.Time, time.Time) {
	start := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())

	switch dType {
	case "day":
		return start.AddDate(0, -1, 0), end
	case "week":
		return start.AddDate(0, 0, -30), end
	case "month":
		return start.AddDate(0, -5, 0), end
	case "year":
		return start.AddDate(-2, 0, 0), end
	default:
		return start.AddDate(0, -1, 0), end
	}
}

// 按时间维度分组
func groupByTimeDimension(dType string) string {
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
func (dao *StatisticsDAO) GetStatistics(userId uint64, dType string) ([]models.Statistics, error) {
	if dType == "day" {
		return dao.getProfitByDay(userId)
	} else if dType == "week" {
		return dao.getProfitByWeek(userId)
	} else if dType == "month" {
		return dao.getProfitByMonth(userId)
	} else if dType == "year" {
		return dao.getProfitByYear(userId)
	}
	return []models.Statistics{}, nil

}
func generateDateSequence(start, end time.Time) []time.Time {
	start = start.UTC()
	end = end.UTC()
	var dates []time.Time
	for date := start; !date.After(end); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date)
	}
	return dates
}

type BuyProfit struct {
	Date   string
	Amount float64
}

func (dao *StatisticsDAO) getProfitByDay(userId uint64) ([]models.Statistics, error) {
	var list []models.Statistics
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := calculateDateRange(now, "day")
	// 分组字段
	// groupBy := groupByTimeDimension("day")

	// 构建日期范围
	dates := make([]string, 0)
	for t := start; !t.After(end); t = t.AddDate(0, 0, 1) {
		dates = append(dates, t.Format("2006-01-02"))
	}

	buyRows, err := dao.db.Table("buy").
		Select("DATE_FORMAT(created_at, '%Y-%m-%d') AS date, SUM(total_amount) AS amount").
		Where("user_id = ?", userId).
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer buyRows.Close()

	// 将结果转换为以日期为 key 的 map
	buyMap := make(map[string]float64)

	for buyRows.Next() {
		var day string
		var value float64
		if err := buyRows.Scan(&day, &value); err != nil {
			return nil, err
		}
		buyMap[day] = value
	}

	sellRows, err := dao.db.Table("sell, buy").
		Select("DATE_FORMAT(sell.created_at, '%Y-%m-%d') AS date, SUM(sell.price*sell.quantity) AS amount, SUM(sell.total_profit) AS total_profit").
		Where("buy.user_id = ? and sell.buy_id = buy.id", userId).
		Where("sell.created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer sellRows.Close()

	// 将结果转换为以日期为 key 的 map
	sellMap := make(map[string]float64)
	profitMap := make(map[string]float64)
	for sellRows.Next() {
		var day string
		var amount float64
		var profit float64

		if err := sellRows.Scan(&day, &amount, &profit); err != nil {
			return nil, err
		}
		sellMap[day] = amount
		profitMap[day] = profit
	}

	// 构建最终结果
	for i := 0; i < len(dates); i++ {
		date := dates[i]
		buyAmount, ok := buyMap[date]
		if !ok {
			buyAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: buyAmount,
			Type:  "买入",
		})
		sellAmount, ok := sellMap[date]
		if !ok {
			sellAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: sellAmount,
			Type:  "卖出",
		})
		profit, ok := profitMap[date]
		if !ok {
			profit = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: profit,
			Type:  "利润",
		})
	}
	// 按照 Label 排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Label < list[j].Label
	})

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.Statistics{}, nil
	}

	return list, nil
}
func (dao *StatisticsDAO) getProfitByWeek(userId uint64) ([]models.Statistics, error) {
	var list []models.Statistics
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := calculateDateRange(now, "week")

	// 构建日期范围
	dates := make([]string, 0)
	for t := start; !t.After(end); t = t.AddDate(0, 0, 7) {
		monday := t.AddDate(0, 0, -int(t.Weekday())+1)
		dates = append(dates, monday.Format("2006-01-02"))
	}
	fmt.Println("dates", dates)

	buyRows, err := dao.db.Table("buy").
		Select("DATE_FORMAT(DATE_SUB(created_at, INTERVAL WEEKDAY(created_at) DAY),'%Y-%m-%d') as monday, SUM(total_amount) AS amount").
		Where("user_id = ?", userId).
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("monday").
		Order("monday").
		Rows()

	if err != nil {
		return nil, err
	}
	defer buyRows.Close()

	// 将结果转换为以日期为 key 的 map
	buyMap := make(map[string]float64)

	for buyRows.Next() {
		var day string
		var value float64
		if err := buyRows.Scan(&day, &value); err != nil {
			return nil, err
		}
		buyMap[day] = value
	}

	sellRows, err := dao.db.Table("sell, buy").
		Select("DATE_FORMAT(DATE_SUB(sell.created_at, INTERVAL WEEKDAY(sell.created_at) DAY),'%Y-%m-%d') as monday, SUM(sell.price*sell.quantity) AS amount, SUM(sell.total_profit) AS total_profit").
		Where("buy.user_id = ? and sell.buy_id = buy.id", userId).
		Where("sell.created_at BETWEEN ? AND ?", start, end).
		Group("monday").
		Order("monday").
		Rows()

	if err != nil {
		return nil, err
	}
	defer sellRows.Close()

	// 将结果转换为以日期为 key 的 map
	sellMap := make(map[string]float64)
	profitMap := make(map[string]float64)
	for sellRows.Next() {
		var day string
		var amount float64
		var profit float64

		if err := sellRows.Scan(&day, &amount, &profit); err != nil {
			return nil, err
		}
		sellMap[day] = amount
		profitMap[day] = profit
	}

	// 构建最终结果
	for i := 0; i < len(dates); i++ {
		date := dates[i]
		buyAmount, ok := buyMap[date]
		if !ok {
			buyAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: buyAmount,
			Type:  "买入",
		})
		sellAmount, ok := sellMap[date]
		if !ok {
			sellAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: sellAmount,
			Type:  "卖出",
		})
		profit, ok := profitMap[date]
		if !ok {
			profit = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: profit,
			Type:  "利润",
		})
	}
	// 按照 Label 排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Label < list[j].Label
	})

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.Statistics{}, nil
	}

	return list, nil
}
func (dao *StatisticsDAO) getProfitByMonth(userId uint64) ([]models.Statistics, error) {
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := calculateDateRange(now, "month")
	// 分组字段
	// groupBy := groupByTimeDimension("month")
	// 构建日期范围
	dates := make([]string, 0)
	for t := start; !t.After(end); t = t.AddDate(0, 1, -1) {
		fmt.Println("t", t)
		dates = append(dates, t.Format("2006-01"))
	}
	fmt.Println("dates", start, end, dates)
	buyRows, err := dao.db.Table("buy").
		Select("DATE_FORMAT(created_at, '%Y-%m') AS date, SUM(total_amount) AS amount").
		Where("user_id = ?", userId).
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer buyRows.Close()

	// 将结果转换为以日期为 key 的 map
	buyMap := make(map[string]float64)

	for buyRows.Next() {
		var day string
		var value float64
		if err := buyRows.Scan(&day, &value); err != nil {
			return nil, err
		}
		buyMap[day] = value
	}
	sellRows, err := dao.db.Table("sell, buy").
		Select("DATE_FORMAT(sell.created_at, '%Y-%m') AS date, SUM(sell.price*sell.quantity) AS amount, SUM(sell.total_profit) AS total_profit").
		Where("buy.user_id = ? and sell.buy_id = buy.id", userId).
		Where("sell.created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer sellRows.Close()
	// 将结果转换为以日期为 key 的 map
	sellMap := make(map[string]float64)
	profitMap := make(map[string]float64)
	for sellRows.Next() {
		var day string
		var amount float64
		var profit float64

		if err := sellRows.Scan(&day, &amount, &profit); err != nil {
			return nil, err
		}
		sellMap[day] = amount
		profitMap[day] = profit
	}
	var list []models.Statistics
	// 构建最终结果
	for i := 0; i < len(dates); i++ {
		date := dates[i]
		buyAmount, ok := buyMap[date]
		if !ok {
			buyAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: buyAmount,
			Type:  "买入",
		})
		sellAmount, ok := sellMap[date]
		if !ok {
			sellAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: sellAmount,
			Type:  "卖出",
		})
		profit, ok := profitMap[date]
		if !ok {
			profit = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: profit,
			Type:  "利润",
		})
	}
	// 按照 Label 排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Label < list[j].Label
	})
	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.Statistics{}, nil
	}
	return list, nil
}
func (dao *StatisticsDAO) getProfitByYear(userId uint64) ([]models.Statistics, error) {
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := calculateDateRange(now, "year")
	// 分组字段
	// groupBy := groupByTimeDimension("month")
	// 构建日期范围
	dates := make([]string, 0)
	for t := start; !t.After(end); t = t.AddDate(1, 0, 0) {
		dates = append(dates, t.Format("2006"))
	}
	buyRows, err := dao.db.Table("buy").
		Select("DATE_FORMAT(created_at, '%Y') AS date, SUM(total_amount) AS amount").
		Where("user_id = ?", userId).
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer buyRows.Close()

	// 将结果转换为以日期为 key 的 map
	buyMap := make(map[string]float64)

	for buyRows.Next() {
		var day string
		var value float64
		if err := buyRows.Scan(&day, &value); err != nil {
			return nil, err
		}
		buyMap[day] = value
	}
	sellRows, err := dao.db.Table("sell, buy").
		Select("DATE_FORMAT(sell.created_at, '%Y') AS date, SUM(sell.price*sell.quantity) AS amount, SUM(sell.total_profit) AS total_profit").
		Where("buy.user_id = ? and sell.buy_id = buy.id", userId).
		Where("sell.created_at BETWEEN ? AND ?", start, end).
		Group("date").
		Order("date").
		Rows()

	if err != nil {
		return nil, err
	}
	defer sellRows.Close()
	// 将结果转换为以日期为 key 的 map
	sellMap := make(map[string]float64)
	profitMap := make(map[string]float64)
	for sellRows.Next() {
		var day string
		var amount float64
		var profit float64

		if err := sellRows.Scan(&day, &amount, &profit); err != nil {
			return nil, err
		}
		sellMap[day] = amount
		profitMap[day] = profit
	}
	var list []models.Statistics
	// 构建最终结果
	for i := 0; i < len(dates); i++ {
		date := dates[i]
		buyAmount, ok := buyMap[date]
		if !ok {
			buyAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: buyAmount,
			Type:  "买入",
		})
		sellAmount, ok := sellMap[date]
		if !ok {
			sellAmount = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: sellAmount,
			Type:  "卖出",
		})
		profit, ok := profitMap[date]
		if !ok {
			profit = 0
		}
		list = append(list, models.Statistics{
			Label: date,
			Value: profit,
			Type:  "利润",
		})
	}
	// 按照 Label 排序
	sort.Slice(list, func(i, j int) bool {
		return list[i].Label < list[j].Label
	})
	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.Statistics{}, nil
	}
	return list, nil
}
func (dao *StatisticsDAO) GetTotalProfit(userId uint64) (float64, error) {
	var totalProfit sql.NullFloat64
	err := dao.db.Model(&models.Buy{}).
		Select("SUM(total_profit)").
		Where("user_id = ?", userId).
		Scan(&totalProfit).Error
	if err != nil {
		return 0, err
	}
	if totalProfit.Valid {
		return totalProfit.Float64, nil
	}
	return 0, nil

}
