package goods

type Goods struct {
	ID       uint64 `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	MinPrice int64  `json:"minPrice" gorm:"column:min_price"`
	MaxPrice int64  `json:"maxPrice" gorm:"column:max_price"`
}
type GoodsUpdate struct {
	Name     string `json:"name"`
	MinPrice int64  `json:"minPrice" gorm:"column:min_price"`
	MaxPrice int64  `json:"maxPrice" gorm:"column:max_price"`
}
