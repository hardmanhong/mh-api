package goods

import (
	"github.com/hardmanhong/utils"
	"gorm.io/gorm"
)

type GoodsDAO interface {
	GetList(query *ListQuery) (*ListResult, error)
	GetItem(id int) (*Goods, error)
	Create(goods *Goods) error
	Exists(id uint64) (bool, error)
	Update(id uint64, goods *GoodsUpdate) error
	Delete(id uint64) error
}

type GoodsDAOImpl struct {
	db *gorm.DB
}

func NewGoodsDAO(db *gorm.DB) GoodsDAO {
	return &GoodsDAOImpl{db}
}

func (dao *GoodsDAOImpl) GetList(query *ListQuery) (*ListResult, error) {
	var result ListResult
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
	err := db.Model(&Goods{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	var goodsList []*Goods
	err = db.Offset(offset).Limit(pageSize).Find(&goodsList).Error
	if err != nil {
		return nil, err
	}

	result.Total = total
	result.List = goodsList

	return &result, nil
}

func (dao *GoodsDAOImpl) GetItem(id int) (*Goods, error) {
	var goods Goods
	err := dao.db.Where("id = ?", id).First(&goods).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &goods, nil
}

func (dao *GoodsDAOImpl) Create(goods *Goods) error {
	return dao.db.Create(goods).Error
}
func (dao *GoodsDAOImpl) Exists(id uint64) (bool, error) {
	var count int64
	err := dao.db.Model(&Goods{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *GoodsDAOImpl) Update(id uint64, goods *GoodsUpdate) error {
	return dao.db.Model(&Goods{}).Where("id = ?", id).Updates(&Goods{
		Name:     goods.Name,
		MinPrice: goods.MinPrice,
		MaxPrice: goods.MaxPrice,
	}).Error
}
func (dao *GoodsDAOImpl) Delete(id uint64) error {
	return dao.db.Where("id = ?", id).Delete(&Goods{}).Error
}
