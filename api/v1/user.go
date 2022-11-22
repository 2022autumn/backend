package v1

import (
	"IShare/model/database"
	"IShare/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetData(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	c.JSON(http.StatusOK, gin.H{
		"message":  "success",
		"username": username,
		"password": password,
	})
}
func Register(c *gin.Context) {
	// 获取请求数据
	username := c.PostForm("username")
	password1 := c.PostForm("password1")
	password2 := c.PostForm("password2")
	if password1 != password2 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "两次输入密码不同",
		})
		return
	}
	var password = password1
	// 用户的用户名已经注册过的情况
	if _, notFound := service.GetUserByUsername(username); !notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户已存在",
		})
		return
	}
	// 成功创建用户
	if err := service.CreateUser(&database.User{Name: username, Password: password}); err != nil {
		panic("CreateUser: create user error")
	}
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  "注册成功",
		"username": username,
		"password": password,
	})
}
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	// 用户的用户名已经注册过的情况
	user, notFound := service.GetUserByUsername(username)
	if notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "登录失败，用户名不存在",
		})
		return
	}
	// 密码错误的情况
	if user.Password != password {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "密码错误",
		})
		return
	}
	// 成功返回响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"token":   666,
	})
}
