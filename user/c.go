user := &models.User{Name: "Alice", Age: 20}
err := userDao.Insert(user)
if err != nil {
	// 处理错误
}

user.Age = 25
err = userDao.Update(user)
if err != nil {
	// 处理错误
}

err = userDao.Delete(user.ID)
if err != nil {
	// 处理错误
}

user, err = userDao.FindById(user.ID)
if err != nil {
	// 处理错误
}
