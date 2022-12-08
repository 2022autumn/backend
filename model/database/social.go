package database

import "time"

type Comment struct {
	CommentID   uint64    `gorm:"primary_key;autoIncrement;" json:"comment_id"`
	Content     string    `gorm:"size:255;" json:"content"`
	UserID      uint64    `gorm:"not null;" json:"user_id"` //评论者的用户id
	CommentTime time.Time `gorm:"column:comment_time;type:datetime" json:"comment_time"`
	PaperID     string    `gorm:"size:64;" json:"paper_id"`
	//PaperTitle  string    `gorm:"type:varchar(256);" json:"paper_title"`
	LikeNum uint64 `gorm:"default:0" json:"like_num"` //	点赞数量
}

type Like struct {
	CommentID uint64 `gorm:"primary_key;" json:"comment_id"`
	UserID    uint64 `gorm:"not null;" json:"user_id"` //评论者的用户id
}

//type Fav struct {
//	FavID  string `gorm:"primary_key;not null;unique;type:varchar(150)" json:"fav_id"`
//	UserID uint64 `gorm:"primary_key;not null;" json:"user_id"` //评论者的用户id
//}

type Tag struct {
	TagID      string    `gorm:"primary_key;not null;unique;autoIncrement" json:"tag_id"`
	TagName    string    `gorm:"type:varchar(50);not null" json:"tag_name"`
	UserID     uint64    `gorm:"not null;" json:"user_id"` //评论者的用户id
	CreateTime time.Time `gorm:"type:datetime" json:"create_time"`
}

type TagPaper struct {
	RelationID uint64    `gorm:"primary_key;type:integer;not null;autoIncrement" json:"relation_id"`
	TagID      string    `json:"tag_id"`
	PaperID    string    `json:"paper_id"`
	CreateTime time.Time `gorm:"type:datetime" json:"create_time"`
}
