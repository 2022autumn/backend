package service

import (
	"IShare/global"
	"IShare/model/database"
	"errors"

	"github.com/jinzhu/gorm"
)

// 数据库操作

// CreateUser 创建用户
func CreateUser(user *database.User) (err error) {
	if err = global.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// GetUserByID 根据用户 ID 查询某个用户
func GetUserByID(ID uint64) (user database.User, notFound bool) {
	err := global.DB.First(&user, ID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else {
		return user, false
	}
}

// GetUserByUsername 根据用户名查询某个用户
func GetUserByUsername(username string) (user database.User, notFound bool) {
	err := global.DB.Where("username = ?", username).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else {
		return user, false
	}
}

// QueryAUserByID 根据用户 ID 查询某个用户
func QueryAUserByID(userID uint64) (user database.User, notFound bool) {
	err := global.DB.Where("user_id = ?", userID).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return user, false
	}
}

// UpdateAUser 更新用户的用户名、密码、个人信息
func UpdateAUser(user *database.User, username string, password string, userInfo string) error {
	user.Username = username
	user.Password = password
	user.UserInfo = userInfo
	err := global.DB.Save(user).Error
	return err
}
func GetAllUser() (num int) {
	users := make([]database.User, 0)
	global.DB.Where("").Find(&users)
	return len(users)
}
