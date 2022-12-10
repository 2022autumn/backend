package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateApplication 申请学者门户
// @Summary     Vera
// @Description 用户可以申请认领自己的学者门户
// @Tags        管理
// @Param       data body     response.CreateApplicationQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"success": true, "application_id": submit.ApplicationID, "message": "申请提交成功", "status": 200}"
// @Failure		401 {string} json "{"success": false, "message": "申请创建失败", "status": 401}"
// @Failure     404 {string} json "{"success": false, "message": "没有该用户", "status": 404}"
// @Failure     405 {string} json "{"success": false, "message": "该作者已被认领", "status": 405, "the_authorname": the_submit.AuthorName}"
// @Failure     406 {string} json "{"success": false, "message": "您已经是认证学者，请勿重复申请", "status": 406}"
// @Router      /application/create [POST]
func CreateApplication(c *gin.Context) {
	var d response.CreateApplicationQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	author_name := d.AuthorName
	instituition_name := d.InstitutionName
	work_email := d.WorkEmail
	//field := d.Field
	//home_page := c.Request.FormValue("home_page")
	author_id := d.AuthorID
	user_id := d.UserID

	_, notFound := service.GetUserByID(user_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户", "status": 404})
		return
	}
	if the_submit, notFound := service.QueryApplicationByAuthor(author_id); !notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "该作者已被认领", "status": 405, "the_authorname": the_submit.AuthorName})
		return
	}
	if _, notFound := service.QueryUserIsScholar(user_id); !notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "您已经是认证学者，请勿重复申请", "status": 406})
		return
	}
	//后续对papers可能需要处理
	//papers := service.GetAuthorAllPaper(author_id)
	//author := service.GetAuthors(append(make([]string, 0), author_id))[0].(map[string]interface{})
	submit := database.Application{UserID: user_id,
		InstitutionName: instituition_name, AuthorName: author_name, Email: work_email,
		//HomePage: home_page,
		AuthorID: author_id,
		//Fields:   field,
		Status: 0, Content: "",
		//PaperCount: int(author["paper_count"].(float64)),
		ApplyTime: time.Now()}

	err := service.CreateApplication(&submit)
	if err != nil {
		panic(err)
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "申请创建失败", "status": 401})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "application_id": submit.ApplicationID, "message": "申请提交成功", "status": 200}) //, "papers": service.GetAuthorSomePapers(author_id, 100)})
}

// HandleApplication 审核学者门户申请
// @Summary     Vera
// @Description 管理员对用户提交的申请进行审核，并给出审核意见content
// @Tags        管理
// @Param       data body     response.HandleApplicationQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"success": true, "message": "申请审批成功", "status": 200}"
// @Failure		401 {string} json "{"success": false, "message": "申请不存在", "status": 404}"
// @Failure     405 {string} json "{"success": false, "message": "没有该用户", "status": 405}"
// @Failure     406 {string} json "{"success": false, "message": "已审核过该申请", "status": 406}"
// @Failure     403 {string} json "{"success": false, "message": "success 不为true false", "status": 403}"
// @Router      /application/handle [POST]
func HandleApplication(c *gin.Context) {
	var d response.HandleApplicationQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	application_id := d.ApplicationID
	userID := d.UserID
	success := d.HandleRes
	content := d.HandleRes
	user, notFound := service.GetUserByID(userID)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "没有该用户", "status": 405})
		return
	}
	application, notFound := service.GetApplicationByID(application_id)
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
	err := global.DB.Save(application).Error
	fmt.Println(application.HandleTime)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "申请审批成功", "status": 200})
}

// UncheckedApplicationList 未审核的学者门户申请列表
// @Summary     Vera
// @Description 显示未审核的申请列表
// @Tags        管理
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"success": true, "message": "获取成功", "status": 200, "submits": submits_arr, "submit_count": len(submits)}"
// @Router      /application/list [POST]
func UncheckedApplicationList(c *gin.Context) {
	submits := make([]database.Application, 0)
	submits, _ = service.QueryUncheckedSubmit()
	submits_arr := make([]interface{}, 0)
	var err error
	for _, obj := range submits {
		// accept_time 是sql.Nulltime 的格式，以下的操作只是为了将这个格式转化为要求的格式罢了
		obj_json, err := json.Marshal(obj)
		if err != nil {
			panic(err)
		}
		submit_map := make(map[string]interface{})
		err = json.Unmarshal(obj_json, &submit_map)
		//submit_map["accept_time"] = submit_map["accept_time"].(map[string]interface{})["Time"]
		//if strings.Index(submit_map["accept_time"].(string), "000") == 0 {
		//	submit_map["accept_time"] = ""
		//}
		submits_arr = append(submits_arr, submit_map)
	}
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "获取成功", "status": 200, "submits": submits_arr, "submit_count": len(submits)})
	return
}
