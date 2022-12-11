package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"encoding/json"
	"fmt"
	"log"
	"math"
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
	user := c.MustGet("user").(database.User)
	var d response.AddUserConceptQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	if _, notFound := service.GetUserByID(user.UserID); notFound {
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
	userConcept, notFound := service.GetUserConcept(user.UserID, d.ConceptID)
	if notFound {
		res, err := service.GetObject("concepts", d.ConceptID)
		if err != nil {
			c.JSON(402, gin.H{"msg": "concept不存在"})
			return
		}
		var tmp map[string]interface{}
		_ = json.Unmarshal(res.Source, &tmp)
		userConcept = database.UserConcept{
			UserID:      user.UserID,
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
// @Summary     获取学者的论文 hr
// @Description 获取学者的论文
// @Description
// @Description 参数说明
// @Description - author_id 作者的id
// @Description
// @Description - page 获取第几页的数据
// @Description
// @Description - page_size 分页的大小
// @Description
// @Description - display 是否显示已删除的论文 -1不显示 1显示
// @Description 返回值说明
// @Description - msg 返回信息
// @Description
// @Description - res 返回该页的works对象数组
// @Description
// @Description - pages 分页总数
// @Tags       	学者主页的论文获取、管理
// @Accept      json
// @Produce     json
// @Param		data body response.GetPersonalWorksQ true "data 是请求参数,包括author_id ,page ,page_size, display"
// @Success     200 {string} json "{"msg":"获取成功","res":{}, "pages":{}}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"作者不存在"}"
// @Failure     402 {string} json "{"msg":"page超出范围"}"
// @Failure     403 {string} json "{"msg":"该作者没有论文"}"
// @Router      /scholar/works/get [POST]
func GetPersonalWorks(c *gin.Context) {
	var d response.GetPersonalWorksQ
	var data []string
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	author_id, page, page_size, display := d.AuthorID, d.Page, d.PageSize, d.Display
	if author_id == "" {
		c.JSON(400, gin.H{"msg": "author_id 为空，参数错误"})
		return
	}
	res, err := service.GetObject("authors", author_id)
	if err != nil {
		c.JSON(401, gin.H{"msg": "作者不存在"})
		return
	}
	var works []database.PersonalWorks
	var notFound bool
	if display == -1 {
		works, notFound = service.GetScholarDisplayWorks(author_id)
	} else {
		works, notFound = service.GetScholarAllWorks(author_id)
	}
	if !notFound && len(works) != 0 { // 能找到则从数据库中获取
		// 按照place排序 从小到大
		sort.Slice(works, func(i, j int) bool {
			return works[i].Place < works[j].Place
		})
		// 总页数,向上取整
		pages := int(math.Ceil(float64(len(works)) / float64(page_size)))
		// 分页
		if page > pages {
			c.JSON(402, gin.H{"msg": "page超出范围(database)"})
			return
		}
		// 页数从1开始
		if page == pages { // 最后一页
			works = works[(page-1)*page_size:]
		} else {
			works = works[(page-1)*page_size : page*page_size]
		}
		// 获取works_id
		for _, work := range works {
			data = append(data, work.WorkID)
		}
		c.JSON(200, gin.H{"msg": "获取成功", "res": data, "pages": pages})
		return
	}
	// 不能找到则从openalex api中获取
	log.Println("从openalex api中获取author works")
	author := res.Source
	var author_map map[string]interface{}
	_ = json.Unmarshal(author, &author_map)
	works_api_url := author_map["works_api_url"].(string)
	works = make([]database.PersonalWorks, 0)
	service.GetAllPersonalWorksByUrl(works_api_url, &works, author_id)
	if len(works) == 0 {
		c.JSON(403, gin.H{"msg": "该作者没有论文"})
		return
	}
	// 总页数,向上取整
	pages := int(math.Ceil(float64(len(works)) / float64(page_size)))
	// 分页
	if page > pages {
		c.JSON(402, gin.H{"msg": "page超出范围"})
		return
	}
	go service.CreateWorks(works)
	// 页数从1开始
	if page == pages { // 最后一页
		works = works[(page-1)*page_size:]
	} else {
		works = works[(page-1)*page_size : page*page_size]
	}
	// 获取works_id
	for _, work := range works {
		data = append(data, work.WorkID)
	}
	c.JSON(200, gin.H{"msg": "获取成功", "res": data, "pages": pages})
}

// IgnoreWork 忽略论文
// @Summary     学者管理主页--忽略论文 hr
// @Description 学者管理主页--忽略论文 通过重复调用该接口可以完成论文的忽略与取消忽略
// @Description
// @Description 参数说明
// @Description - author_id 作者的id
// @Description
// @Description - work_id 论文的id
// @Tags        学者主页的论文获取、管理
// @Accept      json
// @Produce     json
// @Param		data body response.IgnoreWorkQ true "data 是请求参数,包括author_id ,work_id"
// @Success     200 {string} json "{"msg":"修改忽略属性成功"}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"修改忽略属性失败"}"
// @Router      /scholar/works/ignore [POST]
func IgnoreWork(c *gin.Context) {
	var d response.IgnoreWorkQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		log.Println(err)
		return
	}
	author_id, work_id := d.AuthorID, d.WorkID
	err := service.IgnoreWork(author_id, work_id)
	if err != nil {
		c.JSON(401, gin.H{"msg": "修改忽略属性失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "修改忽略属性成功"})
}

// ModifyPlace 修改论文顺序
// @Summary     学者管理主页--修改论文顺序 hr
// @Description 学者管理主页--修改论文顺序
// @Description
// @Description 参数说明
// @Description - author_id 作者的id
// @Description
// @Description - work_id 论文的id
// @Description
// @Description - direction 论文的移动方向，1为向上，-1为向下
// @Tags        学者主页的论文获取、管理
// @Accept      json
// @Produce     json
// @Param		data body response.ModifyPlaceQ true "data 是请求参数,包括author_id ,work_id ,direction"
// @Success     200 {string} json "{"msg":"修改成功"}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"未找到该论文"}"
// @Failure     402 {string} json "{"msg":"论文已经在顶部"}"
// @Failure     403 {string} json "{"msg":"论文已经在底部"}"
// @Failure     404 {string} json "{"msg":"修改失败"}"
// @Router      /scholar/works/modify [POST]
func ModifyPlace(c *gin.Context) {
	var d response.ModifyPlaceQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数数目或类型错误"})
		log.Println(err)
		return
	}
	author_id, work_id, direction := d.AuthorID, d.WorkID, d.Direction
	if direction != 1 && direction != -1 {
		c.JSON(400, gin.H{"msg": "direction参数错误"})
		return
	}
	// 获取当前论文的place
	place, notFound := service.GetWorkPlace(author_id, work_id)
	if notFound || place == -1 {
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
	if place == 0 && direction == 1 {
		c.JSON(402, gin.H{"msg": "论文已经在顶部"})
		return
	}
	if place == total-1 && direction == -1 {
		c.JSON(403, gin.H{"msg": "论文已经在底部"})
		return
	}
	target_place := place - direction
	// 获取目标论文的id
	target_work, notFound := service.GetWorkByPlace(author_id, target_place)
	if notFound || target_work.Place == -1 {
		c.JSON(404, gin.H{"msg": "获取交换目标论文失败,修改失败"})
		return
	}
	// 交换两篇论文的place
	log.Println("target_work.WorkID", target_work.WorkID)
	err = service.SwapWorkPlace(author_id, work_id, target_work.WorkID)
	if err != nil {
		c.JSON(404, gin.H{"msg": "交换ID失败,修改失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "修改成功"})
}

// 置顶论文
// @Summary     学者管理主页--置顶论文 hr
// @Description 学者管理主页--置顶论文 通过重复调用而取消置顶
// @Description
// @Description 参数说明
// @Description - author_id 作者的id
// @Description
// @Description - work_id 论文的id
// @Tags        学者主页的论文获取、管理
// @Accept      json
// @Produce     json
// @Param		data body response.TopWorkQ true "data 是请求参数,包括author_id ,work_id"
// @Success     200 {string} json "{"msg":"置顶成功"}"
// @Failure     400 {string} json "{"msg":"参数错误"}"
// @Failure     401 {string} json "{"msg":"未找到该论文"}"
// @Failure     402 {string} json "{"msg":"修改失败"}"
// @Router      /scholar/works/top [POST]
func TopWork(c *gin.Context) {
	var d response.TopWorkQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数数目或类型错误"})
		return
	}
	author_id, work_id := d.AuthorID, d.WorkID
	// 获取当前论文的place
	place, notFound := service.GetWorkPlace(author_id, work_id)
	if notFound || place == -1 {
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
