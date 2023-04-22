package dao

import (
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
func (dao *StatisticsDAO) GetStatistics(dType string) ([]models.Statistics, error) {
	if dType == "day" {
		return dao.getProfitByDay()
	} else if dType == "week" {
		return dao.getProfitByWeek()
	} else if dType == "month" {
		return dao.getProfitByMonth()
	} else if dType == "year" {
		return dao.getProfitByYear()
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

func (dao *StatisticsDAO) getProfitByDay() ([]models.Statistics, error) {
	var list []models.Statistics
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := calculateDateRange(now, "day")
	// 分组字段
	// groupBy := groupByTimeDimension("day")

	// 构建日期范围
	datess := make([]string, 0)
	dates := make(map[string]bool)
	for t := start; !t.After(end); t = t.AddDate(0, 0, 1) {
		datess = append(datess, t.Format("2006-01-02"))
		dates[t.Format("2006-01-02")] = true
	}

	buyRows, err := dao.db.Table("buy").
		Select("DATE(created_at) AS date, SUM(total_amount) AS amount").
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
		var day time.Time
		var value float64
		if err := buyRows.Scan(&day, &value); err != nil {
			return nil, err
		}
		buyMap[day.Format("2006-01-02")] = value
	}

	sellRows, err := dao.db.Table("sell").
		Select("DATE(created_at) AS date, SUM(price*quantity) AS amount, SUM(total_profit) AS total_profit").
		Where("created_at BETWEEN ? AND ?", start, end).
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
		var day time.Time
		var amount float64
		var profit float64

		if err := sellRows.Scan(&day, &amount, &profit); err != nil {
			return nil, err
		}
		sellMap[day.Format("2006-01-02")] = amount
		profitMap[day.Format("2006-01-02")] = profit
	}

	// 构建最终结果
	for i := 0; i < len(datess); i++ {
		date := datess[i]
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
func (dao *StatisticsDAO) getProfitByWeek() ([]models.Statistics, error) {
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

	var list []models.Statistics
	for rows.Next() {
		var label string
		var value float64
		if err := rows.Scan(&label, &value); err != nil {
			return nil, err
		}
		list = append(list, models.Statistics{Value: value, Label: label})
	}

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.Statistics{}, nil
	}
	return list, nil
}
func (dao *StatisticsDAO) getProfitByMonth() ([]models.Statistics, error) {
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := calculateDateRange(now, "month")
	// 分组字段
	// groupBy := groupByTimeDimension("month")

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

	var list []models.Statistics
	for rows.Next() {
		var label string
		var value float64
		if err := rows.Scan(&label, &value); err != nil {
			return nil, err
		}
		// 2006-01-02 15:04:05 奇葩的格式
		list = append(list, models.Statistics{Value: value, Label: label})
	}

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.Statistics{}, nil
	}
	return list, nil
}
func (dao *StatisticsDAO) getProfitByYear() ([]models.Statistics, error) {
	// 获取当前时间
	now := time.Now()
	// 计算时间范围
	start, end := calculateDateRange(now, "year")
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

	var list []models.Statistics
	for rows.Next() {
		var label string
		var value float64
		if err := rows.Scan(&label, &value); err != nil {
			return nil, err
		}
		// 2006-01-02 15:04:05 奇葩的格式
		list = append(list, models.Statistics{Value: value, Label: label})
	}

	if len(list) == 0 {
		// 如果没有数据则返回一个空的数组
		return []models.Statistics{}, nil
	}
	return list, nil
}
func (dao *StatisticsDAO) GetTotalProfit() (float64, error) {
	var totalProfit float64
	err := dao.db.Model(&models.Buy{}).Select("SUM(total_profit)").Pluck("SUM(total_profit)", &totalProfit).Error
	if err != nil {
		return 0, err
	}
	return totalProfit, nil
}
