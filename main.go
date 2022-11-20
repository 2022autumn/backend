package main

import (
	"github.com/gin-gonic/gin"

	"IShare/global"
	"IShare/initialize"
)

func main() {
	initialize.InitViper()

	initialize.InitMySQL()
	defer initialize.Close()

	initialize.InitMedia()

	r := gin.Default()
	initialize.SetupRouter(r)
	if err := r.Run(":" + global.VP.GetString("port")); err != nil {
		panic(err)
	}
}
