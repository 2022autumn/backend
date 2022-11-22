package router

import (
	v1 "IShare/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	//用户模块
	UserRouter := r.Group("api/v1")
	{
		UserRouter.POST("/register", v1.Register) //注册
		UserRouter.POST("/login", v1.Login)       //登录
	}
	r.POST("/GetData", v1.GetData) //测试数据接收获取
}
