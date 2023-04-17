package dao

import (
	"errors"

	"github.com/hardmanhong/api/models"
	"gorm.io/gorm"
)

type GoodsDAO struct {
	db *gorm.DB
}

func NewGoodsDAO(db *gorm.DB) *GoodsDAO {
	return &GoodsDAO{db}
}

func (dao *GoodsDAO) GetList(query *models.GoodsListQuery) (*models.PaginationResponse, error) {
	response := models.PaginationResponse{
		Total: 0,
		List:  make([]interface{}, 0),
	}
	db := dao.db
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	// 分页查询
	page := query.Page
	if page <= 0 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var total int64
	err := db.Model(&models.Goods{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	var goodsList []*models.Goods
	err = db.Offset(offset).Limit(pageSize).Find(&goodsList).Error
	if err != nil {
		return nil, err
	}

	response.Total = total
	// 将 Goods 类型的切片转换成 interface{} 类型的切片
	for _, g := range goodsList {
		response.List = append(response.List, g)
	}

	return &response, nil
}
func (dao *GoodsDAO) GetItem(id uint64) (*models.Goods, error) {
	var goods models.Goods
	err := dao.db.Where("id = ?", id).First(&goods).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("not found")
		}
		return nil, err
	}
	return &goods, nil
}
func (dao *GoodsDAO) Create(goods *models.Goods) (*models.Goods, error) {
	err := dao.db.Create(goods).Error
	if err != nil {
		return nil, err
	}
	return goods, err
}
func (dao *GoodsDAO) ExistsByName(name string) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Goods{}).Where("name = ?", name).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *GoodsDAO) Exists(id uint64) (bool, error) {
	var count int64
	err := dao.db.Model(&models.Goods{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *GoodsDAO) Update(id uint64, goods *models.GoodsUpdate) error {
	return dao.db.Model(&models.Goods{}).Where("id = ?", id).Updates(&models.Goods{
		Name:     goods.Name,
		MinPrice: goods.MinPrice,
		MaxPrice: goods.MaxPrice,
	}).Error
}

func (dao *GoodsDAO) Delete(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&models.Goods{}).Error

}
