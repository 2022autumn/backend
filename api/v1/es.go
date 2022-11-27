package v1

import (
	"IShare/global"
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
// @Param queryWord formData string true "queryWord"
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
// @Param id path string true "id"
// @Success 200
// @Router /es/get/{id} [GET]
func GetObject(c *gin.Context) {
	id := c.Param("id")
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
		"serachResult": res,
		"status":       200,
	})
}

// BaseSearch
// @Description 基本搜索，根据query字段去查找title和abstract里面含有搜索词的work，词是精确查找
// @Param query_word query string true "搜索词"
// @Router /es/base_query [POST]
func BaseSearch(c *gin.Context) {
	queryWord := c.Query("query_word")
	boolQuery := elastic.NewBoolQuery()
	tiQuery := elastic.NewMatchPhraseQuery("title", queryWord)
	abQuery := elastic.NewMatchPhraseQuery("abstract", queryWord)
	boolQuery.Should(tiQuery, abQuery)
	res, err := service.CommonWorkSearch(1, 10, boolQuery, 0, false)
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

func AdvanceSearch(c *gin.Context) {

}
