package v1

import (
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"github.com/gin-gonic/gin"
)

// AddUserConcept
// @Summary     txc
// @Description 添加user的关注关键词
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param       data    body     response.AddUserConceptQ true "data"
// @Param       x-token header   string                   true "token"
// @Success     200     {string} json                     "{"msg":"添加成功"}"
// @Failure     400     {string} json                     "{"msg":"参数错误"}"
// @Failure     401     {string} json                     "{"msg":"用户不存在"}"
// @Failure     402     {string} json                     "{"msg":"concept不存在"}"
// @Failure     403     {string} json                     "{"msg":"添加失败"}"
// @Failure     404     {string} json                     "{"msg":"删除失败"}"
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
	if _, err := service.GetObject("concepts", d.ConceptID); err != nil {
		c.JSON(402, gin.H{"msg": "concept不存在"})
		return
	}
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
