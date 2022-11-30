package service

import (
	"IShare/global"
	"IShare/model/database"
	"errors"
	"github.com/jinzhu/gorm"
)

func QueryApplicationByAuthor(author_id string) (submit database.Application, notFound bool) {
	err := global.DB.Where("author_id = ? AND status = 1", author_id).First(&submit).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return submit, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return submit, false
	}
}
func QueryUserIsScholar(user_id uint64) (submit database.Application, notFound bool) {
	err := global.DB.Where("user_id = ? AND status = 1", user_id).First(&submit).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return submit, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return submit, false
	}
}
func CreateApplication(submit *database.Application) (err error) {
	if err = global.DB.Create(&submit).Error; err != nil {
		return err
	}
	return nil
}
func GetApplicationByID(application_id uint64) (application database.Application, notFound bool) {
	err := global.DB.First(&application, application_id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return application, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return application, true
	} else {
		return application, false
	}
}
func MakeUserScholar(user database.User, application database.Application) {
	//user.Email = submit.WorkEmail
	//user.AuthorName = submit.AuthorName
	//user.Affiliation = submit.AffiliationName
	//user.UserType = 1
	//user.Fields = submit.Fields
	//user.HomePage = submit.HomePage
	//user.PaperCount += submit.PaperCount
	//user.AuthorID = submit.AuthorID
	//author := GetSimpleAuthors(append(make([]string, 0), submit.AuthorID))[0].(map[string]interface{})
	//user.PaperCount = int(author["paper_count"].(float64))
	//user.CitationCount = int(author["citation_count"].(float64))
	//err := global.DB.Save(&user).Error
	//if err != nil {
	//	panic(err)
	//}
}