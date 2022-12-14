package service

import (
	"IShare/global"
	"IShare/model/database"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

func CreateUserTag(tag *database.Tag) (err error) {
	if err = global.DB.Create(&tag).Error; err != nil {
		return err
	}
	return nil
}

func QueryUserTagByName(user_id uint64, tag_name string) (existed bool) {
	tag := database.Tag{}
	err := global.DB.Where("user_id = ? AND tag_name = ?", user_id, tag_name).First(&tag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return true
	}
}

func QueryTagUser(tag_id uint64, user_id uint64) (tag database.Tag, notFound bool) {
	tag = database.Tag{}
	err := global.DB.Where("tag_id = ? AND user_id = ?", tag_id, user_id).First(&tag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tag, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tag, false
	}
}
func QueryPaperTag(tag_id uint64, paper_id string) (tag_paper database.TagPaper, notFound bool) {
	tag_paper = database.TagPaper{}
	err := global.DB.Where("tag_id = ? AND paper_id = ?", tag_id, paper_id).First(&tag_paper).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tag_paper, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tag_paper, false
	}
}

func GetTagById(tag_id uint64) (tag database.Tag, notFound bool) {
	tag = database.Tag{}
	err := global.DB.Where("tag_id = ?", tag_id).First(&tag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tag, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tag, false
	}
}
func CreateTagPaper(tagpaper *database.TagPaper) (err error) {
	//tagpaper = &database.TagPaper{TagID: tagpaper.TagID, PaperID: tagpaper.PaperID, CreateTime: tagpaper.CreateTime}
	if err = global.DB.Create(&tagpaper).Error; err != nil {
		return err
	}
	return nil
}

// 删除文章-标签收藏关系
func DeleteTagPaper(tag_id uint64, paper_id string) (err error) {
	if err = global.DB.Where("tag_id = ? AND paper_id = ?", tag_id, paper_id).Delete(database.TagPaper{}).Error; err != nil {
		return err
	}
	return nil
}

func QueryTagPaper(tagID uint64) (papers []database.TagPaper) {
	papers = make([]database.TagPaper, 0)
	global.DB.Where("tag_id=?", tagID).Order("create_time desc").Find(&papers)
	return papers
}

// 查询用户所有标签
func QueryTagList(userID uint64) (tags []database.Tag) {
	tags = make([]database.Tag, 0)
	global.DB.Where("user_id=?", userID).Find(&tags)
	return tags
}

func QueryATag(userID uint64, tagName string) (tag database.Tag, notFound bool) {
	db := global.DB
	db = db.Where("user_id = ?", userID)
	db = db.Where("tag_name = ?", tagName)
	err := db.First(&tag).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tag, true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tag, false
	}
}

// 删除标签
func DeleteTag(tagID uint64) (err error) {
	if err = global.DB.Where("tag_id = ?", tagID).Delete(database.Tag{}).Error; err != nil {
		return err
	}
	if err = global.DB.Where("tag_id = ?", tagID).Delete(database.TagPaper{}).Error; err != nil {
		return err
	}
	return nil
}

func RenameTag(new_name string, tag database.Tag) (err error) {
	old_name := tag.TagName
	if old_name == new_name {
		return nil
	}
	tag.TagName = new_name
	if err = global.DB.Save(tag).Error; err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetPaperStarNum(paper_id string) (num int, relate []database.TagPaper) {
	relate = make([]database.TagPaper, 0)
	global.DB.Where("paper_id = ? ", paper_id).Find(&relate)
	return len(relate), relate
}

func QueryPaperIsCollect(paper_id string, tag_id uint64) (tag_paper database.TagPaper, isCollect bool) {
	tag_paper = database.TagPaper{}
	err := global.DB.Where("tag_id = ? AND paper_id = ?", tag_id, paper_id).First(&tag_paper).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return tag_paper, false
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return tag_paper, true
	}
}
