package v1

import (
	"IShare/global"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

// TestEsSearch
// @Param  queryWord formData string true "queryWord"
// @Router /es/test_es [POST]
func TestEsSearch(c *gin.Context) {
	queryWord := c.Request.FormValue("queryWord")
	queryWord = strings.ToLower(queryWord)
	boolQuery := elastic.NewBoolQuery()
	// nameQuery := elastic.NewTermQuery("name", queryWord)
	infoQuery := elastic.NewMatchPhraseQuery("authors.name", queryWord)
	boolQuery.Should(infoQuery)
	age_agg := elastic.NewTermsAggregation().Field("info.keyword")
	searchRes, err := global.ES.Search().
		Index("students").
		Aggregation("nameless", age_agg).
		Query(boolQuery).
		Do(context.Background())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "参数错误", "status": 401})
		panic(fmt.Errorf("es search err"))
	}
	c.JSON(http.StatusOK, gin.H{
		"res":    searchRes,
		"status": 200,
	})
}

// GetObject
// @Summary     txc
// @Description 根据id获取对象，可以是author，work，institution,venue,concept
// @Tags        esSearch
// @Param       id  query    string true "id"
// @Success     200 {string} json   "{"status":200,"res":{obeject}}"
// @Failure     200 {string} json   "{"status":201,"msg":"es get err"}"
// @Failure     200 {string} json   "{"status":201,"msg":"id type error"}"
// @Router      /es/get/ [GET]
func GetObject(c *gin.Context) {
	id := c.Query("id")
	idx, err := utils.TransObjPrefix(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": 202,
			"msg":    "id type error",
		})
		return
	}
	res, err := service.GetObject(idx, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg":    "es get err",
			"status": 201,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"res":    res,
		"status": 200,
	})
}

// BaseSearch
// @Summary     txc
// @Description 基本搜索，Cond里面填筛选条件，key仅包含["type", "author", "institution", "publisher", "venue", "publication_year"]
// @Tags        esSearch
// @Accept      json
// @Produce     json
// @Param       data body     response.BaseSearchQ true "搜索条件"
// @Success     200  {string} json                 "{"status":200,"res":{obeject}}"
// @Failure     200  {string} json                 "{"status":201,"err":"es search err"}"
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
		switch k {
		case "type":
			aggs["types"] = false
			boolQuery.Filter(elastic.NewMatchQuery("type.keyword", v))
		case "institution":
			aggs["institutions"] = false
			boolQuery.Filter(elastic.NewMatchQuery("authorships.institutions.display_name.keyword", v))
		case "publisher":
			aggs["publishers"] = false
			boolQuery.Filter(elastic.NewMatchQuery("host_venue.publisher.keyword", v))
		case "venue":
			aggs["venues"] = false
			boolQuery.Filter(elastic.NewMatchQuery("host_venue.display_name.keyword", v))
		case "author":
			aggs["authors"] = false
			boolQuery.Filter(elastic.NewMatchQuery("authorships.author.display_name.keyword", v))
		case "publication_year":
			aggs["publication_years"] = false
			boolQuery.Filter(elastic.NewMatchQuery("publication_year", v))
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
	var data = response.BaseSearchA{Hits: res.Hits.TotalHits.Value}
	for _, v := range res.Hits.Hits {
		data.Works = append(data.Works, v.Source)
	}
	data.Aggs = make(map[string]interface{})
	var tmp = make(map[string]interface{})
	for k, v := range res.Aggregations {
		by, _ := v.MarshalJSON()
		_ = json.Unmarshal(by, &tmp)
		data.Aggs[k] = tmp["buckets"].([]interface{})
	}
	c.JSON(200, gin.H{
		"status": 200,
		"res":    data,
	})
}

var cond2field = map[string]string{
	"type":             "type.keyword",
	"author":           "authorships.author.display_name.keyword",
	"institution":      "authorships.institutions.display_name.keyword",
	"publisher":        "host_venue.publisher.keyword",
	"venue":            "host_venue.display_name.keyword",
	"publication_year": "publication_years",
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
			aggs[kk] = false
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
	var data = response.BaseSearchA{Hits: res.Hits.TotalHits.Value}
	for _, v := range res.Hits.Hits {
		data.Works = append(data.Works, v.Source)
	}
	data.Aggs = make(map[string]interface{})
	var tmp = make(map[string]interface{})
	for k, v := range res.Aggregations {
		by, _ := v.MarshalJSON()
		_ = json.Unmarshal(by, &tmp)
		data.Aggs[k] = tmp["buckets"].([]interface{})
	}
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
