package service

import (
	"IShare/global"
	"IShare/model/database"
	"errors"
	"github.com/jinzhu/gorm"
)

//func QueryApplicationByAuthor(author_id string) (submit database.Application, notFound bool) {
//	err := global.DB.Where("author_id = ? AND status = 1", author_id).First(&submit).Error
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		return submit, true
//	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		panic(err)
//	} else {
//		return submit, false
//	}
//}
//func QueryUserIsScholar(user_id uint64) (submit database.Application, notFound bool) {
//	err := global.DB.Where("user_id = ? AND status = 1", user_id).First(&submit).Error
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		return submit, true
//	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		panic(err)
//	} else {
//		return submit, false
//	}
//}
//func CreateApplication(submit *database.Application) (err error) {
//	if err = global.DB.Create(&submit).Error; err != nil {
//		return err
//	}
//	return nil
//}
//func GetApplicationByID(application_id uint64) (application database.Application, notFound bool) {
//	err := global.DB.First(&application, application_id).Error
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		return application, true
//	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return application, true
//	} else {
//		return application, false
//	}
//}

//func MakeUserScholar(user database.User, application database.Application) {
//	user.Email = application.Email
//	user.AuthorName = application.AuthorName
//	user.UserType = 1
//	user.Fields = application.Fields
//	user.AuthorID = application.AuthorID
//	err := global.DB.Save(&user).Error
//	if err != nil {
//		panic(err)
//	}
//}

//func QueryAllSubmit() (application []database.Application) {
//	global.DB.Find(&application)
//	return application
//}
//
//func QueryUncheckedSubmit() (applications []database.Application, notFound bool) {
//	err := global.DB.Where("status = ?", 0).Find(&applications).Error
//	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
//		return applications, true
//	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		panic(err)
//	} else {
//		return applications, false
//	}
//}
func GetAppsByUserID(userID uint64, status int) (applications []database.Application, err error) {
	err = global.DB.Where("user_id = ? AND status = ?", userID, status).Find(&applications).Error
	return
}
func GetAppByID(appID uint64) (application database.Application, notFound bool) {
	err := global.DB.First(&application, appID).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return application, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return application, false
	}
}
func GetApps(status int) (applications []database.Application, err error) {
	err = global.DB.Where("status = ?", status).Find(&applications).Error
	return
}
