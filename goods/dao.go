package goods

import (
	"github.com/hardmanhong/db"
	"github.com/hardmanhong/utils"
	"gorm.io/gorm"
)

type GoodsDAO interface {
	GetList(page, pageSize int) (*ListResult, error)
	GetItem(id int) (*Goods, error)
	Create(goods *Goods) error
	Exists(id uint64) (bool, error)
	Update(id uint64, goods *GoodsUpdate) error
	Delete(id uint64) error
}

type GoodsDAOImpl struct{}
type ListResult struct {
	Total int64
	List  []*Goods
}

func (dao *GoodsDAOImpl) GetList(page, pageSize int) (*ListResult, error) {
	// 查询总的条数
	var count int64
	err := db.DB.Model(&Goods{}).Count(&count).Error

	if err != nil {
		return nil, err
	}

	var list []*Goods
	err = db.DB.Limit(pageSize).Offset((page - 1) * pageSize).Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &ListResult{
		Total: count,
		List:  list,
	}, nil
}

func (dao *GoodsDAOImpl) GetItem(id int) (*Goods, error) {
	var goods Goods
	err := db.DB.Where("id = ?", id).First(&goods).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrNotFound
		}
		return nil, err
	}
	return &goods, nil
}

func (dao *GoodsDAOImpl) Create(goods *Goods) error {
	return db.DB.Create(goods).Error
}
func (dao *GoodsDAOImpl) Exists(id uint64) (bool, error) {
	var count int64
	err := db.DB.Model(&Goods{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func (dao *GoodsDAOImpl) Update(id uint64, goods *GoodsUpdate) error {
	return db.DB.Model(&Goods{}).Where("id = ?", id).Updates(&Goods{
		Name:     goods.Name,
		MinPrice: goods.MinPrice,
		MaxPrice: goods.MaxPrice,
	}).Error
}
func (dao *GoodsDAOImpl) Delete(id uint64) error {
	return db.DB.Where("id = ?", id).Delete(&Goods{}).Error
}
