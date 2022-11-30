package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

// Register 注册
// @Summary     ccf
// @Description 填入用户名和密码注册
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       data body     response.RegisterQ true "data"
// @Success     200  {string} json               "{"status":200,"msg":"register success","userid": 666}"
// @Failure     200  {string} json               "{"status":201,"msg":"username exists"}"
// @Router      /register [POST]
func Register(c *gin.Context) {
	// 获取请求数据
	var d response.RegisterQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	// 用户的用户名已经注册过的情况
	if _, notFound := service.GetUserByUsername(d.Username); !notFound {
		c.JSON(http.StatusOK, gin.H{
			"status": 201,
			"msg":    "username exists",
		})
		return
	}
	// 将密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
	if err != nil {
		panic("CreateUser: hash password error")
	}
	user := database.User{
		Username: d.Username,
		Password: string(hashedPassword),
	}
	// 成功创建用户
	if err := service.CreateUser(&user); err != nil {
		panic("CreateUser: create user error")
	}
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "register success",
		"userid": user.UserID,
	})
}

// Login 登录
// @Summary     ccf
// @Description 登录
// @Tags        用户
// @Param       username query string true "username"
// @Param       password query string true "password"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"success":true,"msg":"login success","token": 666}"
// @Failure     200 {string} json "{"status":201,"success":false,"msg":"username doesn't exist"}"
// @Failure     200 {string} json "{"status":202,"success":false,"msg":"password doesn't match"}"
// @Router      /login [POST]
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	//password := c.PostForm("password")	//不用data在hash比较时候会出错？？？
	//password := c.Request.FormValue("password")
	//data := utils.BindJsonAndValid(c, &response.LoginQ{}).(*response.LoginQ)
	//password := data.Password
	//fmt.Print(data)
	// 用户不存在
	user, notFound := service.GetUserByUsername(username)
	if notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "username doesn't exist",
		})
		return
	}
	// 密码错误的情况
	fmt.Print(user.Password)
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//fmt.Print(hashedPassword)
	//fmt.Print(password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "password doesn't match",
		})
		return
	}
	// 成功返回响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "login success",
		"token":   666,
	})
}

// UserInfo 查看用户个人信息
// @Summary     ccf
// @Description 查看用户个人信息
// @Tags        用户
// @Param       user_id      query string true "user_id"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"success":true,"msg":"get UserInfo","data":{object}}"
// @Failure     200 {string} json "{"status":200,"success":false,"msg":"userID not exist"}"
// @Router      user/info [POST]
func UserInfo(c *gin.Context) {
	userID := c.Query("user_id")
	//userID := c.PostForm("userID")
	id, _ := strconv.ParseInt(userID, 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(uint64(id))
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "userID not exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "get UserInfo",
		"data":    user,
	})
}

// ModifyUser 编辑用户信息
// @Summary     ccf
// @Description 编辑用户信息
// @Tags        用户
// @Param       user_id query string true "user_id"
// @Param       user_info    query string true "个性签名"
// @Param       phone_number query string true "电话号码"
// @Param       email        query string true "Email"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"success":true,"msg":"修改成功","data":{object}}"
// @Failure     200 {string} json "{"status":200,"success":false,"msg":"用户ID不存在"}"
// @Failure     200 {string} json "{"status":200,"success":false,"msg":err.Error()}"
// @Router      user/mod [POST]
func ModifyUser(c *gin.Context) {
	userId := c.Query("user_id")
	userID, _ := strconv.ParseUint(userId, 0, 64)
	//userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	userInfo := c.Query("user_info")
	phoneNum := c.Query("phone_number")
	email := c.Query("email")

	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     "用户ID不存在",
		})
		return
	}
	if len(userInfo) != 0 {
		user.UserInfo = userInfo
	}
	if len(phoneNum) != 0 {
		user.Phone = phoneNum
	}
	if len(email) != 0 {
		user.Email = email
	}
	err := global.DB.Save(user).Error
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"msg":     err.Error(),
		})
		return
	}
	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "修改成功",
		"data":    user,
	})
}

// ModifyPassword 编辑用户信息
// @Summary     ccf
// @Description 编辑用户信息
// @Tags        用户
// @Accept      json
// @Produce     json
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
