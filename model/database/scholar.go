package database

import (
	"database/sql"
	"time"
)

type Author struct {
	AuthorId            string `gorm:"primary_key;type:varchar(150);not null;" json:"author_id"`
	AuthorName          string `gorm:"not null;size :64;type:varchar(100)" json:"author_name"`
	AuthorRank          string `gorm:"size:32;" json:"author_rank"`
	PaperNum            int    `gorm:"type:int;" json:"paper_num"`
	CitationNum         int    `gorm:"type:int;" json:"citation_num"`
	AuthorInstitutionID string `gorm:"size:32;" json:"author_institution_id"`
}

type Institution struct {
	InstitutionID   string `gorm:"primary_key;type:varchar(150);not null;" json:"institution_id"`
	InstitutionName string `gorm:"type:varchar(150);not null;" json:"institution_name"`
	HomePageURL     string `gorm:"type:varchar(150);" json:"homepage_url"`
	CountryCode     string `gorm:"type:varchar(150);" json:"country_code"`
	WorksCount      int    `gorm:"type:int;" json:"works_count"`
	CitedByCount    int    `gorm:"type:int;" json:"cited_by_count"`
}
type Venue struct {
	VenueID      string    `gorm:"primary_key;type:varchar(150);not null;" json:"venue_id"`
	ISSN         string    `gorm:"type:varchar(50);unique;" json:"issn"`
	DisplayName  string    `gorm:"type:varchar(150);not null" json:"Venue_display_name"`
	WorksCount   int       `gorm:"type:int;not null" json:"works_count"`
	CitedByCount int       `gorm:"type:int;not null" json:"cited_by_count"`
	HomePageURL  string    `gorm:"type:varchar(150);" json:"homepage_url"`
	VenueType    uint64    `gorm:"default:0;" json:"venue_type"` //0:journal 1:repository 2:conference 3:ebook_platform
	UpdatedTime  time.Time `gorm:"column:updated_time;type:datetime" json:"updated_time"`
	CreatedTime  time.Time `gorm:"column:created_time;type:datetime" json:"created_time"`
}

type AuthorConnection struct {
	ConnectionID uint64 `gorm:"primary_key; not null" json:"connection_id"`
	AuthorID1    string `gorm:"type:varchar(32);" json:"author_id1"`
	AuthorID2    string `gorm:"type:varchar(32)" json:"author_id2"`
}

type Application struct {
	ApplicationID uint64 `gorm:"primary_key;not null;" json:"application_id"`
	UserID        uint64 `gorm:"not null;" json:"user_id"` //申请者的用户id
	AuthorName    string `gorm:"not null;type:varchar(100);" json:"author_name"`
	AuthorID      string `gorm:"type:varchar(150);not null;" json:"author_id"`
	//Fields          string       `gorm:"type:varchar(256);" json:"fields"`
	Status  int    `gorm:"not null;default:0" json:"status"` //0:未处理；1：通过申请 2：未通过申请
	Content string `gorm:"type:varchar(255)" json:"content"`
	//WorksCount      int          `gorm:"type:int;" json:"works_count"`
	InstitutionName string       `gorm:"type:varchar(150);not null;" json:"institution_name"`
	Email           string       `gorm:"size:32;" json:"email"` //邮箱
	ApplyTime       time.Time    `gorm:"type:datetime;" json:"apply_time"`
	HandleTime      sql.NullTime `gorm:"type:datetime;" json:"handle_time"`
}
