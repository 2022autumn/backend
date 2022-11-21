package main

import (
	"IShare/global"
	"IShare/initialize"
	"github.com/gin-gonic/gin"
)

func main() {
	initialize.InitViper()

	initialize.InitMySQL()
	defer initialize.CloseMySQL()

	initialize.InitMedia()
	initialize.InitElasticSearch()

	r := gin.Default()
	initialize.SetupRouter(r)
	if err := r.Run(":" + global.VP.GetString("port")); err != nil {
		panic(err)
	}
}
