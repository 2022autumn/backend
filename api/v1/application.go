package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/service"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

func CreateApplication(c *gin.Context) {
	// 创建申请称为学者的表，等待平台管理申请
	author_name := c.Request.FormValue("author_name")
	instituition_name := c.Request.FormValue("institution_name")
	work_email := c.Request.FormValue("work_email")
	//fields := c.Request.FormValue("fields")
	//home_page := c.Request.FormValue("home_page")
	author_id := c.Request.FormValue("author_id")
	user_id := c.Request.FormValue("user_id")
	user_id_u64, err := strconv.ParseUint(user_id, 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户ID不为正整数", "status": 402})
		return
	}
	_, notFound := service.GetUserByID(user_id_u64)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户", "status": 404})
		return
	}
	if the_submit, notFound := service.QueryApplicationByAuthor(author_id); !notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "该作者已被认领", "status": 405, "the_authorname": the_submit.AuthorName})
		return
	}
	if _, notFound := service.QueryUserIsScholar(user_id_u64); !notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "您已经是认证学者，请勿重复申请", "status": 406})
		return
	}
	//后续对papers可能需要处理
	//papers := service.GetAuthorAllPaper(author_id)
	//author := service.GetAuthors(append(make([]string, 0), author_id))[0].(map[string]interface{})
	submit := database.Application{UserID: user_id_u64,
		InstitutionName: instituition_name, AuthorName: author_name, Email: work_email,
		//HomePage: home_page,
		AuthorID: author_id,
		//Fields: fields,
		Status: 0, Content: "",
		//PaperCount: int(author["paper_count"].(float64)),
		ApplyTime: time.Now()}

	err = service.CreateApplication(&submit)
	if err != nil {
		panic(err)
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "申请创建失败", "status": 401})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请提交成功", "status": 200}) //, "papers": service.GetAuthorSomePapers(author_id, 100)})
	return
}

func HandleApplication(c *gin.Context) {
	application_id := c.Request.FormValue("application_id")
	application_id_u64, err1 := strconv.ParseUint(application_id, 10, 64)
	user_id := c.Request.FormValue("user_id")
	success := c.Request.FormValue("success")
	content := c.Request.FormValue("content")
	if err1 != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "申请表ID不为正整数", "status": 402})
		return
	}
	user_id_u64, err := strconv.ParseUint(user_id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "用户ID不为正整数", "status": 401})
		return
	}
	user, notFound := service.GetUserByID(user_id_u64)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户", "status": 405})
		return
	}
	application, notFound := service.GetApplicationByID(application_id_u64)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "申请不存在", "status": 404})
	}
	fmt.Println("check user application", user.UserID)
	if application.Status != 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "已审核过该申请", "status": 406})
		return
	}
	if success == "false" {
		application.Status = 2
		application.Content = content
	} else if success == "true" {
		application.Status = 1
		application.Content = content
		service.MakeUserScholar(user, application)
		application.HandleTime = sql.NullTime{Time: time.Now(), Valid: true}
	} else {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "success 不为true false", "status": 403})
		return
	}
	err = global.DB.Save(application).Error
	fmt.Println(application.HandleTime)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请审批成功", "status": 200})
	return
}
