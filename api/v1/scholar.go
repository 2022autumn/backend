package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AddUserConcept
// @Summary     添加user的关注关键词 txc
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
		res, err := service.GetObject("concepts", d.ConceptID)
		if err != nil {
			c.JSON(402, gin.H{"msg": "concept不存在"})
			return
		}
		var tmp map[string]interface{}
		_ = json.Unmarshal(res.Source, &tmp)
		userConcept = database.UserConcept{
			UserID:      d.UserID,
			ConceptID:   d.ConceptID,
			ConceptName: tmp["display_name"].(string),
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
// @Summary     获取用户推荐的论文 请勿使用 txc
// @Description 获取用户推荐的论文 请勿使用
// @Tags        scholar
// @Param       userid query    string false "userid"
// @Success     200    {string} json   "{"msg":"获取成功","data":{}}"
// @Router      /scholar/roll [GET]
func RollWorks(c *gin.Context) {
	userID := c.Query("userid")
	ret := make([]map[string]interface{}, 0)
	rand.Seed(time.Now().UnixNano())
	filter := utils.InitWorksfilter()
	if userID != "" {
		userID, _ := strconv.ParseUint(userID, 10, 64)
		ucs, _ := service.GetUserConcepts(userID)
		ufs, _ := service.GetUserFollows(userID)
		if len(ucs) != 0 {
			uc := ucs[rand.Intn(len(ucs))]
			url := "https://api.openalex.org/works?filter=concepts.id:" + uc.ConceptID
			works := make([]map[string]interface{}, 0)
			_, err := service.GetWorksByUrl(url, 1, &works)
			if err != nil {
				c.JSON(400, gin.H{"msg": "获取失败"})
				return
			}
			workids := make([]string, 0)
			for i, work := range works {
				workids = append(workids, utils.RemovePrefix(work["id"].(string)))
				if i > 8 {
					break
				}
			}
			rand.Shuffle(len(workids), func(i, j int) { workids[i], workids[j] = workids[j], workids[i] })
			res, err := service.GetObjects("works_v1", workids)
			if err == nil {
				for _, work := range res.Docs {
					if work.Found {
						var tmp map[string]interface{}
						_ = json.Unmarshal(work.Source, &tmp)
						utils.FilterData(&tmp, &filter)
						ret = append(ret, map[string]interface{}{
							"work":   tmp,
							"source": "concept",
							"name":   uc.ConceptName,
						})
						break
					}
				}
			}
		}
		if len(ufs) != 0 {
			uf := ufs[rand.Intn(len(ufs))]
			url := "https://api.openalex.org/works?filter=author.id:" + uf.AuthorID
			works := make([]map[string]interface{}, 0)
			_, err := service.GetWorksByUrl(url, 1, &works)
			if err != nil {
				c.JSON(400, gin.H{"msg": "获取失败"})
				return
			}
			workids := make([]string, 0)
			for i, work := range works {
				workids = append(workids, utils.RemovePrefix(work["id"].(string)))
				if i > 8 {
					break
				}
			}
			rand.Shuffle(len(workids), func(i, j int) { workids[i], workids[j] = workids[j], workids[i] })
			res, err := service.GetObjects("works_v1", workids)
			if err == nil {
				for _, work := range res.Docs {
					if work.Found {
						var tmp map[string]interface{}
						_ = json.Unmarshal(work.Source, &tmp)
						utils.FilterData(&tmp, &filter)
						ret = append(ret, map[string]interface{}{
							"work":   tmp,
							"source": "author",
							"name":   uf.AuthorName,
						})
						break
					}
				}
			}
		}
	}
	if len(ret) < 6 {
		var count int
		global.DB.Table("work_views").Count(&count)
		ids := rand.Perm(count)
		workids := make([]string, 0)
		for i := len(ret); i < 6; i++ {
			var work database.WorkView
			global.DB.Table("work_views").Offset(ids[i]).First(&work)
			workids = append(workids, work.WorkID)
		}
		res, err := service.GetObjects("works_v1", workids)
		if err == nil {
			for _, work := range res.Docs {
				ret = append(ret, map[string]interface{}{
					"work":   work.Source,
					"source": "random",
					"name":   "",
				})
			}
		}
	}
	c.JSON(200, gin.H{"msg": "获取成功", "data": ret})
}

// GetHotWorks
// @Summary     获取热门论文（根据访问量） txc
// @Description 获取热门论文（根据访问量）
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

// GetPersonalWorks
// @Summary     获取学者的论文 hr 未测试
// @Description 获取学者的论文
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param		author_id body string true "author_id 是作者的id"
// @Param		page body int true "page 获取第几页的数据"
// @Param		page_size body int true "page_size 是分页的大小"
// @Success     200 {string} json "{"msg":"获取成功","data":{}}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"作者不存在"}"
// @Failure     402 {string} json "{"msg":"该作者没有论文"}"
// @Router      /scholar/works [GET]
func GetPersonalWorks(c *gin.Context) {
	var d response.GetPersonalWorksQ
	author_id, page, page_size := d.AuthorID, d.Page, d.PageSize
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	res, err := service.GetObject("authors", author_id)
	if err != nil {
		c.JSON(401, gin.H{"msg": "作者不存在"})
		return
	}
	works, notFound := service.GetScholarWorks(author_id)
	if !notFound { // 能找到则从数据库中获取
		// 按照place排序 从小到大
		sort.Slice(works, func(i, j int) bool {
			return works[i].Place < works[j].Place
		})
		// 分页
		works = works[(page-1)*page_size : page*page_size]
		c.JSON(200, gin.H{"msg": "获取成功", "data": works})
		return
	}
	// 不能找到则从openalex api中获取
	author := res.Source
	var author_map map[string]interface{}
	_ = json.Unmarshal(author, &author_map)
	works_api_url := author_map["works_api_url"].(string)
	works = make([]database.PersonalWorks, 0)
	works_detail := make([]map[string]interface{}, 0)
	service.GetAllWorksByUrl(works_api_url, &works_detail)
	for i, work := range works_detail {
		presonal_work := database.PersonalWorks{
			AuthorID: author_id,
			WorkID:   work["id"].(string),
			Place:    i,
		}
		_ = service.AddScholarWork(&presonal_work)
		works = append(works, presonal_work)
	}
	// 分页
	works = works[(page-1)*page_size : page*page_size]
	if len(works) == 0 {
		c.JSON(402, gin.H{"msg": "该作者没有论文"})
		return
	}
	c.JSON(200, gin.H{"msg": "获取成功", "data": works})
}

// IgnoreWork 忽略论文
// @Summary     学者管理主页--忽略论文 hr 未测试
// @Description 学者管理主页--忽略论文 通过重复调用该接口可以完成论文的忽略与取消忽略
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param		author_id body string true "author_id 是作者的id"
// @Param		work_id body string true "work_id 是论文的id"
// @Param		ignore body bool true "ignore 是当前论文的忽略状态"
// @Success     200 {string} json "{"msg":"忽略成功"}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"忽略失败"}"
// @Router      /scholar/ignore [POST]
func IgnoreWork(c *gin.Context) {
	var d response.IgnoreWorkQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	author_id, work_id, ignore := d.AuthorID, d.WorkID, d.Ignore
	err := service.UpdateWorkIgnore(author_id, work_id, !ignore)
	if err != nil {
		c.JSON(401, gin.H{"msg": "忽略失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "忽略成功"})
}

// ModifyPlace 修改论文顺序
// @Summary     学者管理主页--修改论文顺序 hr 未测试
// @Description 学者管理主页--修改论文顺序
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param		author_id body string true "author_id 是作者的id"
// @Param		work_id body string true "work_id 是论文的id"
// @Param		direction body int true "direction 是论文的移动方向，1为向上，-1为向下"
// @Success     200 {string} json "{"msg":"修改成功"}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"未找到该论文"}"
// @Failure     402 {string} json "{"msg":"论文已经在顶部"}"
// @Failure     403 {string} json "{"msg":"论文已经在底部"}"
// @Failure     404 {string} json "{"msg":"修改失败"}"
// @Router      /scholar/modify [POST]
func ModifyPlace(c *gin.Context) {
	var d response.ModifyPlaceQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数数目或类型错误"})
		return
	}
	author_id, work_id, direction := d.AuthorID, d.WorkID, d.Direction
	if direction != 1 && direction != -1 {
		c.JSON(400, gin.H{"msg": "direction参数错误"})
		return
	}
	// 获取当前论文的place
	place, notFound := service.GetWorkPlace(author_id, work_id)
	if notFound {
		c.JSON(401, gin.H{"msg": "未找到该论文"})
		return
	}
	// 获取论文总数
	total, err := service.GetScholarWorksCount(author_id)
	if err != nil {
		c.JSON(404, gin.H{"msg": "查询论文总数出错，修改失败"})
		return
	}
	// 判断论文是否在顶部或底部
	if place == 0 && direction == -1 {
		c.JSON(402, gin.H{"msg": "论文已经在顶部"})
		return
	}
	if place == total-1 && direction == 1 {
		c.JSON(403, gin.H{"msg": "论文已经在底部"})
		return
	}
	target_place := place + direction
	// 获取目标论文的id
	target_work, notFound := service.GetWorkByPlace(author_id, target_place)
	if notFound {
		c.JSON(404, gin.H{"msg": "获取目标论文失败,修改失败"})
		return
	}
	// 交换两篇论文的place
	err = service.SwapWorkPlace(author_id, work_id, target_work.WorkID)
	if err != nil {
		c.JSON(404, gin.H{"msg": "交换ID失败,修改失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "修改成功"})
}

// 置顶论文
// @Summary     学者管理主页--置顶论文 hr 未测试
// @Description 学者管理主页--置顶论文 通过重复调用而取消置顶
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param		author_id body string true "author_id 是作者的id"
// @Param		work_id body string true "work_id 是论文的id"
// @Success     200 {string} json "{"msg":"置顶成功"}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"未找到该论文"}"
// @Failure     402 {string} json "{"msg":"修改失败"}"
// @Router      /scholar/top [POST]
func TopWork(c *gin.Context) {
	var d response.TopWorkQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数数目或类型错误"})
		return
	}
	author_id, work_id := d.AuthorID, d.WorkID
	// 获取当前论文的place
	_, notFound := service.GetWorkPlace(author_id, work_id)
	if notFound {
		c.JSON(401, gin.H{"msg": "未找到该论文"})
		return
	}
	// 置顶论文
	err := service.TopWork(author_id, work_id)
	if err != nil {
		c.JSON(402, gin.H{"msg": "修改失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "置顶成功"})
}

// UploadAuthorHeadshot
// @Summary     上传作者头像 txc
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
