package v1

import (
	"IShare/global"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"context"
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
// @Param       id  query     string true "id"
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
// @Param       data body response.BaseSearchQ true "搜索条件"
// @Success     200        {string} json   "{"status":200,"res":{obeject}}"
// @Failure     200        {string} json   "{"status":201,"err":"es search err"}"
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
	res, err := service.CommonWorkSearch(boolQuery, d.Page, d.Size, 0, false, aggs)
	if err != nil {
		c.JSON(200, gin.H{
			"status": 201,
			"msg":    "es search err",
			"err":    err,
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"res":    res,
	})
}

// BaseSearch2
// @Summary     txc
// @Description 基本搜索，Cond里面填筛选条件，key仅包含["type", "author", "institution", "publisher", "venue", "publication_year"]
// @Tags        esSearch
// @Accept      json
// @Produce     json
// @Param       data body response.BaseSearchQ true "搜索条件"
// @Success     200        {string} json   "{"status":200,"res":{obeject}}"
// @Failure     200        {string} json   "{"status":201,"err":"es search err"}"
// @Router      /es/search/base2 [POST]
func BaseSearch2(c *gin.Context) {
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
	res, err := service.CommonWorkSearch(boolQuery, d.Page, d.Size, 0, false, aggs)
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
	//for k, v := range res.Aggregations {
	//	data.Aggs[k] = v
	//}
	c.JSON(200, gin.H{
		"status": 200,
		"res":    data,
	})
}

// AdvanceSearch
// @Description 高级搜索，搜索条件通过body传入，未完成
// @Router      /es/search/advance [POST]
func AdvanceSearch(c *gin.Context) {
}

// DoiSearch
// @Summary     txc
// @Description 使用doi查找work，未测试，请勿使用
// @Tags 	  esSearch
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
