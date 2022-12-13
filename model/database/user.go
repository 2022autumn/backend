package database

import "time"

// User 用户
type User struct {
	UserID   uint64 `gorm:"primary_key; autoIncrement; not null;" json:"user_id"`
	Username string `gorm:"size:32; not null; unique;" json:"username"` //用户名
	Password string `gorm:"size:256; not null;" json:"password"`        //密码

	UserType uint64 `gorm:"default:0" json:"user_type"` // 0: 普通用户，1 管理员

	HeadShot string `gorm:"default:'default.jpg'" json:"head_shot"` //头像url
	UserInfo string `gorm:"size:64;" json:"user_info"`              //个性签名
	Name     string `gorm:"size:32;" json:"name"`                   //真实姓名
	Phone    string `gorm:"size:32;" json:"phone"`                  //电话号码
	Email    string `gorm:"size:32;" json:"email"`                  //邮箱
	Fields   string `gorm:"size:256;" json:"fields"`                //研究领域

	AuthorName string `gorm:"size:64;" json:"author_name"`        //被申请作者姓名
	AuthorID   string `gorm:"type:varchar(32);" json:"author_id"` // 被申请的作者ID
	Verified   int    `gorm:"default:0" json:"verified"`          //是否已经认证
}
type BrowseHistory struct {
	UserID          uint64    `gorm:"not null;" json:"user_id"`
	WorkID          string    `gorm:"not null;" json:"work_id"`
	Title           string    `gorm:"not null;type:text" json:"title"`
	HostVenue       string    `gorm:"not null;type:text" json:"host_venue"`
	PublicationYear string    `gorm:"not null;" json:"publication_year"`
	BrowseTime      time.Time `gorm:"not null;" json:"browse_time"`
}
