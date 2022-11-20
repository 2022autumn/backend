package initialize

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"IShare/global"
)

func InitMySQL() {
	// 配置数据
	var addr, port, user, password, dbname string
	addr = global.VP.GetString("db.addr")
	fmt.Println("addr is: " + addr)
	port = global.VP.GetString("db.port")
	user = global.VP.GetString("db.user")
	password = global.VP.GetString("db.password")
	dbname = global.VP.GetString("db.dbname")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, addr, port, dbname)

	// 连接数据库
	var err error
	global.DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Errorf("数据库出问题啦: %s \n", err))
		return
	}

	// 迁移
	global.DB.AutoMigrate(
	//base
	// &database.User{},
	)

	// 检查数据库连接是否存在, 好像没啥用
	err = global.DB.DB().Ping()
	if err != nil {
		panic(fmt.Errorf("数据库出问题啦: %s \n", err))
		return
	}

	return
}

func Close() {
	err := global.DB.Close()
	if err != nil {
		panic(fmt.Errorf("数据库出问题啦: %s \n", err))
		return
	}
}
