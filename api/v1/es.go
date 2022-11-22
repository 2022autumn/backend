package v1

import (
	"IShare/global"
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
	nameQuery := elastic.NewTermQuery("name", queryWord)
	infoQuery := elastic.NewMatchPhraseQuery("info", queryWord)
	boolQuery.Should(nameQuery, infoQuery)
	searchRes, err := global.ES.Search().
		Index("students").
		Query(boolQuery).
		Do(context.Background())
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"success": false, "message": "参数错误", "status": 401})
		panic(fmt.Errorf("es search err"))
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"res":     searchRes,
	})
}
