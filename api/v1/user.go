package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register 注册
// @Summary     注册 ccf
// @Description 注册
// @Description 填入用户名和密码注册
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       data body     response.RegisterQ true "data"
// @Success     200  {string} json               "{"status":200,"msg":"注册成功"}"
// @Failure     400  {string} json               "{"status":400,"msg":"用户名已存在"}"
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
		UserInfo: "这个用户很懒什么都没有留下",
	}
	// 成功创建用户
	if err := service.CreateUser(&user); err != nil {
		panic("CreateUser: create user error")
	}

	//为用户创建默认收藏夹
	tag := database.Tag{UserID: user.UserID, TagName: "默认收藏夹", CreateTime: time.Now()}
	_ = service.CreateUserTag(&tag)

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "注册成功",
	})
}

// Login 登录
// @Summary     登录 ccf
// @Description 登录
// @Description 填入用户名和密码
// @Tags        用户
// @Accept      json
// @Produce     json
// @Param       data body     response.LoginQ true "data"
// @Success     200  {string} json            "{"status":200,"msg":"登录成功","token": token,"ID": user.UserID}"
// @Failure     400  {string} json            "{"status":400,"msg":"用户名不存在"}"
// @Failure     401  {string} json            "{"status":401,"msg":"密码错误"}"
// @Router      /login [POST]
func Login(c *gin.Context) {
	var d response.LoginQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	// 用户不存在
	user, notFound := service.GetUserByUsername(d.Username)
	if notFound {
		c.JSON(400, gin.H{
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
	//token := utils.GenerateToken(user.UserID)
	token := user.UserID
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"msg":    "登录成功",
		"token":  token,
		"ID":     user.UserID,
		"user":   user,
	})
}

// UserInfo 查看用户个人信息
// @Summary     查看用户个人信息 ccf
// @Description 查看用户个人信息
// @Tags        用户
// @Param       user_id query string true "user_id"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"msg":"获取用户信息成功","data":{object}}"
// @Failure     400      {string} json   "{"status":400,"msg":"用户ID不存在"}"
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
// @Summary     编辑用户信息 ccf
// @Description 编辑用户信息
// @Tags        用户
// @Param       data body response.ModifyQ true "data"
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
// @Summary     编辑用户密码 ccf
// @Description 编辑用户密码
// @Tags        用户
// @Param       data body response.PwdModifyQ true "data"
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
// @Summary     上传用户头像 ccf
// @Description 上传用户头像
// @Tags        用户
// @Param       user_id  formData string true "用户ID"
// @Param       Headshot formData file   true "新头像"
// @Success     200      {string} json   "{"status":200,"msg":"修改成功","data":{object}}"
// @Failure     400 {string} json "{"status":400,"msg":"用户ID不存在"}"
// @Failure     401      {string} json   "{"status":401,"msg":"头像文件上传失败"}"
// @Failure     402      {string} json   "{"status":402,"msg":"文件保存失败"}"
// @Failure     403      {string} json   "{"status":403,"msg":"保存文件路径到数据库中失败"}"
// @Router      /user/headshot [POST]
func UploadHeadshot(c *gin.Context) {
	userId, _ := strconv.ParseUint(c.Request.FormValue("user_id"), 0, 64)
	user, notFoundUserByID := service.QueryAUserByID(userId)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	//1、获取上传的文件
	file, err := c.FormFile("Headshot")
	if err != nil {
		c.JSON(401, gin.H{"msg": "头像文件上传失败"})
		return
	}
	raw := fmt.Sprintf("%d", userId) + time.Now().String() + file.Filename
	md5 := utils.GetMd5(raw)
	suffix := strings.Split(file.Filename, ".")[1]
	saveDir := "./media/headshot/"
	saveName := md5 + "." + suffix
	savePath := path.Join(saveDir, saveName)
	if err = c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(402, gin.H{"msg": "文件保存失败"})
		return
	}
	//3、将文件对应路径更新到数据库中
	//user.HeadShot = "http://116.204.107.117:8000/api/media/headshot/" + saveName
	user.HeadShot = saveName
	err = global.DB.Save(user).Error
	if err != nil {
		c.JSON(403, gin.H{"msg": "保存文件路径到数据库中失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "修改用户头像成功", "data": user})
}

// GetRegisterUserNum 获取网站所有注册用户数量
// @Summary     获取网站所有注册用户数量 Vera
// @Description 统计网站信息，该接口不需要前端参数
// @Tags        网站信息
// @Success     200 {string} json "{"status": 200, "register_num": int}"
// @Router      /info/register_num [GET]
func GetRegisterUserNum(c *gin.Context) {
	num := service.GetAllUser()
	c.JSON(http.StatusOK, gin.H{"status": 200, "register_num": num})
}

// GetBrowseHistory
// @Summary     txc
// @Description 获取用户浏览历史
// @Tags        用户
// @Param       token header   string                     true "token"
// @Param       data  body     response.GetBrowseHistoryQ true "data"
// @Success     200   {string} json                       "{"status": 200, "msg": "获取成功", "data": {object}}"
// @Router      /user/history [POST]
func GetBrowseHistory(c *gin.Context) {
	user := c.MustGet("user").(database.User)
	var d response.GetBrowseHistoryQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 400, "msg": "参数错误"})
		return
	}
	var history []database.BrowseHistory
	global.DB.Where("user_id = ?", user.UserID).Order("browse_time desc").Find(&history)
	count := len(history)
	if count < d.Page*d.Size {
		history = history[(d.Page-1)*d.Size:]
	} else {
		history = history[(d.Page-1)*d.Size : d.Page*d.Size]
	}
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "获取成功", "data": history, "count": count})
}
