package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
	"time"
)

// AddUserConcept
// @Summary     txc
// @Description 添加user的关注关键词
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param       data  body     response.AddUserConceptQ true "data"
// @Param       token header   string                   true "token"
// @Success     200   {string} json                     "{"msg":"添加成功"}"
// @Failure     400   {string} json                     "{"msg":"参数错误"}"
// @Failure     401   {string} json                     "{"msg":"用户不存在"}"
// @Failure     402   {string} json                     "{"msg":"concept不存在"}"
// @Failure     403   {string} json                     "{"msg":"添加失败"}"
// @Failure     404   {string} json                     "{"msg":"删除失败"}"
// @Router      /scholar/concept [POST]
func AddUserConcept(c *gin.Context) {
	var d response.AddUserConceptQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
	}
	if _, notFound := service.GetUserByID(d.UserID); notFound {
		c.JSON(401, gin.H{"msg": "用户不存在"})
		return
	}
	index, err := utils.TransObjPrefix(d.ConceptID)
	if err != nil || index != "concepts" {
		c.JSON(402, gin.H{"msg": "concept参数错误"})
		return
	}
	//if _, err := service.GetObject("concepts", d.ConceptID); err != nil {
	//	c.JSON(402, gin.H{"msg": "concept不存在"})
	//	return
	//}
	userConcept, notFound := service.GetUserConcept(d.UserID, d.ConceptID)
	if notFound {
		userConcept = database.UserConcept{
			UserID:    d.UserID,
			ConceptID: d.ConceptID,
		}
		if err := service.CreateUserConcept(&userConcept); err != nil {
			c.JSON(403, gin.H{"msg": "添加失败"})
			return
		}
		c.JSON(200, gin.H{"msg": "添加成功"})
		return
	}
	if err := service.DeleteUserConcept(&userConcept); err != nil {
		c.JSON(404, gin.H{"msg": "删除失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "删除成功"})
}

// RollWorks
// @Summary     txc
// @Description 获取用户推荐的文章 请勿使用
// @Tags        scholar
// @Success     200 {string} json "{"msg":"获取成功","data":{}}"
// @Router      /scholar/roll [GET]
func RollWorks(c *gin.Context) {

	c.JSON(200, gin.H{"msg": "error"})
}

// GetHotWorks
// @Summary     txc
// @Description 获取热门文章（根据访问量）
// @Tags        scholar
// @Success     200 {string} json "{"msg":"获取成功","data":{}}"
// @Failure     400 {string} json "{"msg":"获取失败"}"
// @Router      /scholar/hot [GET]
func GetHotWorks(c *gin.Context) {
	works, err := service.GetHotWorks(10)
	if err != nil {
		c.JSON(400, gin.H{"msg": "获取失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "获取成功", "data": works})
}

// UploadAuthorHeadshot
// @Summary     txc
// @Description 上传作者头像
// @Tags        scholar
// @Param       author_id formData string true "用户ID"
// @Param       Headshot  formData file   true "新头像"
// @Router      /scholar/author/headshot [POST]
func UploadAuthorHeadshot(c *gin.Context) {
	authorID := c.Request.FormValue("author_id")
	author, notFound := service.GetAuthor(authorID)
	if notFound {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "学者未被认领"})
		return
	}
	file, err := c.FormFile("Headshot")
	if err != nil {
		c.JSON(401, gin.H{"msg": "头像文件获取失败"})
		return
	}
	raw := fmt.Sprintf("%d", authorID) + time.Now().String() + file.Filename
	md5 := utils.GetMd5(raw)
	suffix := strings.Split(file.Filename, ".")[1]
	saveDir := "./media/headshot/"
	saveName := md5 + "." + suffix
	savePath := path.Join(saveDir, saveName)
	if err = c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(402, gin.H{"msg": "文件保存失败"})
		return
	}
	author.HeadShot = saveName
	err = global.DB.Save(author).Error
	if err != nil {
		c.JSON(403, gin.H{"msg": "保存文件路径到数据库中失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "修改用户头像成功", "data": author})
	return
}
