package v1

import (
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTag 创建收藏夹
// @Summary     用户可以按照需要建立收藏夹 Vera
// @Description 用户可以按照需要建立收藏夹
// @Tags        社交
// @Param       data body response.TagCreation true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status": 200, "msg": "收藏夹创建成功", "tag_id": tag.TagID}"
// @Failure     400 {string} json "{"status": 400, "msg":"用户ID不存在"}"
// @Failure     401 {string} json "{"status": 401, "msg":"收藏夹已存在，换个名字吧～"}"
// @Failure     402 {string} json "{"status": 402, "msg": "创建失败"}"
// @Router      /social/tag/create [POST]
func CreateTag(c *gin.Context) {
	var d response.TagCreation
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	tag_name := d.TagName

	//验证用户是否存在
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}

	//nameExisted := service.QueryUserTagByName(user_id, tag_name)
	//if nameExisted {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": 401,
	//		"msg":    "收藏夹已存在，换个名字吧～",
	//	})
	//	return
	//}
	tag, notFoundTag := service.QueryATag(user_id, tag_name)
	if !notFoundTag {
		c.JSON(http.StatusOK, gin.H{
			//"success": false,
			"status": 401,
			"msg":    "收藏夹已存在，换个名字吧～",
		})
		return
	}
	tag = database.Tag{UserID: user_id, TagName: tag_name, CreateTime: time.Now()}
	err := service.CreateUserTag(&tag)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": 402, "msg": "创建失败"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "收藏夹创建成功", "tag_id": tag.TagID})
		return
	}
}

// AddTagToPaper 收藏文献
// @Summary     将某篇文献加入到某一收藏夹下 Vera
// @Description 将某篇文献加入到某一收藏夹下
// @Tags        社交
// @Param       data body response.AddTagToPaper true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status": 200, "msg": "收藏成功"}"
// @Failure     400 {string} json "{"status": 400, "msg":"用户ID不存在"}"
// @Failure     401 {string} json "{"status": 401, "msg": "用户无此收藏夹"}"
// @Failure     402 {string} json "{"status": 402, "msg": "文章已在此收藏夹下"}"
// @Failure     403 {string} json "{"status": 403, "msg": "收藏失败"}"
// @Router      /social/tag/collectPaper [POST]
func AddTagToPaper(c *gin.Context) {
	var d response.AddTagToPaper
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	paper_id := d.PaperID
	tag_id := d.TagID
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	_, notFound := service.QueryTagUser(tag_id, user_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"status": 401, "msg": "用户无此收藏夹"})
		return
	}
	_, notFound = service.QueryPaperTag(tag_id, paper_id)
	if !notFound {
		c.JSON(http.StatusOK, gin.H{"status": 402, "msg": "文章已在此收藏夹下"})
		return
	} else {
		tag_paper := database.TagPaper{
			PaperID:    paper_id,
			TagID:      tag_id,
			CreateTime: time.Now(),
		}
		err := service.CreateTagPaper(&tag_paper)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 403, "msg": "收藏失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "收藏成功"})
		return
	}
}

// RemovePaperTag 将某篇文献取消收藏
// @Summary     Vera
// @Description 将某篇文献从收藏夹中移除
// @Tags        社交
// @Param       data body response.AddTagToPaper true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status": 200, "msg": "取消收藏成功"}"
// @Failure     400 {string} json "{"status": 400,"msg":"用户ID不存在"}"
// @Failure     401 {string} json "{"status": 401, "msg": "用户无此收藏夹"}"
// @Failure     402 {string} json "{"status": 402, "msg": "文章不在此收藏夹下"}"
// @Failure     403 {string} json "{"status": 403, "msg": "取消收藏失败"}"
// @Router      /social/tag/cancelCollectPaper [POST]
func RemovePaperTag(c *gin.Context) {
	var d response.AddTagToPaper
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	paper_id := d.PaperID
	tag_id := d.TagID
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	_, notFound := service.QueryTagUser(tag_id, user_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"status": 401, "msg": "用户无此收藏夹"})
		return
	}
	_, notFound = service.QueryPaperTag(tag_id, paper_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"status": 402, "msg": "文章不在此收藏夹下"})
		return
	} else {
		err := service.DeleteTagPaper(tag_id, paper_id)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"status": 403, "msg": "取消收藏失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "取消收藏成功"})
	}
}

// ShowTagPaperList 查看收藏夹内的文献列表
// @Summary     返回某一收藏夹内的文献信息 Vera
// @Description 返回某一收藏夹内的文献信息
// @Tags        社交
// @Param       data body response.TagPaperListQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"success": true, "status": 200,"num":int, "paper_list": paper_list,"msg": "查询成功"}"
// @Success     210 {string} json "{"success": true, "status": 402,"num":0, "msg": "标签下没有文章"}"
// @Failure     400 {string} json "{"success": false,"status": 400, "msg":"用户ID不存在"}"
// @Failure     401 {string} json "{"success": false,"status": 401, "msg": "用户无此收藏夹"}"
// @Failure     404 {string} json "{"success": false, "status": 404, "msg":"查询失败"}"
// @Failure     403 {string} json "{"status": 403, "msg": "收藏失败"}"
// @Router      /social/tag/sublist [POST]
func ShowTagPaperList(c *gin.Context) {
	var d response.TagPaperListQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	tag_id := d.TagID
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  400,
			"msg":     "用户ID不存在",
		})
		return
	}
	tag, notFound := service.QueryTagUser(tag_id, user_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 401, "msg": "用户无此收藏夹"})
		return
	}

	var paper_ids []string
	papers := service.QueryTagPaper(tag.TagID)
	if papers == nil || len(papers) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"status":  210,
			"num":     0,
			"msg":     "标签下没有文章",
		})
		return
	}
	for _, paper := range papers {
		paper_ids = append(paper_ids, paper.PaperID)
	}
	var paper_list []interface{}
	for _, id := range paper_ids {
		//fmt.Println(id)
		idx, err := utils.TransObjPrefix(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": 400,
				"msg":    "id type error",
			})
			return
		}
		//res, err := service.GetObject("works", id)
		res, err, _ := service.GetObject2(idx, id)
		//fmt.Println(res)
		if err != nil {
			c.JSON(404, gin.H{
				"success": false,
				"status":  404,
				"msg":     "查询失败",
			})
			return
		}
		var tmp = make(map[string]interface{})
		tmp = res
		tmp["star_num"], _ = service.GetPaperStarNum(id)
		paper_list = append(paper_list, tmp)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "num": len(paper_ids), "paper_list": paper_list, "msg": "查询成功"})
}

// ShowUserTagList 用户收藏夹列表
// @Summary     显示用户建立的所有收藏夹 Vera
// @Description 显示用户建立的所有收藏夹
// @Tags        社交
// @Param       data body response.UserInfo true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"success": true, "status":  200, "msg": "查看收藏夹列表成功", "data":tags}"
// @Failure     400 {string} json "{"success": false,"status": 400, "msg":"用户ID不存在"}"
// @Failure     403 {string} json "{"success": false,"status": 403, "msg":"未查询到结果"}"
// @Router      /social/tag/taglist [POST]
func ShowUserTagList(c *gin.Context) {
	var d response.UserInfo
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	//验证用户是否存在
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  400,
			"msg":     "用户ID不存在",
		})
		return
	}
	tags := service.QueryTagList(user_id)
	if tags == nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"msg":     "未查询到结果",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true, "status": 200, "msg": "查看收藏夹列表成功", "data": tags,
	})
}

// DeleteTag 删除标签
// @Summary     删除标签 Vera
// @description 删除标签
// @Tags        社交
// @Param       data body response.TagPaperListQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} string "{"success": true,"status":200, "msg": "标签删除成功"}"
// @Failure     400 {string} string "{"success": false,"status":400, "msg": "用户ID不存在"}"
// @Failure     403 {string} string "{"success": false,"status":403, "msg": "标签不存在"}"
// @Router      /social/tag/delete [POST]
func DeleteTag(c *gin.Context) {
	var d response.TagPaperListQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	tag_id := d.TagID

	//验证用户是否存在
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"status":  400,
			"msg":     "用户ID不存在",
		})
		return
	}
	tag, notFoundTag := service.GetTagById(tag_id)

	if notFoundTag {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 403, "msg": "标签不存在"})
		return
	}
	//tagPapers := service.QueryTagPaper(tag.TagID)
	service.DeleteTag(tag.TagID)
	//for _, paper := range tagPapers {
	//	collect, _ := service.QueryACollect(user_id, paper.PaperID)
	//	collect.TagCount--
	//	if collect.TagCount == 0 {
	//		service.DeleteACollect(collect.ID)
	//	} else {
	//		service.UpdateACollect(&collect)
	//	}
	//	service.DeleteTagPaper(paper.ID)
	//}
	c.JSON(http.StatusOK, gin.H{"success": true, "status": 200, "msg": "标签删除成功"})
}

// RenameTag 收藏夹重命名
// @Summary     收藏夹重命名 Vera
// @description 对原有的收藏夹重命名
// @Tags        社交
// @Param       data body response.RenameTagQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"msg":"修改成功"}"
// @Failure     400 {string} json "{"status": 400, "msg": "用户ID不存在"}"
// @Failure     401 {string} json "{"status": 401, "msg": "用户无此收藏夹"}"
// @Failure     402 {string} json "{"status": 402, "msg": "名称已存在，换个名字吧～"}"
// @Failure     403 {string} json "{"status":403,"msg":"修改失败"}"
// @Router      /social/tag/rename [POST]
func RenameTag(c *gin.Context) {
	var d response.RenameTagQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	tag_id := d.TagID
	new_name := d.NewTagName

	//验证用户是否存在
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	tag, notFound := service.QueryTagUser(tag_id, user_id)
	if notFound {
		c.JSON(http.StatusOK, gin.H{"status": 401, "msg": "用户无此收藏夹"})
		return
	}
	same_name_tag, notFoundTag := service.QueryATag(user_id, new_name)
	if !notFoundTag && same_name_tag != tag {
		c.JSON(http.StatusOK, gin.H{
			//"success": false,
			"status": 402,
			"msg":    "名称已存在，换个名字吧～",
		})
		return
	}
	err := service.RenameTag(new_name, tag)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 403, "msg": "修改失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "修改成功"})
}

// PaperTagList 显示某篇文章从属的收藏夹信息
// @Summary     查询文章所在的收藏夹 Vera
// @description 返回用户所有的收藏夹列表，并一一标记文章是否从属于收藏夹
// @Tags        社交
// @Param       data body response.PaperBelongingQ true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"status":200,"msg":"获取成功","data":data}"
// @Failure     400 {string} json "{"status": 400, "msg": "用户ID不存在"}"
// @Failure     401 {string} json "{"status": 401, "msg": "未查询到结果"}"
// @Router      /social/tag/paperTagList [POST]
func PaperTagList(c *gin.Context) {
	var d response.PaperBelongingQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	user_id := d.UserID
	paper_id := d.PaperID

	//验证用户是否存在
	_, notFoundUserByID := service.QueryAUserByID(user_id)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	tags := service.QueryTagList(user_id)
	if tags == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 401,
			"msg":    "未查询到结果",
		})
		return
	}

	var tag_list []interface{}
	for _, value := range tags {
		tag_id := value.TagID
		var tag_info = make(map[string]interface{})
		//tag_info["create_time"]
		tag_info["create_time"] = value.CreateTime
		tag_info["tag_id"] = tag_id
		tag_info["tag_name"] = value.TagName
		tag_info["user_id"] = value.UserID
		_, isCollect := service.QueryPaperIsCollect(paper_id, tag_id)
		tag_info["isCollect"] = isCollect
		tag_list = append(tag_list, tag_info)
	}
	c.JSON(http.StatusOK, gin.H{"status": 200, "msg": "获取成功", "data": tag_list})
}
