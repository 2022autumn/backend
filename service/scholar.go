package service

import (
	"IShare/global"
	"IShare/model/database"
)

func CreateUserConcept(uc *database.UserConcept) (err error) {
	err = global.DB.Create(uc).Error
	return err
}
func GetUserConcept(user_id uint64, concept_id string) (uc database.UserConcept, notFound bool) {
	notFound = global.DB.Where("user_id = ? AND concept_id = ?", user_id, concept_id).
		First(&uc).RecordNotFound()
	return uc, notFound
}
func GetUserConcepts(user_id uint64) (ucs []database.UserConcept, err error) {
	err = global.DB.Where("user_id = ?", user_id).Find(&ucs).Error
	return ucs, err
}
func DeleteUserConcept(uc *database.UserConcept) (err error) {
	err = global.DB.Delete(uc).Error
	return err
}
func GetWorkView(work_id string) (work database.WorkView, notFound bool) {
	notFound = global.DB.Where("work_id = ?", work_id).First(&work).RecordNotFound()
	return work, notFound
}
func SaveWorkView(work *database.WorkView) (err error) {
	err = global.DB.Save(work).Error
	return err
}
func CreateWorkView(work *database.WorkView) (err error) {
	err = global.DB.Create(work).Error
	return err
}
func GetHotWorks(size int) (works []database.WorkView, err error) {
	err = global.DB.Order("views desc").Limit(size).Find(&works).Error
	return works, err
}
func GetAuthor(author_id string) (author database.Author, notFound bool) {
	notFound = global.DB.Where("author_id = ?", author_id).First(&author).RecordNotFound()
	return author, notFound
}
