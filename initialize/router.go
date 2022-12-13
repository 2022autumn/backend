package initialize

import (
	v1 "IShare/api/v1"
	v2 "IShare/api/v2"
	"IShare/docs"
	"IShare/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine) {
	r.Use(middleware.Cors()) // 跨域
	// r.Use(middleware.LoggerToFile()) // 日志

	docs.SwaggerInfo.Title = "ishare backend doc"
	docs.SwaggerInfo.Version = "v1"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/api/test", testGin)

	baseGroup := r.Group("/api")
	{
		//用户模块
		baseGroup.POST("/register", v1.Register)            //注册
		baseGroup.POST("/login", v1.Login)                  //登录
		baseGroup.GET("/user/info", v1.UserInfo)            //个人中心
		baseGroup.POST("/user/mod", v1.ModifyUser)          //编辑个人信息
		baseGroup.POST("/user/pwd", v1.ModifyPassword)      //重置用户密码
		baseGroup.POST("/user/headshot", v1.UploadHeadshot) //上传用户头像
		baseGroup.GET("/info/register_num", v1.GetRegisterUserNum)
		baseGroup.GET("/info/verified_num", v1.GetVerifiedUserNum)
		baseGroup.POST("/user/history", middleware.AuthRequired(), v1.GetBrowseHistory)
		baseGroup.Static("/media", "./media")
	}
	ApplicationRouter := baseGroup.Group("/application")
	{
		ApplicationRouter.POST("/create", v1.CreateApplication)
		ApplicationRouter.POST("/handle", v1.HandleApplication)
		ApplicationRouter.GET("/list", v1.UncheckedApplicationList)
		ApplicationRouter.POST("/code", v1.SendVerifyEmail)
		ApplicationRouter.GET("/testCode", v1.TestCodeGen)
	}
	SocialRouter := baseGroup.Group("/social")
	{
		SocialRouter.POST("/comment/create", middleware.AuthRequired(), v1.CreateComment)
		SocialRouter.POST("/comment/like", v1.LikeComment)
		SocialRouter.POST("/comment/unlike", v1.UnLikeComment)
		SocialRouter.POST("/comment/list", v1.ShowPaperCommentList)
		//SocialRouter.POST("/follow", v1.FollowAuthor)
		//SocialRouter.POST("/follow/list", v1.GetUserFollows)
		SocialRouter.POST("/tag/create", v1.CreateTag)
		SocialRouter.POST("/tag/collectPaper", v1.AddTagToPaper)
		SocialRouter.POST("/tag/sublist", v1.ShowTagPaperList)
		SocialRouter.POST("/tag/taglist", v1.ShowUserTagList)
		SocialRouter.POST("/tag/delete", v1.DeleteTag)
		SocialRouter.POST("/tag/cancelCollectPaper", v1.RemovePaperTag)
		SocialRouter.POST("/tag/rename", v1.RenameTag)
		SocialRouter.POST("/tag/paperTagList", v1.PaperTagList)
		SocialRouter.POST("/follow", middleware.AuthRequired(), v1.FollowAuthor)
		SocialRouter.POST("/follow/list", middleware.AuthRequired(), v1.GetUserFollows)
	}
	esGroup := baseGroup.Group("/es")
	{
		esGroup.GET("/statistic", v1.GetStatistics)
		esGroup.GET("/get/", v1.GetObject)
		esGroup.GET("/get2/", v2.GetObject2)
		esGroup.POST("/search/base", v1.BaseSearch)
		esGroup.POST("/search/doi", v1.DoiSearch)
		esGroup.POST("/search/advanced", v1.AdvancedSearch)
		esGroup.GET("/search/author", v1.AuthorSearch)
		esGroup.GET("/search/author2", v1.AuthorSearch2)
		esGroup.GET("/getAuthorRelationNet", v1.GetAuthorRelationNet)
		esGroup.POST("/prefix", v1.GetPrefixSuggestions)
	}
	scholarGroup := baseGroup.Group("/scholar")
	{
		scholarGroup.POST("/concept", middleware.AuthRequired(), v1.AddUserConcept)
		scholarGroup.GET("/concept", middleware.AuthRequired(), v1.GetUserConcepts)
		scholarGroup.GET("/roll", v1.RollWorks)
		scholarGroup.GET("/hot", v1.GetHotWorks)
		scholarGroup.POST("/author/headshot", middleware.AuthRequired(), v1.UploadAuthorHeadshot)
		scholarGroup.POST("/author/intro", middleware.AuthRequired(), v1.ModifyAuthorIntro)
	}
	// 学者主页论文
	personalWorksGroup := scholarGroup.Group("/works")
	{
		personalWorksGroup.POST("/get", v1.GetPersonalWorks)
		personalWorksGroup.POST("/ignore", v1.IgnoreWork)
		personalWorksGroup.POST("/modify", v1.ModifyPlace)
		personalWorksGroup.POST("/top", v1.TopWork)
		personalWorksGroup.POST("/untop", v1.UnTopWork)
		personalWorksGroup.POST("/upload", v1.UploadPaperPDF)
		personalWorksGroup.POST("/unupload", v1.UnUploadPaperPDF)
		personalWorksGroup.POST("/getpdf", v1.GetPaperPDF)
	}
}

func testGin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"success": true,
	})
}
