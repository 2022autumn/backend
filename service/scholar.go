package service

import (
	"IShare/global"
	"IShare/model/database"
	"log"

	"github.com/jinzhu/gorm"
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

// 添加学者的作品
func AddScholarWork(work *database.PersonalWorks) (err error) {
	err = global.DB.Create(work).Error
	return err
}

// 查询学者的作品
func GetScholarWorks(author_id string) (works []database.PersonalWorks, notFound bool) {
	notFound = global.DB.Where("author_id = ?", author_id).Find(&works).RecordNotFound()
	return
}

// 修改作品的place
func UpdateWorkPlace(work_id string, place string) (err error) {
	err = global.DB.Model(&database.PersonalWorks{}).Where("work_id = ?", work_id).Update("place", place).Error
	return err
}

// 修改作品的ignore
func UpdateWorkIgnore(author_id string, work_id string, ignore bool) (err error) {
	err = global.DB.Model(&database.PersonalWorks{}).Where("author_id = ? AND work_id = ?", author_id, work_id).Update("ignore", !ignore).Error
	return err
}

// 获取作品的当前place
func GetWorkPlace(author_id string, work_id string) (place int, notFound bool) {
	var work database.PersonalWorks
	notFound = global.DB.Where("author_id = ? AND work_id = ?", author_id, work_id).First(&work).RecordNotFound()
	return work.Place, notFound
}

// 获取学者的作品总数
func GetScholarWorksCount(author_id string) (count int, err error) {
	err = global.DB.Model(&database.PersonalWorks{}).Where("author_id = ?", author_id).Count(&count).Error
	return count, err
}

// 加锁，交换两个作品的place
func SwapWorkPlace(author_id string, work_id1 string, work_id2 string) (err error) {
	tx := global.DB.Begin()
	var work1 database.PersonalWorks
	var work2 database.PersonalWorks
	err = tx.Where("author_id = ? AND work_id = ?", author_id, work_id1).First(&work1).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Where("author_id = ? AND work_id = ?", author_id, work_id2).First(&work2).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	p1 := work1.Place
	p2 := work2.Place
	log.Println(p1, p2)
	err = tx.Model(&database.PersonalWorks{}).Where("author_id = ? AND work_id = ?", author_id, work_id1).Update("place", p2).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&database.PersonalWorks{}).Where("author_id = ? AND work_id = ?", author_id, work_id2).Update("place", p1).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// 通过place获取作品
func GetWorkByPlace(author_id string, place int) (work database.PersonalWorks, notFound bool) {
	notFound = global.DB.Where("author_id = ? AND place = ?", author_id, place).First(&work).RecordNotFound()
	return work, notFound
}

// 置顶作品
func TopWork(author_id string, work_id string) (err error) {
	tx := global.DB.Begin()
	var work database.PersonalWorks
	err = tx.Where("author_id = ? AND work_id = ?", author_id, work_id).First(&work).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// err = tx.Model(&database.PersonalWorks{}).Where("author_id = ? AND work_id = ?", author_id, work_id).Update("place", 0).Error
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	err = tx.Model(&database.PersonalWorks{}).Where("author_id = ? AND place < ?", author_id, work.Place).Update("place", gorm.Expr("place + 1")).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&database.PersonalWorks{}).Where("author_id = ? AND work_id = ?", author_id, work_id).Update("place", 0).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetAuthor(author_id string) (author database.Author, notFound bool) {
	notFound = global.DB.Where("author_id = ?", author_id).First(&author).RecordNotFound()
	return author, notFound
}

// 批量创建作者的作品, 加锁
func CreateWorks(works []database.PersonalWorks) (err error) {
	tx := global.DB.Begin()
	for _, work := range works {
		err = tx.Create(&work).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return err
}

// 修改作品ignore属性 加锁
func IgnoreWork(author_id string, work_id string) (err error) {
	tx := global.DB.Begin()
	var work database.PersonalWorks
	err = tx.Where("author_id = ? AND work_id = ?", author_id, work_id).First(&work).Error
	if err != nil {
		tx.Rollback()
		log.Println("work not found", author_id, work_id, err)
		return err
	}
	preIgnore := work.Ignore
	err = tx.Model(&work).Update("ignore", !preIgnore).Error
	if err != nil {
		tx.Rollback()
		log.Println("update ignore failed", author_id, work_id, err)
		return err
	}
	tx.Commit()
	return nil
}
