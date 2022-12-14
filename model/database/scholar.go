package database

import (
	"time"
)

type Author struct {
	AuthorID string `gorm:"primary_key;not null;" json:"author_id"`
	HeadShot string `gorm:"default:'author_default.jpg'" json:"head_shot"` //头像url
	Intro    string `gorm:"type:text;" json:"intro"`
}

//type Institution struct {
//	InstitutionID   string `gorm:"primary_key;type:varchar(150);not null;" json:"institution_id"`
//	InstitutionName string `gorm:"type:varchar(150);not null;" json:"institution_name"`
//	HomePageURL     string `gorm:"type:varchar(150);" json:"homepage_url"`
//	CountryCode     string `gorm:"type:varchar(150);" json:"country_code"`
//	WorksCount      int    `gorm:"type:int;" json:"works_count"`
//	CitedByCount    int    `gorm:"type:int;" json:"cited_by_count"`
//}
//type Venue struct {
//	VenueID      string    `gorm:"primary_key;type:varchar(150);not null;" json:"venue_id"`
//	ISSN         string    `gorm:"type:varchar(50);unique;" json:"issn"`
//	DisplayName  string    `gorm:"type:varchar(150);not null" json:"Venue_display_name"`
//	WorksCount   int       `gorm:"type:int;not null" json:"works_count"`
//	CitedByCount int       `gorm:"type:int;not null" json:"cited_by_count"`
//	HomePageURL  string    `gorm:"type:varchar(150);" json:"homepage_url"`
//	VenueType    uint64    `gorm:"default:0;" json:"venue_type"` //0:journal 1:repository 2:conference 3:ebook_platform
//	UpdatedTime  time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
//	CreatedTime  time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
//}

//type AuthorConnection struct {
//	ConnectionID uint64 `gorm:"primary_key; not null" json:"connection_id"`
//	AuthorID1    string `gorm:"type:varchar(32);" json:"author_id1"`
//	AuthorID2    string `gorm:"type:varchar(32)" json:"author_id2"`
//}

type Application struct {
	ApplicationID uint64    `gorm:"primary_key;not null;" json:"application_id"`
	UserID        uint64    `gorm:"not null;" json:"user_id"` //申请者的用户id
	Username      string    `gorm:"not null;" json:"username"`
	RealName      string    `gorm:"not null;type:varchar(100);" json:"real_name"`
	AuthorID      string    `gorm:"not null;" json:"author_id"`
	Status        int       `gorm:"not null;default:0" json:"status"` //0:未处理；1：通过申请 2：未通过申请
	Content       string    `gorm:"type:text" json:"content"`
	Institution   string    `gorm:"not null;" json:"institution"`
	Email         string    `gorm:"not null;" json:"email"` //邮箱
	VerifyCode    string    `gorm:"not null;" json:"verify_code"`
	ApplyTime     time.Time `gorm:"type:datetime;default:Now()" json:"apply_time"`
	HandleTime    time.Time `gorm:"type:datetime;default:Now()" json:"handle_time"`
	HandleContent string    `gorm:"type:text;" json:"handle_content"`
}
type UserConcept struct {
	UserID      uint64 `gorm:"not null;" json:"user_id"`
	ConceptID   string `gorm:"not null;" json:"concept_id"`
	ConceptName string `gorm:"not null;" json:"concept_name"`
}
type WorkView struct {
	WorkID    string `gorm:"primary_key;not null;" json:"work_id"`
	Views     int    `gorm:"not null;default:0" json:"views"`
	WorkTitle string `gorm:"not null;" json:"work_title"`
}

type PersonalWorks struct {
	AuthorID string `gorm:"not null;" json:"author_id"`
	WorkID   string `gorm:"not null;" json:"work_id"`
	Place    int    `gorm:"not null;" json:"place"`
	Ignore   bool   `gorm:"not null;default:false" json:"ignore"`
	PDF      string `gorm:"default:''" json:"pdf"`
	Top      int    `gorm:"not null;default:-1" json:"top"`
}

type PersonalWorksCount struct {
	AuthorID string `gorm:"primary_key;not null;" json:"author_id"`
	Count    int    `gorm:"not null;default:0" json:"count"`
}
