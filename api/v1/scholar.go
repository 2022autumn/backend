package v1

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"path"
	"sort"
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
		res, err, _ := service.GetObject2("concepts", d.ConceptID)
		if err != nil {
			c.JSON(402, gin.H{"msg": "concept不存在"})
			return
		}
		userConcept = database.UserConcept{
			UserID:      user.UserID,
			ConceptID:   d.ConceptID,
			ConceptName: res["display_name"].(string),
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

// GetUserConcepts
// @Summary     获取用户关注的关键词 txc
// @Description 获取用户关注的关键词
// @Tags        scholar
// @Param       token header   string                      true "token"
// @Success     200   {string} json   "{"msg":"获取成功","data":{}}"
// @Failure     401   {string} json   "{"msg":"数据库获取失败"}"
// @Router      /scholar/concept [GET]
func GetUserConcepts(c *gin.Context) {
	user := c.MustGet("user").(database.User)
	concepts, err := service.GetUserConcepts(user.UserID)
	if err != nil {
		c.JSON(401, gin.H{"msg": "数据库获取失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "获取成功", "data": concepts})
}

// RollWorks
// @Summary     获取用户推荐的论文 请勿使用 txc
// @Description 获取用户推荐的论文 请勿使用
// @Tags        scholar
// @Param       concept_id query    string false "concept_id"
// @Success     200        {string} json   "{"msg":"获取成功","data":{}}"
// @Failure     401        {string} json   "{"msg":"openalex获取失败"}"
// @Router      /scholar/roll [GET]
func RollWorks(c *gin.Context) {
	ret := make([]map[string]interface{}, 0)
	retSize := 6
	rand.Seed(time.Now().UnixNano())
	conceptID := c.Query("concept_id")
	var pages = rand.Perm(5)
	if conceptID != "" {
		if conceptID[0] != 'C' {
			c.JSON(400, gin.H{"msg": "concept_id参数错误"})
			return
		}
		url := "https://api.openalex.org/works?filter=concepts.id:" + conceptID
		page := 1
		for true {
			works := make([]map[string]interface{}, 0)
			total_pages, err := service.GetWorksByUrl(url, pages[page]+1, &works)
			if err != nil {
				c.JSON(401, gin.H{"msg": "openalex获取失败"})
				return
			}
			workids := make([]string, 0)
			for _, work := range works {
				workids = append(workids, utils.RemovePrefix(work["id"].(string)))
			}
			rand.Shuffle(len(workids), func(i, j int) { workids[i], workids[j] = workids[j], workids[i] })
			res, err := service.GetObjects2("works", workids)
			if err == nil {
				works := res["results"].([]interface{})
				for _, work := range works {
					work := work.(map[string]interface{})
					work["id"] = utils.RemovePrefix(work["id"].(string))
					if work["abstract_inverted_index"] != nil {
						work["abstract"] = utils.TransInvertedIndex2String(work["abstract_inverted_index"].(map[string]interface{}))
						work["abstract_inverted_index"] = nil
					}
					if work["abstract"] == nil {
						work["abstract"] = ""
					}
					ret = append(ret, map[string]interface{}{
						"work": work,
					})
					if len(ret) == retSize {
						c.JSON(200, gin.H{"msg": "获取成功", "data": ret})
						return
					}
				}
			}
			page++
			if page == total_pages {
				break
			}
		}
	}
	if len(ret) < retSize {
		var count int
		global.DB.Table("work_views").Count(&count)
		ids := rand.Perm(count)
		workids := make([]string, 0)
		for i := len(ret); i < retSize; i++ {
			var work database.WorkView
			global.DB.Table("work_views").Offset(ids[i]).First(&work)
			workids = append(workids, work.WorkID)
		}
		res, err := service.GetObjects2("works", workids)
		if err == nil {
			works := res["results"].([]interface{})
			for _, work := range works {
				work := work.(map[string]interface{})
				work["id"] = utils.RemovePrefix(work["id"].(string))
				if work["abstract_inverted_index"] != nil {
					work["abstract"] = utils.TransInvertedIndex2String(work["abstract_inverted_index"].(map[string]interface{}))
					work["abstract_inverted_index"] = nil
				}
				if work["abstract"] == nil {
					work["abstract"] = ""
				}
				ret = append(ret, map[string]interface{}{
					"work": work,
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
// @Description - page 获取第几页的数据, START FROM 1
// @Description
// @Description - page_size 分页的大小, 不能为0
// @Description
// @Description - display 是否显示已删除的论文 -1不显示 1显示
// @Description
// @Description 返回值说明
// @Description - msg 返回信息
// @Description
// @Description - data 返回该页的works对象数组
// @Description
// @Description - pages 分页总数，一共有多少页
// @Description
// @Description - total 论文总数
// @Tags        学者主页的论文获取、管理
// @Accept      json
// @Produce     json
// @Param       data body     response.GetPersonalWorksQ true "data 是请求参数,包括author_id ,page ,page_size, display"
// @Success     200  {string} json                       "{"msg":"获取成功","data":{}, "pages":{}, "total":{}}"
// @Failure     400  {string} json                       "{"msg":"参数错误"}"
// @Failure     401  {string} json                       "{"msg":"作者不存在"}"
// @Failure     402  {string} json                       "{"msg":"page超出范围"}"
// @Failure     403  {string} json                       "{"msg":"该作者没有论文"}"
// @Router      /scholar/works/get [POST]
func GetPersonalWorks(c *gin.Context) {
	var d response.GetPersonalWorksQ
	var works_ids []string
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	author_id, page, page_size, display := d.AuthorID, d.Page, d.PageSize, d.Display
	if author_id == "" {
		c.JSON(400, gin.H{"msg": "author_id 为空，参数错误"})
		return
	}
	res, err, _ := service.GetObject2("authors", author_id)
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
		// 先按照Top排序，从大到小，再按照place排序 从小到大
		sort.Slice(works, func(i, j int) bool {
			if works[i].Top == works[j].Top {
				return works[i].Place < works[j].Place
			}
			return works[i].Top > works[j].Top
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
			works_ids = append(works_ids, work.WorkID)
		}
		objects, err := service.GetObjects2("works", works_ids)
		if err != nil {
			c.JSON(500, gin.H{"msg": "获取objects失败"})
			return
		}
		// 获取论文总数
		total, err := service.GetScholarWorksCount(author_id)
		if err != nil {
			c.JSON(404, gin.H{"msg": "查询论文总数出错，修改失败"})
			return
		}
		data := make([]map[string]interface{}, len(works))
		_works := objects["results"].([]interface{})
		if objects != nil {
			//for i, v := range objects {
			//	if v.Found {
			//		json.Unmarshal(v.Source, &data[i])
			//		data[i]["Top"] = works[i].Top
			//		data[i]["find"] = true
			//	} else {
			//		data[i] = make(map[string]interface{})
			//		data[i]["find"] = false
			//		println(works_ids[i] + " not found")
			//	}
			//}
			workfilter := utils.InitWorksfilter()
			for i, v := range works_ids {
				for j, _v := range _works {
					work := _v.(map[string]interface{})
					id := work["id"].(string)
					if v == utils.RemovePrefix(id) {
						utils.FilterData(&work, &workfilter)
						if work["abstract_inverted_index"] != nil {
							work["abstract"] = utils.TransInvertedIndex2String(work["abstract_inverted_index"].(map[string]interface{}))
							work["abstract_inverted_index"] = nil
						}
						data[i] = work
						data[i]["Top"] = works[i].Top
						data[i]["find"] = true
						data[i]["pdf"] = works[i].PDF
						if works[i].PDF != "" {
							data[i]["isupdatepdf"] = 1
						} else {
							data[i]["isupdatepdf"] = 0
						}
						_works = append(_works[:j], _works[j+1:]...)
						break
					}
				}
			}
		}
		c.JSON(200, gin.H{"msg": "获取成功", "data": data, "pages": pages, "total": total})
		return
	} else {
		// 不能找到则从openalex api中获取
		log.Println("从openalex api中获取author works")
		//author := res[]
		//var author_map map[string]interface{}
		//_ = json.Unmarshal(author, &author_map)
		works_api_url := res["works_api_url"].(string)
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
		total := len(works)
		// 获取论文总数
		count, err := service.GetScholarWorksCount(author_id)
		if count == 0 || err != nil { // 之前没有插入过
			service.UpdateScholarWorksCount(author_id, total)
			go service.CreateWorks(works)
		}
		// 页数从1开始
		if page == pages { // 最后一页
			works = works[(page-1)*page_size:]
		} else {
			works = works[(page-1)*page_size : page*page_size]
		}
		// 获取works_id
		for _, work := range works {
			works_ids = append(works_ids, work.WorkID)
		}
		objects, err := service.GetObjects2("works", works_ids)
		if err != nil {
			c.JSON(500, gin.H{"msg": "获取objects失败"})
			return
		}
		data := make([]map[string]interface{}, len(works))
		_works := objects["results"].([]interface{})
		if objects != nil {
			//for i, v := range objects.Docs {
			//	if v.Found {
			//		json.Unmarshal(v.Source, &data[i])
			//		data[i]["Top"] = works[i].Top
			//		data[i]["find"] = true
			//	} else {
			//		data[i] = make(map[string]interface{})
			//		data[i]["find"] = false
			//		println(works_ids[i] + " not found")
			//	}
			//}
			workfilter := utils.InitWorksfilter()
			for i, v := range works_ids {
				for j, _v := range _works {
					work := _v.(map[string]interface{})
					id := work["id"].(string)
					if v == utils.RemovePrefix(id) {
						utils.FilterData(&work, &workfilter)
						if work["abstract_inverted_index"] != nil {
							work["abstract"] = utils.TransInvertedIndex2String(work["abstract_inverted_index"].(map[string]interface{}))
							work["abstract_inverted_index"] = nil
						}
						data[i] = work
						data[i]["Top"] = works[i].Top
						data[i]["find"] = true
						data[i]["pdf"] = works[i].PDF
						if works[i].PDF != "" {
							data[i]["isupdatepdf"] = 1
						} else {
							data[i]["isupdatepdf"] = 0
						}
						_works = append(_works[:j], _works[j+1:]...)
						break
					}
				}
			}
		}
		c.JSON(200, gin.H{"msg": "获取成功", "data": data, "pages": pages, "total": total})
	}
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
// @Param       data body     response.IgnoreWorkQ true "data 是请求参数,包括author_id ,work_id"
// @Success     200  {string} json                 "{"msg":"修改忽略属性成功"}"
// @Failure     400  {string} json                 "{"msg":"参数错误"}"
// @Failure     401  {string} json                 "{"msg":"修改忽略属性失败"}"
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
// @Param       data body     response.ModifyPlaceQ true "data 是请求参数,包括author_id ,work_id ,direction"
// @Success     200  {string} json                  "{"msg":"修改成功"}"
// @Failure     400  {string} json                  "{"msg":"参数错误"}"
// @Failure     401  {string} json                  "{"msg":"未找到该论文"}"
// @Failure     402  {string} json                  "{"msg":"论文已经在顶部"}"
// @Failure     403  {string} json                  "{"msg":"论文已经在底部"}"
// @Failure     404  {string} json                  "{"msg":"修改失败"}"
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
// @Param       data body     response.TopWorkQ true "data 是请求参数,包括author_id ,work_id"
// @Success     200  {string} json              "{"msg":"置顶成功"}"
// @Failure     400  {string} json              "{"msg":"参数错误"}"
// @Failure     401  {string} json              "{"msg":"未找到该论文"}"
// @Failure     402  {string} json              "{"msg":"修改失败"}"
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

// 取消置顶论文
// @Summary     学者管理主页--取消置顶论文 hr
// @Description 学者管理主页--取消置顶论文
// @Description
// @Description 参数说明
// @Description - author_id 作者的id
// @Description
// @Description - work_id 论文的id
// @Tags        学者主页的论文获取、管理
// @Accept      json
// @Produce     json
// @Param       data body     response.TopWorkQ true "data 是请求参数,包括author_id ,work_id"
// @Success     200  {string} json              "{"msg":"取消置顶成功"}"
// @Failure     400  {string} json                  "{"msg":"参数错误"}"
// @Failure     401  {string} json                  "{"msg":"未找到该论文"}"
// @Failure     402  {string} json              "{"msg":"修改失败"}"
// @Router      /scholar/works/untop [POST]
func UnTopWork(c *gin.Context) {
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
	err := service.UnTopWork(author_id, work_id)
	if err != nil {
		c.JSON(402, gin.H{"msg": "修改失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "取消置顶成功"})
}

// UploadAuthorHeadshot
// @Summary     上传作者头像 txc
// @Description 上传作者头像
// @Tags        scholar
// @Param       token     header   string true "token"
// @Param       author_id formData string true "学者ID"
// @Param       Headshot  formData file   true "新头像"
// @Success     200       {string} json   "{"msg":"上传成功","data": author}"
// @Failure     400       {string} json   "{"msg":"学者未被认领"}"
// @Failure     401       {string} json   "{"msg":"无权限"}"
// @Failure     402       {string} json   "{"msg":"头像文件获取失败"}"
// @Router      /scholar/author/headshot [POST]
func UploadAuthorHeadshot(c *gin.Context) {
	authorID := c.Request.FormValue("author_id")
	author, notFound := service.GetAuthor(authorID)
	if notFound {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "学者未被认领"})
		return
	}
	user := c.MustGet("user").(database.User)
	if user.AuthorID != authorID {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "无权限"})
		return
	}
	file, err := c.FormFile("Headshot")
	if err != nil {
		c.JSON(402, gin.H{"msg": "头像文件获取失败"})
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
}

// ModifyAuthorIntro
// @Summary     txc
// @Description 修改作者简介
// @Tags        scholar
// @Accept      json
// @Produce     json
// @Param       token header   string true "token"
// @Param       data  body     response.ModifyAuthorIntroQ true "data"
// @Success     200   {string} json                        "{"msg":"修改成功"}"
// @Failure     400   {string} json                        "{"msg":"参数错误"}"
// @Failure     401   {string} json                        "{"msg":"无权限"}"
// @Failure     404   {string} json                        "{"msg":"学者未被认领"}"
// @Router      /scholar/author/intro [POST]
func ModifyAuthorIntro(c *gin.Context) {
	var d response.ModifyAuthorIntroQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	author, notFound := service.GetAuthor(d.AuthorID)
	if notFound {
		c.JSON(404, gin.H{"msg": "学者未被认领"})
		return
	}
	user := c.MustGet("user").(database.User)
	if user.AuthorID != d.AuthorID {
		c.JSON(401, gin.H{"msg": "无权限"})
		return
	}
	author.Intro = d.Intro
	err := global.DB.Save(author).Error
	if err != nil {
		c.JSON(402, gin.H{"msg": "数据库修改失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "修改成功"})
}

// 上传作品PDF
// @Summary     学者管理主页--上传作品PDF hr
// @Description 学者管理主页--上传作品PDF
// @Description
// @Description 参数说明
// @Description - author_id 作者的id
// @Description
// @Description - work_id 论文的id
// @Description
// @Description - PDF 上传的PDF文件
// @Tags        学者主页的论文获取、管理
// @Param       author_id formData string true "学者ID"
// @Param       work_id   formData string true "论文ID"
// @Param       PDF       formData file   true "PDF"
// @Router      /scholar/works/upload [POST]
func UploadPaperPDF(c *gin.Context) {
	authorID := c.Request.FormValue("author_id")
	workID := c.Request.FormValue("work_id")
	log.Println("authorID: ", authorID)
	log.Println("workID: ", workID)
	_, notFound := service.GetPersonalWork(authorID, workID)
	if notFound {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "论文不存在"})
		return
	}
	pdf, err := c.FormFile("PDF")
	if err != nil {
		c.JSON(401, gin.H{"msg": "文件上传失败"})
		return
	}
	raw := fmt.Sprintf("%d", authorID) + time.Now().String() + pdf.Filename
	md5 := utils.GetMd5(raw)
	suffix := strings.Split(pdf.Filename, ".")[1]
	saveDir := "./media/pdf/"
	saveName := md5 + "." + suffix
	log.Printf("saveName: %s", saveName)
	savePath := path.Join(saveDir, saveName)
	if err = c.SaveUploadedFile(pdf, savePath); err != nil {
		c.JSON(402, gin.H{"msg": "文件保存失败"})
		return
	}
	err = service.UpdateWorkPdf(authorID, workID, saveName)
	if err != nil {
		c.JSON(403, gin.H{"msg": "保存文件路径到数据库中失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "上传成功", "data": saveName})
}

// UnUploadPaperPDF
// @Summary     学者管理主页--取消上传作品PDF txc
// @Description 学者管理主页--取消上传作品PDF
// @Description
// @Description 参数说明
// @Description - author_id 作者的id
// @Description
// @Description - work_id 论文的id
// @Tags        学者主页的论文获取、管理
// @Param       author_id formData string true "学者ID"
// @Param       work_id   formData string true "论文ID"
// @Success     200       {string} json   "{"msg":"取消上传成功"}"
// @Failure     400       {string} json   "{"msg":"论文不存在"}"
// @Failure     403       {string} json   "{"msg":"保存文件路径到数据库中失败"}"
// @Router      /scholar/works/unupload [POST]
func UnUploadPaperPDF(c *gin.Context) {
	authorID := c.Request.FormValue("author_id")
	workID := c.Request.FormValue("work_id")
	_, notFound := service.GetPersonalWork(authorID, workID)
	if notFound {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "论文不存在"})
		return
	}
	err := service.UpdateWorkPdf(authorID, workID, "")
	if err != nil {
		c.JSON(403, gin.H{"msg": "保存文件路径到数据库中失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "取消上传成功"})
}

// 获取论文PDF地址
// @Summary     获取学者上传的文章PDF地址 hr
// @Description 获取学者上传的文章PDF地址
// @Description
// @Description 参数说明
// @Description - work_id 论文的id
// @Description
// @Description 返回说明
// @Description - pdf地址,直接使用即可，无需拼接
// @Tags        学者主页的论文获取、管理
// @Accept      json
// @Produce     json
// @Param       data body     response.GetPaperPDFQ true "data 是请求参数work_id"
// @Success     200  {string} json                  "{"msg":"获取成功", "data": "pdf地址"}"
// @Failure     400  {string} json              "{"msg":"参数错误"}"
// @Failure     401  {string} json              "{"msg":"未找到该论文"}"
// @Failure     402  {string} json                  "{"msg":"未上传PDF"}"
// @Router      /scholar/works/getpdf [POST]
func GetPaperPDF(c *gin.Context) {
	var d response.GetPaperPDFQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(400, gin.H{"msg": "参数错误"})
		return
	}
	works, err := service.GetWorksByWorkID(d.WorkID)
	if err != nil {
		c.JSON(401, gin.H{"msg": "出现err,未找到该论文", "err": err})
		return
	}
	if len(works) == 0 {
		c.JSON(401, gin.H{"msg": "未找到该论文"})
		return
	}
	// log.Println(works)
	for _, work := range works {
		// log.Println(work)
		if work.PDF != "" {
			url := "http://ishare.horik.cn:8000/api/media/pdf/" + work.PDF
			c.JSON(200, gin.H{"msg": "获取成功", "data": url})
			return
		}
	}
	c.JSON(402, gin.H{"msg": "未上传PDF"})
}
