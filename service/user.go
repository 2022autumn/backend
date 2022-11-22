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
	err := global.DB.Where("name = ?", username).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return user, true
	} else {
		return user, false
	}
}
