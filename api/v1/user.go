package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
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

// Register 注册
func Register(c *gin.Context) {
	// 获取请求数据
	username := c.PostForm("username")
	password1 := c.PostForm("password_1")
	password2 := c.PostForm("password_2")
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
	// 将密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("CreateUser: hash password error")
	}
	// 成功创建用户
	if err := service.CreateUser(&database.User{Username: username, Password: string(hashedPassword)}); err != nil {
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

// Login 登录
func Login(c *gin.Context) {
	username := c.PostForm("username")
	//password := c.PostForm("password")	//不用data在hash比较时候会出错？？？
	//password := c.Request.FormValue("password")
	data := utils.BindJsonAndValid(c, &response.LoginQ{}).(*response.LoginQ)
	fmt.Print(data)
	// 用户不存在
	user, notFound := service.GetUserByUsername(username)
	if notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "登录失败，用户名不存在",
		})
		return
	}
	// 密码错误的情况
	fmt.Print(user.Password)
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//fmt.Print(hashedPassword)
	//fmt.Print(password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
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

// UserInfo 查看用户个人信息
func UserInfo(c *gin.Context) {
	userID := c.PostForm("userID")
	id, _ := strconv.ParseInt(userID, 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(uint64(id))
	if notFoundUserByID {
		c.JSON(404, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"message": "查看用户信息成功",
		"data":    user,
	})
}

// ModifyUser 编辑用户信息
func ModifyUser(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	userInfo := c.Request.FormValue("user_info")
	phoneNum := c.Request.FormValue("phone_number")
	email := c.Request.FormValue("email")

	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return
	}

	user.UserInfo = userInfo
	user.Phone = phoneNum
	user.Email = email
	err := global.DB.Save(user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"status":  500,
			"message": err.Error(),
		})
		return
	}
	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改成功",
		"status":  200,
		"data":    user,
	})
}

// ModifyPassword 编辑用户信息
func ModifyPassword(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	passwordOld := c.Request.FormValue("password_old")
	passwordNew1 := c.Request.FormValue("password_new1")
	passwordNew2 := c.Request.FormValue("password_new2")

	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  404,
			"message": "用户ID不存在",
		})
		return
	}

	if user.Password != passwordOld {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  401,
			"message": "原密码输入错误",
		})
		return
	}

	if passwordNew1 != passwordNew2 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "两次输入密码不一致",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "修改成功",
		"status":  200,
		"data":    user,
	})
}

//
