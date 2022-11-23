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
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/api/test", testGin)

	//用户模块
	UserRouter := r.Group("api/v1")
	{
		UserRouter.POST("/register", v1.Register)  //注册
		UserRouter.POST("/login", v1.Login)        //登录
		UserRouter.POST("/userinfo", v1.UserInfo)  //个人中心
		UserRouter.POST("/usermod", v1.ModifyUser) //编辑个人信息
	}
	r.POST("/GetData", v1.GetData) //测试数据接收获取

	// baseGroup := r.Group("/api/v1")
	// {
	// 	baseGroup.POST("/register", v1.Register)
	// 	baseGroup.POST("/login", v1.Login)
	// 	baseGroup.Static("/media", "./media")
	// }

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
