package main

import (
	"2022autumn/global"
	"2022autumn/initialize"

	"github.com/gin-gonic/gin"
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
