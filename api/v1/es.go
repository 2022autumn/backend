package v1

import (
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"encoding/json"
	"net/http"

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

// GetObject
// @Summary     txc
// @Description 根据id获取对象，可以是author，work，institution,venue,concept
// @Tags        esSearch
// @Param       id  query    string true "id"
// @Success     200 {string} json   "{"status":200,"res":{obeject}}"
// @Failure     201 {string} json   "{"status":201,"msg":"es get err"}"
// @Failure     400 {string} json   "{"status":400,"msg":"id type error"}"
// @Router      /es/get/ [GET]
func GetObject(c *gin.Context) {
	id := c.Query("id")
	idx, err := utils.TransObjPrefix(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": 400,
			"msg":    "id type error",
		})
		return
	}
	res, err := service.GetObject(idx, id)
	if err != nil {
		c.JSON(201, gin.H{
			"msg":    "es get err",
			"status": 201,
		})
		return
	}
	if idx == "works" && res.Hits.TotalHits.Value == 1 {
		var tmp = make(map[string]interface{})
		by, _ := res.Hits.Hits[0].Source.MarshalJSON()
		_ = json.Unmarshal(by, &tmp)
		referenced_works := tmp["referenced_works"].([]interface{})
		var newReferencedWorks []map[string]string
		for _, v := range referenced_works {
			res, _ := service.GetObject("works", v.(string))
			if res.Hits.TotalHits.Value == 1 {
				newReferencedWorks = append(newReferencedWorks, map[string]string{
					"id":    v.(string),
					"cited": GetWorkCited(res.Hits.Hits[0].Source),
				})
			}
		}
		tmp["referenced_works"] = newReferencedWorks
		related_works := tmp["related_works"].([]interface{})
		var newRelatedWorks []map[string]string
		for _, v := range related_works {
			res, _ := service.GetObject("works", v.(string))
			if res.Hits.TotalHits.Value == 1 {
				newRelatedWorks = append(newRelatedWorks, map[string]string{
					"id":    v.(string),
					"cited": GetWorkCited(res.Hits.Hits[0].Source),
				})
			}
		}
		tmp["related_works"] = newRelatedWorks
		c.JSON(http.StatusOK, gin.H{
			"data":   tmp,
			"status": 200,
		})
		return
	}
	var data = response.GetObjectA{
		RawMessage: res.Hits.Hits[0].Source,
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   data,
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
