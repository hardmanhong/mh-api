package goods

type Goods struct {
	ID       uint64  `json:"id" gorm:"primary_key"`
	Name     string  `json:"name"`
	MinPrice float64 `json:"minPrice" gorm:"column:min_price;type:decimal(10,2)"`
	MaxPrice float64 `json:"maxPrice" gorm:"column:max_price;type:decimal(10,2)"`
}
type GoodsUpdate struct {
	Name     string  `json:"name"`
	MinPrice float64 `json:"minPrice" gorm:"column:min_price;type:decimal(10,2)"`
	MaxPrice float64 `json:"maxPrice" gorm:"column:max_price;type:decimal(10,2)"`
}
type ListQuery struct {
	Name     string `json:"name"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}
type ListResult struct {
	Total int64
	List  []*Goods
}
