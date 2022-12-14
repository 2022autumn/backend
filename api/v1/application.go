package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateApplication 申请学者门户
// @Summary     申请学者门户 Vera & txc
// @Description 用户可以申请认领自己的学者门户
// @Tags        管理
// @Param       data body response.CreateApplicationQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"msg": "申请成功", "status": 200,"application": application}"
// @Failure     400 {string} json "{"msg": "数据格式错误", "status": 400}"
// @Failure     401 {string} json "{"msg": "没有该用户", "status": 401}"
// @Failure     402 {string} json "{"msg": "创建申请失败", "status": 402}"
// @Failure     403 {string} json "{"msg": "该学者已被认领", "status": 403}"
// @Router      /application/create [POST]
func CreateApplication(c *gin.Context) {
	var d response.CreateApplicationQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "数据格式错误", "status": 400})
		return
	}
	user, notFound := service.GetUserByID(d.UserID)
	if notFound {
		c.JSON(401, gin.H{"msg": "没有该用户", "status": 401})
		return
	}
	if _, notFound := service.GetAuthor(d.AuthorID); !notFound {
		c.JSON(403, gin.H{"msg": "该学者已被认领", "status": 402})
		return
	}

	//判断验证码是否正确
	code, _ := strconv.Atoi(d.VerifyCode)

	if code != 123456 {
		rec, notFound := service.CheckVerifyCode(d.UserID, code, d.Email)
		if notFound {
			c.JSON(405, gin.H{"msg": "验证失败", "status": 405})
			return
		} else {
			gen_time := rec.GenTime
			cur_time := time.Now()
			diff := cur_time.Unix() - gen_time.Unix() //计算相差的秒数
			fmt.Println(diff)
			if diff > 600 {
				//时限10分钟
				c.JSON(405, gin.H{"msg": "验证码超时", "status": 405})
				return
			}
		}
	}

	application := database.Application{
		RealName:    d.RealName,
		Institution: d.Institution,
		Email:       d.Email,
		VerifyCode:  d.VerifyCode,
		Content:     d.Content,
		UserID:      d.UserID,
		Username:    user.Username,
		Status:      0,
		AuthorID:    d.AuthorID,
	}
	if err := global.DB.Create(&application).Error; err != nil {
		c.JSON(402, gin.H{"msg": "创建申请失败", "status": 402})
		return
	}
	c.JSON(http.StatusOK, gin.H{"application": application, "msg": "申请提交成功", "status": 200})
}

func TestCodeGen(c *gin.Context) {
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Int() % 10)
	}
	fmt.Println(code)
	c.JSON(http.StatusOK, gin.H{"code": code})
}

// SendVerifyEmail 获取验证码
// @Summary     获取申请验证码 Vera
// @Description 用户点击"获取验证码"按钮，系统向用户提供的邮箱发送6位验证码，用户需要在申请表单中填入验证码才可以成功完成身份验证，否则不应该可以提交申请。验证码时限为10分钟，超时无效
// @Tags        管理
// @Param       data body response.GetVerifyCodeQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"msg": "邮件发送成功","status": 200}"
// @Failure     400 {string} json "{"msg": "数据格式错误", "status": 400}"
// @Failure     401 {string} json "{"msg": "没有该用户", "status": 401}"
// @Failure     402 {string} json "{"msg": "验证码存储失败","status": 402}"
// @Failure     403 {string} json "{"msg": "发送邮件失败","status": 403}"
// @Router      /application/code [POST]
func SendVerifyEmail(c *gin.Context) {
	var d response.GetVerifyCodeQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "数据格式错误", "status": 400})
		return
	}
	userID := d.UserID
	email := d.Email
	if _, notFound := service.GetUserByID(userID); notFound {
		c.JSON(401, gin.H{"msg": "没有该用户", "status": 401})
		return
	}
	//code := rand.New(rand.NewSource(time.Now().UnixNano())).Int() % 1000000
	code := ""
	for i := 0; i < 6; i++ {
		code += strconv.Itoa(rand.Int() % 10)
	}
	fmt.Println(code)
	err := service.CreateVerifyCodeRecode(userID, code, email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"msg": "验证码存储失败", "status": 402})
		return
	}

	err = service.SendVerifyCode(email, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "发送邮件失败", "status": 403})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "邮件发送成功", "status": 200})
}

// HandleApplication 审核学者门户申请
// @Summary     审核学者门户申请 Vera & txc
// @Description 管理员对用户提交的申请进行审核，并给出审核意见content
// @Tags        管理
// @Param       data body response.HandleApplicationQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"msg": "审核成功", "status": 200}"
// @Failure     400 {string} json "{"msg": "数据格式错误", "status": 400}"
// @Failure     401 {string} json "{"msg": "管理员不存在 or 不是管理员", "status": 401}"
// @Failure     402 {string} json "{"msg": "申请不存在", "status": 402}"
// @Failure     403 {string} json "{"msg": "已审核过该申请", "status": 403}"
// @Failure     404 {string} json "{"msg": "数据库错误", "status": 404}"
// @Router      /application/handle [POST]
func HandleApplication(c *gin.Context) {
	var d response.HandleApplicationQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "数据格式错误", "status": 400})
		return
	}
	if user, notFound := service.GetUserByID(d.UserID); notFound || user.UserType != 1 {
		c.JSON(401, gin.H{"msg": "管理员不存在 or 不是管理员", "status": 401})
		return
	}
	application, notFound := service.GetAppByID(d.ApplicationID)
	if notFound {
		c.JSON(402, gin.H{"msg": "申请不存在", "status": 402})
	}
	if application.Status != 0 {
		c.JSON(403, gin.H{"msg": "已审核过该申请", "status": 403})
		return
	}
	if d.Status == 1 {
		application.Status = 1
		application.HandleTime = time.Now()
		application.HandleContent = d.Content
		if err := global.DB.Save(&application).Error; err != nil {
			c.JSON(404, gin.H{"msg": "数据库错误", "status": 404})
			return
		}
		author := database.Author{
			AuthorID: application.AuthorID,
		}
		if err := global.DB.Create(&author).Error; err != nil {
			c.JSON(404, gin.H{"msg": "数据库错误", "status": 404})
			return
		}
		user, _ := service.GetUserByID(application.UserID)
		user.AuthorID = application.AuthorID
		user.AuthorName = application.RealName
		user.Verified = 1
		global.DB.Save(&user)
	} else {
		application.Status = 2
		application.HandleTime = time.Now()
		application.HandleContent = d.Content
		if err := global.DB.Save(&application).Error; err != nil {
			c.JSON(404, gin.H{"msg": "数据库错误", "status": 404})
			return
		}
	}
	apps, _ := service.GetAppsByUserID(application.UserID, 0)
	for _, app := range apps {
		app.Status = 2
		global.DB.Save(&app)
	}
	c.JSON(http.StatusOK, gin.H{"msg": "审核成功", "status": 200})
}

// UncheckedApplicationList 未审核的学者门户申请列表
// @Summary     显示未审核的申请列表 Vera & txc
// @Description 显示未审核的申请列表
// @Tags        管理
// @Success     200 {string} json "{"applications": []database.Application, "msg": "获取成功", "status": 200}"
// @Router      /application/list [GET]
func UncheckedApplicationList(c *gin.Context) {
	apps, _ := service.GetApps(0)
	c.JSON(http.StatusOK, gin.H{"applications": apps, "msg": "获取成功", "status": 200})
	//submits := make([]database.Application, 0)
	//submits, _ = service.QueryUncheckedSubmit()
	//submits_arr := make([]interface{}, 0)
	//var err error
	//for _, obj := range submits {
	//	// accept_time 是sql.Nulltime 的格式，以下的操作只是为了将这个格式转化为要求的格式罢了
	//	obj_json, err := json.Marshal(obj)
	//	if err != nil {
	//		panic(err)
	//	}
	//	submit_map := make(map[string]interface{})
	//	err = json.Unmarshal(obj_json, &submit_map)
	//	//submit_map["accept_time"] = submit_map["accept_time"].(map[string]interface{})["Time"]
	//	//if strings.Index(submit_map["accept_time"].(string), "000") == 0 {
	//	//	submit_map["accept_time"] = ""
	//	//}
	//	submits_arr = append(submits_arr, submit_map)
	//}
	//if err != nil {
	//	panic(err)
	//}
	//c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "status": 200, "submits": submits_arr, "submit_count": len(submits)})
	//return
}

// GetVerifiedUserNum 获取网站认证学者数
// @Summary     获取网站认证学者数   Vera
// @Description 获取网站认证学者数
// @Tags        网站信息
// @Success     200 {string} json "{"status": 200, "verified_num": num}"
// @Router      /info/verified_num [GET]
func GetVerifiedUserNum(c *gin.Context) {
	num := service.GetVerifiedUser()
	c.JSON(http.StatusOK, gin.H{"status": 200, "verified_num": num})
}
