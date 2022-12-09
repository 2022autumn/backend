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
	"io"
	"net/http"
	"os"
	"strconv"
)

// Register 注册
// @Summary     ccf
// @Description 注册
// @Description 填入用户名和密码注册
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       data body     response.RegisterQ true "data"
// @Success     200  {string} json "{"status":200,"msg":"注册成功"}"
// @Failure     400  {string} json "{"status":400,"msg":"用户名已存在"}"
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
			"msg":    "用户名已存在",
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
		"msg":    "注册成功",
	})
}

// Login 登录
// @Summary     ccf
// @Description 登录
// @Description 填入用户名和密码
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       data body     response.LoginQ true "data"
// @Success     200 {string} json "{"status":200,"msg":"登录成功","token": token,"ID": user.UserID}"
// @Failure     400 {string} json "{"status":400,"msg":"用户名不存在"}"
// @Failure     401 {string} json "{"status":401,"msg":"密码错误"}"
// @Router      /login [POST]
func Login(c *gin.Context) {
	var d response.LoginQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	// 用户不存在
	user, notFound := service.GetUserByUsername(d.Username)
	if notFound {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户名不存在",
		})
		return
	}
	// 密码错误的情况
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(d.Password)); err != nil {
		c.JSON(401, gin.H{
			"status": 401,
			"msg":    "密码错误",
		})
		return
	}
	// 成功返回响应
	//token := 666
	token := utils.GenerateToken(user.UserID)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "登录成功",
		"token":  token,
		"ID":     user.UserID,
	})
}

// UserInfo 查看用户个人信息
// @Summary     ccf
// @Description 查看用户个人信息
// @Tags        用户
// @Param       user_id query string true "user_id"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"msg":"获取用户信息成功","data":{object}}"
// @Failure     400 {string} json "{"status":400,"msg":"用户ID不存在"}"
// @Router      /user/info [GET]
func UserInfo(c *gin.Context) {
	//GET
	userID := c.Query("user_id")
	id, _ := strconv.ParseInt(userID, 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(uint64(id))
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "获取用户信息成功",
		"data":   user,
	})
}

// ModifyUser 编辑用户信息
// @Summary     ccf
// @Description 编辑用户信息
// @Tags        用户
// @Param       data body	response.ModifyQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"msg":"修改个人信息成功","data":{object}}"
// @Failure     400 {string} json "{"status":400,"msg":"用户ID不存在"}"
// @Failure     401 {string} json "{"status":401,"msg":err.Error()}"
// @Router      /user/mod [POST]
func ModifyUser(c *gin.Context) {
	//userId := c.Query("user_id")
	//获取修改信息
	var d response.ModifyQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	userId := d.ID
	userInfo := d.UserInfo
	name := d.Name
	phoneNum := d.Phone
	email := d.Email
	fields := d.Fields
	// 用户不存在
	userID, _ := strconv.ParseUint(userId, 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	// 修改用户信息
	if len(userInfo) != 0 {
		user.UserInfo = userInfo
	}
	if len(name) != 0 {
		user.Name = name
	}
	if len(phoneNum) != 0 {
		user.Phone = phoneNum
	}
	if len(email) != 0 {
		user.Email = email
	}
	if len(fields) != 0 {
		user.Fields = fields
	}
	//成功修改数据库
	err := global.DB.Save(user).Error
	if err != nil {
		c.JSON(401, gin.H{
			"status": 401,
			"msg":    err.Error(),
		})
		return
	}
	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "修改个人信息成功",
		"data":   user,
	})
}

// ModifyPassword 编辑用户密码
// @Summary     ccf
// @Description 编辑用户信息
// @Tags        用户
// @Param       data body	response.PwdModifyQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"msg":"修改密码成功","data":{object}}"
// @Failure     400 {string} json "{"status":400,"msg":"用户ID不存在"}"
// @Failure     401 {string} json "{"status":401,"msg":"原密码输入错误"}"
// @Failure     402 {string} json "{"status":402,"msg":err1.Error()}"
// @Router      /user/pwd [POST]
func ModifyPassword(c *gin.Context) {
	//userId := c.Query("user_id")
	//userID, _ := strconv.ParseUint(userId, 0, 64)

	var d response.PwdModifyQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	userId := d.ID
	passwordOld := d.PasswordOld
	passwordNew := d.PasswordNew
	//用户ID不存在
	userID, _ := strconv.ParseUint(userId, 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordOld)); err != nil {
		c.JSON(401, gin.H{
			"status": 401,
			"msg":    "原密码输入错误",
		})
		return
	}
	var password = passwordNew
	// 将密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic("CreateUser: hash password error")
	}
	user.Password = string(hashedPassword)
	err1 := global.DB.Save(user).Error
	if err1 != nil {
		c.JSON(402, gin.H{
			"status": 402,
			"msg":    err1.Error(),
		})
		return
	}
	//修改成功
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "修改密码成功",
		"data":   user,
	})
}

// UploadHeadshot 上传用户头像
// @Summary     ccf
// @Description 上传用户头像
// @Tags        用户
// @Param       user_id formData string true "用户ID"
// @Param       Headshot formData file true "新头像"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"msg":"修改成功","data":{object}}"
// @Failure     400 {string} json "{"status":400,"msg":"用户ID不存在"}"
// @Failure     401 {string} json "{"status":401,"msg":"头像文件上传失败"}"
// @Failure     402 {string} json "{"status":402,"msg":"文件保存失败"}"
// @Failure     403 {string} json "{"status":403,"msg":"保存文件路径到数据库中失败"}"
// @Router      /user/headshot [POST]
func UploadHeadshot(c *gin.Context) {
	//userId := c.Query("user_id")
	userId := c.Request.FormValue("user_id")
	userID, _ := strconv.ParseUint(userId, 0, 64)
	/*
		var d response.AvatarQ
		if err := c.ShouldBind(&d); err != nil {
			panic(err)
		}
		userId := d.ID
		userID, _ := strconv.ParseUint(userId, 0, 64)
	*/
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	//1、获取上传的文件
	file, header, fileErr := c.Request.FormFile("Headshot")
	if fileErr != nil {
		c.JSON(401, gin.H{
			"status": 401,
			"msg":    "头像文件上传失败",
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
			c.JSON(402, gin.H{
				"status": 402,
				"msg":    "文件保存失败",
			})
		}
		return
	}
	//defer out.Close()
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {

		}
	}(out)
	_, err := io.Copy(out, file)
	if err != nil {
		fmt.Println(err)
		return
	}
	//3、将文件对应路径更新到数据库中
	user.HeadShot = "http://116.204.107.117:8000/media/headshot/" + header.Filename
	err = global.DB.Save(user).Error
	if err != nil {
		c.JSON(403, gin.H{
			"status": 403,
			"msg":    "保存文件路径到数据库中失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "修改用户头像成功",
		"data":   user,
	})
	return
}
