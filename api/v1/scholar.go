package v1

import (
	"IShare/model/database"
	"IShare/service"
	"github.com/gin-gonic/gin"
)

// AddUserConcept
// @Summary     txc
// @Description 添加user的关注关键词
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param       data body     response.AddUserConceptQ true "data"
// @Success     200  {string} json                     "{"msg":"添加成功"}"
// @Failure     401  {string} json                     "{"msg":"用户不存在"}"
// @Failure     402  {string} json                     "{"msg":"concept不存在"}"
// @Failure     403  {string} json                     "{"msg":"添加失败"}"
// @Failure     404  {string} json                     "{"msg":"删除失败"}"
// @Router      /scholar/concept [POST]
func AddUserConcept(c *gin.Context) {
	var d database.UserConcept
	_ = c.ShouldBindJSON(&d)
	if _, notFound := service.GetUserByID(d.UserID); notFound {
		c.JSON(401, gin.H{"msg": "用户不存在"})
		return
	}
	if _, err := service.GetObject("concepts", d.Concept); err != nil {
		c.JSON(402, gin.H{"msg": "concept不存在"})
		return
	}
	userConcept, notFound := service.GetUserConcept(d.UserID, d.Concept)
	if notFound {
		userConcept = database.UserConcept{
			UserID:  d.UserID,
			Concept: d.Concept,
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
