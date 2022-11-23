package global

import (
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
)

var (
	DB *gorm.DB
	VP *viper.Viper
	ES *elastic.Client
)
