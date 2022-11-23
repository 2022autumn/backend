package database

import "time"

// User 用户
type User struct {
	UserID   uint64    `gorm:"primary_key; autoIncrement; not null;" json:"user_id"`
	Username string    `gorm:"size:32; not null; unique;" json:"username"` //用户名
	Password string    `gorm:"size:128; not null;" json:"password"`        //密码
	RegTime  time.Time `gorm:"autoCreateTime" json:"regTime"`              //注册时间

	UserType uint64 `gorm:"default:0" json:"user_type"` // 0: 普通用户，1: 认证机构用户,2 管理员

	HeadShot    string `gorm:"size:64;" json:"head_shot"` //头像url
	UserInfo    string `gorm:"size:64;" json:"user_info"` //个性签名
	Name        string `gorm:"size:32;" json:"name"`      //真实姓名
	Phone       string `gorm:"size:32;" json:"phone"`
	Email       string `gorm:"size:32;" json:"email"`         //邮箱
	Fields      string `gorm:"size:256;" json:"fields"`       //研究领域
	InterestTag string `gorm:"size:256;" json:"interest_tag"` //兴趣词

	AuthorName string `gorm:"size:64;" json:"author_name"`        //被申请作者姓名
	AuthorID   string `gorm:"type:varchar(32);" json:"author_id"` // 被申请的作者ID
}
