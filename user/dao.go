package dao

import "github.com/hardmanhong/user/models"

// UserDAO 用户数据访问接口
type UserDAO interface {
	Insert(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
	FindById(id int) (*models.User, error)
}

// UserDAOImpl 用户数据访问实现
type UserDAOImpl struct {
}

// Insert 新增用户
func (d *UserDAOImpl) Insert(user *models.User) error {
	// TODO: 插入数据库
	return nil
}

// Update 更新用户
func (d *UserDAOImpl) Update(user *models.User) error {
	// TODO: 更新数据库
	return nil
}

// Delete 删除用户
func (d *UserDAOImpl) Delete(id int) error {
	// TODO: 删除数据库中指定 id 的用户
	return nil
}

// FindById 根据 id 查找用户
func (d *UserDAOImpl) FindById(id int) (*models.User, error) {
	// TODO: 查找数据库中指定 id 的用户并返回
	return nil, nil
}
