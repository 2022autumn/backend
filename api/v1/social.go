package v1

import (
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateComment 创建评论
// @Summary     用户可以在某一篇文献的评论区中发表自己的评论 Vera
// @Description 用户可以在某一篇文献的评论区中发表自己的评论
// @Tags        社交
// @Param       data  body   response.CommentCreation true "data"
// @Param       token header string                   true "token"
// @Accept      json
// @Produce     json
// @Success     200 {string} json "{"success":true,"status":200,"msg":"评论创建成功","comment_id":string}"
// @Failure     400 {string} json "{"success":false,"status":400,"msg":"用户ID不存在"}"
// @Failure     403 {string} json "{"success":false,"status":403,"msg":"评论创建失败"}"
// @Router      /social/comment/create [POST]
func CreateComment(c *gin.Context) {
	//user_id := c.Query("user_id")
	//paper_id := c.Query("paper_id")
	var d response.CommentCreation
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	content := d.Content
	userID := d.UserID
	paper_id := d.PaperID
	authuser, exist := c.Get("user")
	if exist {
		userID = authuser.(database.User).UserID
	}
	//userID, _ := strconv.ParseUint(user_id, 0, 64)
	//验证用户是否存在
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}

	comment := database.Comment{UserID: user.UserID, PaperID: paper_id,
		CommentTime: time.Now(), Content: content}
	// new a comment auto generate comment_id
	err := service.CreateComment(&comment)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "status": 403, "msg": "评论创建失败"})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success":    true,
			"status":     200,
			"msg":        "评论创建成功",
			"comment_id": comment.CommentID})
	}
}

// LikeComment 点赞评论
// @Summary     用户可以对某一评论进行点赞 Vera
// @Description 用户可以对某一评论进行点赞
// @Tags        社交
// @Param       data body     response.CommentUser true "data"
// @Success     200  {string} json                 "{"success": true,"status":200,"msg": "操作成功"}"
// @Failure     400  {string} json                 "{"success": false,"status":400,"msg":"用户ID不存在"}"
// @Failure     402  {string} json                 "{"success": false,"status":402, "msg": "用户已赞过该评论"}"
// @Failure     403  {string} json                 "{"success": false,"status":403, "msg": "评论不存在"}"
// @Router      /social/comment/like [POST]
func LikeComment(c *gin.Context) {
	//user_id := c.Query("user_id")
	//comment_id := c.Query("comment_id")
	//fmt.Printf("debug 0")
	var d response.CommentUser
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	userID := d.UserID
	commentID := d.CommentID
	//userID, _ := strconv.ParseUint(user_id, 0, 64)
	//commentID, _ := strconv.ParseUint(comment_id, 0, 64)
	//验证用户是否存在
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	//fmt.Printf("debug 1")
	//commentID, _ := strconv.ParseUint(comment_id, 0, 64)
	comment, notFound := service.GetCommentByID(commentID)
	if notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"msg":     "评论不存在",
		})
		return
	}
	//fmt.Printf("debug 2")

	isLike := service.GetLike_Rel(commentID, userID)
	if isLike {
		c.JSON(http.StatusOK, gin.H{
			"status": 402,
			"msg":    "用户已赞过该评论",
		})
		return
	}
	//fmt.Printf("debug 3")

	service.UpdateCommentLike(comment, user)
	c.JSON(http.StatusOK, gin.H{"success": true,
		"status": 200,
		"msg":    "操作成功"})
}

// UnLikeComment  取消点赞
// @Summary     取消点赞 Vera
// @Description 取消点赞
// @Tags        社交
// @Param       data body     response.CommentUser true "data"
// @Success     200  {string} json                 "{"success": true,"status":200,"msg": "已取消点赞"}"
// @Failure     400  {string} json                 "{"success": false,"status":400,"msg":"用户ID不存在"}"
// @Failure     402  {string} json                 "{"success": false,"status":402, "msg": "用户未点赞"}"
// @Failure     403  {string} json                 "{"success": false,"status":403, "msg": "评论不存在"}"
// @Router      /social/comment/unlike [POST]
func UnLikeComment(c *gin.Context) {
	//user_id := c.Query("user_id")
	//comment_id := c.Query("comment_id")
	//
	//userID, _ := strconv.ParseUint(user_id, 0, 64)
	//commentID, _ := strconv.ParseUint(comment_id, 10, 64)
	var d response.CommentUser
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	userID := d.UserID
	commentID := d.CommentID
	//验证用户是否存在
	user, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	comment, notFound := service.GetCommentByID(commentID)
	if notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  403,
			"msg":     "评论不存在",
		})
		return
	}
	//isLike := service.GetLike_Rel(comment_id, userID)
	//if !isLike {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": 402,
	//		"msg":    "用户未赞过该评论",
	//	})
	//	return
	//}
	notFound = service.CancelLike(comment, user)
	if !notFound {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"status":  200,
			"msg":     "已取消点赞",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"status":  402,
			"message": "用户未点赞",
		})
		return
	}
}

// ShowPaperCommentList 取消点赞
// @Summary     显示文献评论列表，时间倒序 Vera
// @Description 显示文献评论列表，时间倒序
// @Tags        社交
// @Param       data body response.CommentListQuery true "data"
// @Accept      json
// @Produce     json
// @Success     200 {string} string "{"data":{"comments":[],"paper_id":"string"},"message":"查找成功"}"
// @Failure     404 {string} string "{"success": false, "status":  403,"message": "评论用户不存在"}"
// @Failure     400 {string} string "{"status": 400, "msg": "用户ID不存在"}"
// @Router      /social/comment/list [POST]
func ShowPaperCommentList(c *gin.Context) {
	var d response.CommentListQuery
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	userID := d.UserID
	paper_id := d.PaperID
	//userID, _ := strconv.ParseUint(user_id, 0, 64)
	//验证用户是否存在
	_, notFoundUserByID := service.QueryAUserByID(userID)
	if notFoundUserByID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "用户ID不存在",
		})
		return
	}
	comments := service.GetCommentsByPaperId(paper_id)
	//fmt.Println(comments)
	//if len(comments) == 0 {
	//	c.JSON(http.StatusOK, gin.H{
	//		"success": false,
	//		"status":  403,
	//		"message": "评论不存在",
	//	})
	//	return
	//}
	var dataList = make([]map[string]interface{}, 0)
	for _, comment := range comments {
		var com = make(map[string]interface{})
		//com["id"] = comment.CommentID
		//com["like"] = comment.LikeNum
		//com["is_animating"] = false
		//com["is_like"] = false
		//if !err && service.GetLike_Rel(comment.CommentID, userID) {
		//	com["is_like"] = true
		//}
		com["user_id"] = comment.UserID
		//ccf-return Comment_User's username and headshot
		user, notFoundUserByID := service.QueryAUserByID(comment.UserID)
		if notFoundUserByID {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"status":  404,
				"message": "评论用户不存在",
			})
			return
		}
		com["username"] = user.Username
		com["headshot"] = user.HeadShot
		com["userinfo"] = user.UserInfo
		//com["username"] = comment.Username
		com["content"] = comment.Content
		com["time"] = comment.CommentTime
		//com["reply_count"] = comment.ReplyCount
		// fmt.Println(com)
		dataList = append(dataList, com)
	}
	// fmt.Println(dataList)
	var data = make(map[string]interface{})
	data["paper_id"] = paper_id
	//data["paper_title"] = comments[0].PaperTitle
	data["comments"] = dataList
	c.JSON(http.StatusOK, gin.H{
		"message": "查找成功",
		"data":    data,
	})
}

// FollowAuthor 关注学者
// @Summary     关注学者 包括了关注和取消关注（通过重复调用来实现） txc
// @Description 关注学者 包括了关注和取消关注（通过重复调用来实现）
// @Tags        社交
// @Accept      json
// @Produce     json
// @Param       data  body     response.FollowAuthorQ true "data"
// @Param       token header   string                 true "token"
// @Success     200   {string} string                 "{"msg": "取消关注成功/关注成功"}"
// @Failure     400   {string} string                 "{"err":err,"msg": "参数错误"}"
// @Router      /social/follow [POST]
// @security    ApiKeyAuth
func FollowAuthor(c *gin.Context) {
	//authuser, _ := c.Get("user")
	var d response.FollowAuthorQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": err,
			"msg": "参数错误",
		})
		return
	}
	_, notFound := service.QueryAUserByID(d.UserID)
	if notFound {
		c.JSON(401, gin.H{"msg": "用户ID不存在"})
		return
	}
	author, err, _ := service.GetObject2("authors", d.AuthorID)
	if err != nil {
		c.JSON(402, gin.H{"msg": "学者不存在"})
		return
	}
	uf, notFound := service.GetUserFollow(d.UserID, d.AuthorID)
	if notFound {
		uf = database.UserFollow{
			UserID:     d.UserID,
			AuthorID:   d.AuthorID,
			AuthorName: author["display_name"].(string),
		}
		err := service.CreateUserFollow(&uf)
		if err != nil {
			c.JSON(403, gin.H{"msg": "关注失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "关注成功"})
		return
	}
	err = service.DeleteUserFollow(&uf)
	if err != nil {
		c.JSON(404, gin.H{"msg": "取消关注失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "取消关注成功"})
}

// GetUserFollows
// @Summary     获取用户关注的学者 txc
// @Description 获取用户关注的学者
// @Tags        社交
// @Accept      json
// @Produce     json
// @Param       data  body     response.GetUserFollowsQ true "data"
// @Param       token header   string                   true "token"
// @Success     200   {string} string                   "{"msg": "查找成功","data":data,"count":count}"
// @Failure     400   {string} string                   "{"err":err,"msg": "参数错误"}"
// @Failure     401   {string} string                   "{"msg": "用户ID不存在"}"
// @Router      /social/follow/list [POST]
func GetUserFollows(c *gin.Context) {
	var d response.GetUserFollowsQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err, "msg": "参数错误"})
		return
	}
	_, notFound := service.QueryAUserByID(d.UserID)
	if notFound {
		c.JSON(401, gin.H{"msg": "用户ID不存在"})
		return
	}
	ufs, _ := service.GetUserFollows(d.UserID)
	c.JSON(http.StatusOK, gin.H{
		"message": "查找成功",
		"data":    ufs,
		"count":   len(ufs),
	})
}
