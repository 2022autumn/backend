package database

import "time"

type VerifyCode struct {
	RecID   uint64    `gorm:"primary key;autoIncrement; not null;" json:"rec_id"`
	UserID  uint64    `gorm:"not null;" json:"user_id"`
	Code    string    `gorm:"not null" json:"code"`
	Email   string    `gorm:"size:32;" json:"email"` //邮箱
	GenTime time.Time `gorm:"type:datetime" json:"create_time"`
}
