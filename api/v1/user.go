package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Register 注册
// @Summary     ccf
// @Description 填入用户名和密码注册
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       data body     response.RegisterQ true "data"
// @Success     200  {string} json               "{"status":200,"msg":"register success"}"
// @Failure     400  {string} json               "{"status":400,"msg":"username exists"}"
// @Router      /register [POST]
func Register(c *gin.Context) {
	// 获取请求数据
	var d response.RegisterQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	// 用户的用户名已经注册过的情况
	if _, notFound := service.GetUserByUsername(d.Username); !notFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
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
	})
}

// Login 登录
// @Summary     ccf
// @Description 登录
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       data body     response.LoginQ true "data"
// @Success     200 {string} json "{"status":200,"success":true,"msg":"login success","token": 666}"
// @Failure     400 {string} json "{"status":400,"success":false,"msg":"username doesn't exist"}"
// @Failure     401 {string} json "{"status":401,"success":false,"msg":"password doesn't match"}"
// @Router      /login [POST]
func Login(c *gin.Context) {
	var d response.LoginQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	//password := c.PostForm("password")	//不用data在hash比较时候会出错？？？
	//password := c.Request.FormValue("password")
	//data := utils.BindJsonAndValid(c, &response.LoginQ{}).(*response.LoginQ)
	//password := data.Password
	//fmt.Print(data)
	// 用户不存在
	user, notFound := service.GetUserByUsername(d.Username)
	if notFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "username doesn't exist",
		})
		return
	}
	// 密码错误的情况
	//fmt.Print(user.Password)
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//fmt.Print(hashedPassword)
	//fmt.Print(password)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(d.Password)); err != nil {
		c.JSON(401, gin.H{
			"status": 401,
			"msg":    "password doesn't match",
		})
		return
	}
	// 成功返回响应
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "login success",
		"token":  666,
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
// @Failure     200 {string} json "{"status":201,"success":false,"msg":"userID not exist"}"
// @Router      /user/info [POST]
func UserInfo(c *gin.Context) {
	userID := c.Query("user_id")
	//userID := c.PostForm("userID")
	id, _ := strconv.ParseInt(userID, 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(uint64(id))
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  201,
			"msg":     "userID not exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
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
// @Failure     200 {string} json "{"status":201,"success":false,"msg":"用户ID不存在"}"
// @Failure     400 {string} json "{"status":202,"success":false,"msg":err.Error()}"
// @Router      /user/mod [POST]
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
			"status":  201,
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
			"status":  202,
			"msg":     err.Error(),
		})
		return
	}
	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  200,
		"msg":     "修改成功",
		"data":    user,
	})
}

// ModifyPassword 编辑用户信息
// @Summary     ccf
// @Description 编辑用户信息
// @Tags        用户
// @Param       user_id query string true "user_id"
// @Param       password_old    query string true "旧密码"
// @Param       password_new1 query string true "新密码1"
// @Param       password_new2 query string true "新密码2"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"success":true,"msg":"修改成功","data":{object}}"
// @Failure     200 {string} json "{"status":201,"success":false,"msg":"用户ID不存在"}"
// @Failure     200 {string} json "{"status":202,"success":false,"msg":"两次输入密码不一致"}"
// @Failure     200 {string} json "{"status":203,"success":false,"msg":"原密码输入错误"}"
// @Failure     400 {string} json "{"status":204,"success":false,"msg":err1.Error()}"
// @Router      /user/pwd [POST]
func ModifyPassword(c *gin.Context) {
	userId := c.Query("user_id")
	//userID, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	passwordOld := c.Query("password_old")
	passwordNew1 := c.Query("password_new1")
	passwordNew2 := c.Query("password_new2")
	userID, _ := strconv.ParseUint(userId, 0, 64)

	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  201,
			"msg":     "用户ID不存在",
		})
		return
	}
	if passwordNew1 != passwordNew2 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  202,
			"msg":     "两次输入密码不一致",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordOld)); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  203,
			"msg":     "原密码输入错误",
		})
		return
	}
	var password = passwordNew1
	// 将密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("CreateUser: hash password error")
	}
	user.Password = string(hashedPassword)
	err1 := global.DB.Save(user).Error
	if err1 != nil {
		c.JSON(400, gin.H{
			"success": false,
			"status":  204,
			"msg":     err1.Error(),
		})
		return
	}
	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "修改成功",
		"status":  200,
		"data":    user,
	})
}

// UploadHeadshot 上传用户头像
// @Summary     ccf
// @Description 上传用户头像
// @Tags        用户
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"success":true,"msg":"修改成功","data":{object}}"
// @Failure     200 {string} json "{"status":201,"success":false,"msg":"用户ID不存在"}"
// @Failure     200 {string} json "{"status":202,"success":false,"msg":"头像文件上传失败"}"
// @Failure     200 {string} json "{"status":203,"success":false,"msg":"文件保存失败"}"
// @Failure     200 {string} json "{"status":204,"success":false,"msg":"保存文件路径到数据库中失败"}"
// @Router      /user/headshot [POST]
func UploadHeadshot(c *gin.Context) {
	userId := c.Query("user_id")
	userID, _ := strconv.ParseUint(userId, 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  201,
			"msg":     "用户ID不存在",
		})
		return
	}
	//1、获取上传的文件
	file, header, fileErr := c.Request.FormFile("headshot")
	if fileErr != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  202,
			"msg":     "头像文件上传失败",
		})
		return
	}
	//2、将文件保存到本地
	filePath := "./media/headshot/" + header.Filename
	out, e := os.Create(filePath)
	if e != nil {
		fmt.Println(e)
		_ = os.Mkdir(filePath, 777)
		out, e = os.Create(filePath)
		if e != nil {
			fmt.Println(e)
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"status":  203,
				"msg":     "文件保存失败",
			})
		}
		return
	}
	defer out.Close()
	_, err := io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	//3、将文件对应路径更新到数据库中
	user.HeadShot = filePath
	err = global.DB.Save(user).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  204,
			"msg":     "保存文件路径到数据库中失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "头像上传成功",
		"status":  200,
		"data":    user,
	})
	return
}
