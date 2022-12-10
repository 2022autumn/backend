package v1

import (
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

func GetWorkCited(w json.RawMessage) string {
	var work = make(map[string]interface{})
	_ = json.Unmarshal(w, &work)
	var cited string
	for _, v := range work["authorships"].([]interface{}) {
		authorship := v.(map[string]interface{})
		if authorship["author_position"] == "first" {
			author := authorship["author"].(map[string]interface{})
			cited += author["display_name"].(string) + ", "
		}
	}
	cited += "\"" + work["title"].(string) + "\""
	if work["host_venue"] != nil {
		if work["host_venue"].(map[string]interface{})["display_name"] != nil {
			cited += "," + work["host_venue"].(map[string]interface{})["display_name"].(string)
		}
	}
	cited += "."
	return cited
}

func TransRefs2Cited(refs []interface{}) []map[string]string {
	var newReferencedWorks = make([]map[string]string, 0)
	var ids []string
	for _, v := range refs {
		ids = append(ids, v.(string))
	}
	works, _ := service.GetObjects("works_v1", ids)
	if works != nil {
		for i, v := range works.Docs {
			if v.Found == true {
				newReferencedWorks = append(newReferencedWorks, map[string]string{
					"id":    ids[i],
					"cited": GetWorkCited(v.Source),
				})
			} else {
				println(ids[i] + " not found")
			}
		}
	}
	return newReferencedWorks
}

// GetObject
// @Summary     txc
// @Description 根据id获取对象，可以是author，work，institution,venue,concept W4237558494,W2009180309,W2984203759
// @Tags    esSearch
// @Param       id     query    string true  "对象id"
// @Param       userid query    string false "用户id"
// @Success     200    {string} json   "{"status":200,"res":{}}"
// @Failure     404    {string} json   "{"status":201,"msg":"es get err or not found"}"
// @Failure     400    {string} json   "{"status":400,"msg":"id type error"}"
// @Router      /es/get/ [GET]
func GetObject(c *gin.Context) {
	id := c.Query("id")
	userid := c.Query("userid")
	idx, err := utils.TransObjPrefix(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "id type error"})
		return
	}
	res, err := service.GetObject(idx, id)
	if err != nil {
		c.JSON(404, gin.H{"msg": "es get err or not found"})
		return
	}
	var tmp = make(map[string]interface{})
	_ = json.Unmarshal(res.Source, &tmp)
	if userid != "" {
		userid, _ := strconv.ParseUint(userid, 0, 64)
		ucs, err := service.GetUserConcepts(userid)
		if err != nil {
			c.JSON(405, gin.H{"msg": "get user concepts err"})
			return
		}
		if tmp["concepts"] != nil {
			concepts := tmp["concepts"].([]interface{})
			for _, c := range concepts {
				concept := c.(map[string]interface{})
				conceptid := concept["id"].(string)
				var flag = false
				for i, uc := range ucs {
					if conceptid == uc.ConceptID {
						concept["islike"] = true
						flag = true
						ucs = append(ucs[:i], ucs[i+1:]...)
						break
					}
				}
				if flag == false {
					concept["islike"] = false
				}
			}
		}
	}
	if idx == "works_v1" {
		referenced_works := tmp["referenced_works"].([]interface{})
		tmp["referenced_works"] = TransRefs2Cited(referenced_works)
		related_works := tmp["related_works"].([]interface{})
		tmp["related_works"] = TransRefs2Cited(related_works)
		wv, notFound := service.GetWorkView(id)
		if notFound {
			wv = database.WorkView{
				WorkID:    id,
				Views:     1,
				WorkTitle: tmp["title"].(string),
			}
			err := service.CreateWorkView(&wv)
			if err != nil {
				println("create work view err")
			}
		} else {
			wv.Views += 1
			err := service.SaveWorkView(&wv)
			if err != nil {
				println("save work view err")
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   tmp,
		"status": 200,
	})
}

var cond2field = map[string]string{
	"type":             "type.keyword",
	"author":           "authorships.author.display_name.keyword",
	"institution":      "authorships.institutions.display_name.keyword",
	"publisher":        "host_venue.publisher.keyword",
	"venue":            "host_venue.display_name.keyword",
	"publication_year": "publication_year",
}
var query2field = map[string]string{
	"title":       "title",
	"abstract":    "abstract",
	"venue":       "host_venue.display_name",
	"publisher":   "host_venue.publisher",
	"author":      "authorships.author.display_name",
	"institution": "authorships.institutions.display_name",
	"concept":     "concepts.display_name",
}

// BaseSearch
// @Summary     txc
// @Description 基本搜索，Cond里面填筛选条件，key仅包含["type", "author", "institution", "publisher", "venue", "publication_year"]
// @Tags        esSearch
// @Accept      json
// @Produce     json
// @Param       data body     response.BaseSearchQ true "搜索条件"
// @Success     200  {string} json                 "{"status":200,"res":{obeject}}"
// @Failure     201  {string} json                 "{"status":201,"err":"es search err"}"
// @Router      /es/search/base [POST]
func BaseSearch(c *gin.Context) {
	var d response.BaseSearchQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	boolQuery := elastic.NewBoolQuery()
	tiQuery := elastic.NewMatchPhraseQuery("title", d.QueryWord)
	abQuery := elastic.NewMatchPhraseQuery("abstract", d.QueryWord)
	b2Query := elastic.NewBoolQuery()
	b2Query.Should(tiQuery, abQuery)
	boolQuery.Must(b2Query)
	var aggs = make(map[string]bool)
	var aggList = [6]string{"types", "authors", "institutions", "publishers", "venues", "publication_years"}
	for _, k := range aggList {
		aggs[k] = true
	}
	for k, v := range d.Conds {
		if kk, ok := cond2field[k]; ok {
			aggs[k+"s"] = false
			boolQuery.Filter(elastic.NewMatchQuery(kk, v))
		}
	}
	res, err := service.CommonWorkSearch(boolQuery, d.Page, d.Size, d.Sort, d.Asc, aggs)
	if err != nil {
		c.JSON(201, gin.H{
			"status": 201,
			"msg":    "es search err",
			"err":    err,
		})
		return
	}
	var data = response.BaseSearchA{}
	data.Hits, data.Works, data.Aggs, _ = utils.NormalizationSearchResult(res)
	c.JSON(http.StatusOK, gin.H{
		"status": 200,
		"res":    data,
	})
}

// AdvancedSearch
// @Summary     txc
// @Description 高级搜索，query是一个map列表， 每个map包含"content" "field" "logic"
// @Description logic 仅包含["and", "or", "not"]
// @Description field 仅包含["title", "abstract", "venue", "publisher", "author", "institution", "concept"]
// @Description 对于年份的筛选，在query里面 field是"publication_date" logic默认为and， 该map下有"begin" "end"分别是开始和结束
// @Description sort=0为默认排序（降序） =1为按引用数降序 =2按发表日期由近到远
// @Description asc=0为降序 =1为升序
// @Description { "asc": false,"conds": {"venue":"International Journal for Research in Applied Science and Engineering Technology","author": "Zenith Nandy"},"page": 1,"query": [{"field": "title","content": "python","logic": "and"},{"field": "publication_date","begin": "2021-12-01","end":"2022-06-01","logic": "and"}],"size": 8,"sort": 0}
// @Tags        esSearch
// @Accept      json
// @Produce     json
// @Param       data body response.AdvancedSearchQ true "data"
// @Router      /es/search/advanced [POST]
func AdvancedSearch(c *gin.Context) {
	// author title abstract venue institution publisher publication_year concept
	var d response.AdvancedSearchQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	boolQuery := elastic.NewBoolQuery()
	subQuery := elastic.NewBoolQuery()
	for _, i := range d.Query {
		if i["logic"] == "and" {
			if i["field"] == "publication_date" {
				subQuery.Must(elastic.NewRangeQuery("publication_date").Gte(i["begin"]).Lte(i["end"]))
			} else {
				subQuery.Must(elastic.NewMatchPhraseQuery(query2field[i["field"]], i["content"]))
			}
		} else if i["logic"] == "or" {
			subQuery.Should(elastic.NewMatchPhraseQuery(query2field[i["field"]], i["content"]))
		} else if i["logic"] == "not" {
			subQuery.MustNot(elastic.NewMatchPhraseQuery(query2field[i["field"]], i["content"]))
		}
	}
	boolQuery.Must(subQuery)
	var aggs = make(map[string]bool)
	var aggList = [6]string{"types", "authors", "institutions", "publishers", "venues", "publication_years"}
	for _, k := range aggList {
		aggs[k] = true
	}
	for k, v := range d.Conds {
		if kk, ok := cond2field[k]; ok {
			aggs[k+"s"] = false
			boolQuery.Filter(elastic.NewMatchQuery(kk, v))
		}
	}
	res, err := service.CommonWorkSearch(boolQuery, d.Page, d.Size, d.Sort, d.Asc, aggs)
	if err != nil {
		c.JSON(200, gin.H{
			"status": 201,
			"msg":    "es search err",
			"err":    err,
		})
		return
	}
	var data = response.BaseSearchA{}
	data.Hits, data.Works, data.Aggs, _ = utils.NormalizationSearchResult(res)
	c.JSON(200, gin.H{
		"status": 200,
		"res":    data,
	})
}

// DoiSearch
// @Summary     txc
// @Description 使用doi查找work，未测试，请勿使用
// @Tags        esSearch
// @Param       doi query string true "doi"
// @Router      /es/search/doi [POST]
func DoiSearch(c *gin.Context) {
	doi := c.Query("doi")
	boolQuery := elastic.NewBoolQuery()
	doiQuery := elastic.NewMatchQuery("doi", doi)
	boolQuery.Must(doiQuery)
	res, err := service.GetWork(boolQuery)
	if err != nil {
		c.JSON(200, gin.H{
			"status": 201,
			"err":    err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"res":    res,
	})
}

// GetAuthorRelationNet
// @Summary     hr
// @Description 根据author的id获取专家关系网络, 目前会返回Top N的关系网，N=10，后续可以讨论修改N的大小或者传参给我
// @Description
// @Description 目前接口时延约为1s, 后续考虑把计算出来的结果存入数据库，二次查询时延降低
// @Description
// @Tags        esSearch
// @Param   author_id query    string                     true "author_id" Enums(A2764814280, A2900471938, A2227665069)
// @Success 200       {object} response.AuthorRelationNet "{"res":{ "Vertex_set":[], "Edge_set":[]}}"
// @Failure 201       {string} json                       "{"msg":"Get Author Relation Net Error"}"
// @Router  /es/getAuthorRelationNet [GET]
func GetAuthorRelationNet(c *gin.Context) {
	author_id := c.Query("author_id")
	var err error
	data := response.AuthorRelationNet{}
	data.Vertex_set, data.Edge_set, err = service.GetAuthorRelationNet(author_id)
	if err != nil {
		c.JSON(201, gin.H{
			"msg": "Get Author Relation Net Error",
			"err": err,
		})
		return
	}
	c.JSON(200, gin.H{
		"res": data,
	})
}

func GetWorksOfAuthorByUrl(c *gin.Context) {
	var d response.BaseSearchQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	c.JSON(200, gin.H{
		"status": 200,
	})
}

// GetStatistics
// @Summary     txc
// @Description 获取统计信息
// @Tags        esSearch
// @Success     200 {string} json "{"res":{}}"
// @Failure     301 {string} json "{"err":{}}"
// @Router      /es/statistic [GET]
func GetStatistics(c *gin.Context) {
	res, err := service.GetStatistics()
	if err != nil {
		c.JSON(301, gin.H{"err": err})
		return
	}
	c.JSON(200, gin.H{"res": res})
}

// GetPrefixSuggestion doc
// @Summary     hr
// @description 根据前缀得到搜索建议，返回results 字符串数组
// @Tags esSearch
// @Param name query string true "name 表示需要查询的字段名"
// @Param prefix query string true "prefix 表示用户已经输入的前缀"
// @Success 200 {string} string "{"success": true, "msg": "获取成功"}"
// @Failure 400 {string} string "{"success": false, "msg": 参数错误"}"
// @Failure 402 {string} string "{"success": false, "msg": "es服务出错"}"
// @Router /es/prefix [POST]
func GetPrefixSuggestion(c *gin.Context) {
	var d response.PrefixSuggestionQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "参数错误"})
		panic(err)
	}
	index, field, prefix, topN := "work_v1", d.Prefix, d.Field, 10
	prefixResult, err := service.PrefixSearch(index, field, prefix, topN)
	if err != nil {
		c.JSON(402, gin.H{"success": false, "msg": "es服务出错"})
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "成功", "res": prefixResult})
}
