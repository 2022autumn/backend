package service

import (
	"IShare/global"
	"IShare/model/database"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
)

func CreateComment(comment *database.Comment) (err error) {
	if err := global.DB.Create(&comment).Error; err != nil {
		return err
	}
	return nil
}
func GetCommentByID(comment_id uint64) (comment *database.Comment, notFound bool) {
	comment = &database.Comment{}
	err := global.DB.First(&comment, comment_id).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return comment, true
	} else {
		return comment, false
	}
}

func GetLike_Rel(comment_id uint64, user_id uint64) (isLike bool) {
	like := database.Like{}
	err := global.DB.Where("user_id = ? AND comment_id = ?", user_id, comment_id).First(&like).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		return true
	}
}
func UpdateCommentLike(comment *database.Comment, user database.User) (err error) {
	fmt.Println(comment.LikeNum)
	comment.LikeNum++
	err = global.DB.Save(comment).Error
	if err != nil {
		return err
	}

	like := database.Like{UserID: user.UserID, CommentID: comment.CommentID}
	err = global.DB.Create(&like).Error
	return err
}

// 取消点赞
func CancelLike(comment *database.Comment, user database.User) (notFound bool) {
	like := database.Like{}
	err := global.DB.Where("user_id = ? AND comment_id = ?", user.UserID, comment.CommentID).First(&like).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	} else {
		global.DB.Delete(&like)
		comment.LikeNum--
		global.DB.Save(&comment)
		return false
	}
}

// 根据文献id获取文献所有评论
func GetCommentsByPaperId(paperId string) (comments []database.Comment) {
	comments = make([]database.Comment, 0)
	global.DB.Where(map[string]interface{}{"paper_id": paperId}).Order("comment_time desc").Find(&comments)
	return comments
}

func GetUserFollow(userID uint64, authorID string) (uf database.UserFollow, notFound bool) {
	notFound = global.DB.Where("user_id = ? AND author_id = ?", userID, authorID).
		First(&uf).RecordNotFound()
	return uf, notFound
}
func CreateUserFollow(uf *database.UserFollow) (err error) {
	err = global.DB.Create(uf).Error
	return err
}
func DeleteUserFollow(uf *database.UserFollow) (err error) {
	err = global.DB.Delete(uf, "user_id = ? AND author_id = ?", uf.UserID, uf.AuthorID).Error
	return err
}
func GetUserFollows(userID uint64) (userFollows []database.UserFollow, err error) {
	err = global.DB.Where("user_id = ?", userID).Find(&userFollows).Error
	return
}
