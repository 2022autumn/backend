package initialize

import (
	v1 "IShare/api/v1"
	"IShare/docs"
	"IShare/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine) {
	r.Use(middleware.Cors()) // 跨域
	// r.Use(middleware.LoggerToFile()) // 日志

	docs.SwaggerInfo.Title = "?"
	docs.SwaggerInfo.Version = "v1"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/api/test", testGin)

	baseGroup := r.Group("/api")
	{
		//用户模块
		baseGroup.POST("/register", v1.Register)  //注册
		baseGroup.POST("/login", v1.Login)        //登录
		baseGroup.POST("/userinfo", v1.UserInfo)  //个人中心
		baseGroup.POST("/usermod", v1.ModifyUser) //编辑个人信息
	}
	ApplicationRouter := baseGroup.Group("/application")
	{
		ApplicationRouter.POST("/create", v1.CreateApplication)
	}
	// {
	// 	baseGroup.Static("/media", "./media")
	// }
	esGroup := baseGroup.Group("/es")
	{
		esGroup.POST("test_es", v1.TestEsSearch)
		esGroup.GET("get/:id", v1.GetObject)
		esGroup.POST("basequery", v1.BaseSearch)
	}
	// userGroup := baseGroup.Group("/user", middleware.AuthRequired())
	// {
	// 	userGroup.POST("/upload_avatar", v1.UploadAvatar)
	// }

}

func testGin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"success": true,
	})
}
