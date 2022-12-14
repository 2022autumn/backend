package v1

import (
	"IShare/model/database"
	"IShare/model/response"
	"IShare/service"
	"IShare/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

func GetWorkCited(w json.RawMessage) string {
	var work = make(map[string]interface{})
	_ = json.Unmarshal(w, &work)
	var cited string
	for i, v := range work["authorships"].([]interface{}) {
		authorship := v.(map[string]interface{})
		if i <= 2 {
			author := authorship["author"].(map[string]interface{})
			cited += author["display_name"].(string) + ", "
		} else {
			break
		}
	}
	cited += "\"" + work["title"].(string) + "\""
	if work["host_venue"] != nil {
		if work["host_venue"].(map[string]interface{})["display_name"] != nil {
			cited += "," + work["host_venue"].(map[string]interface{})["display_name"].(string)
		}
	}
	cited += "," + strconv.Itoa(int(work["publication_year"].(float64))) + "."
	return cited
}
func GetWorkAuthorNames(w map[string]interface{}) (names []string) {
	for _, v := range w["authorships"].([]interface{}) {
		authorship := v.(map[string]interface{})
		author := authorship["author"].(map[string]interface{})
		names = append(names, author["display_name"].(string))
	}
	return
}
func GenAPACited(work map[string]interface{}) string {
	var authorNames = GetWorkAuthorNames(work)
	var cited string
	if len(authorNames) > 6 {
		cited += authorNames[0] + " et al."
	} else {
		for i, v := range authorNames {
			if i == len(authorNames)-2 {
				cited += v + ", & "
			} else if i == len(authorNames)-1 {
				cited += v
			} else {
				cited += v + ", "
			}
		}
	}
	cited += " (" + strconv.Itoa(int(work["publication_year"].(float64))) + "). "
	cited += work["title"].(string) + ". "
	if work["host_venue"] != nil {
		if work["host_venue"].(map[string]interface{})["display_name"] != nil {
			cited += work["host_venue"].(map[string]interface{})["display_name"].(string) + ". "
		}
	}
	return cited
}
func GenMLACited(work map[string]interface{}) string {
	var authorNames = GetWorkAuthorNames(work)
	var cited string
	if len(authorNames) > 3 {
		cited += authorNames[0] + " et al."
	} else {
		cited += strings.Join(authorNames, ", ")
	}
	cited += ". \"" + work["title"].(string) + ".\" "
	if work["host_venue"] != nil {
		if work["host_venue"].(map[string]interface{})["display_name"] != nil {
			cited += work["host_venue"].(map[string]interface{})["display_name"].(string) + ", "
		}
	}
	cited += strconv.Itoa(int(work["publication_year"].(float64))) + ", "
	return cited
}
func GenGBCited(work map[string]interface{}) string {
	var authorNames = GetWorkAuthorNames(work)
	var cited string
	if len(authorNames) > 3 {
		cited += authorNames[0] + ", " + authorNames[1] + ", " + authorNames[2] + ", et al."
	} else {
		cited += strings.Join(authorNames, ", ")
	}
	cited += ". " + work["title"].(string) + ". "
	if work["type"] != nil {
		if strings.Index(work["type"].(string), "journal") != -1 {
			cited += "[J]. "
			if work["host_venue"] != nil {
				if work["host_venue"].(map[string]interface{})["display_name"] != nil {
					cited += work["host_venue"].(map[string]interface{})["display_name"].(string) + ", "
				}
				cited += strconv.Itoa(int(work["publication_year"].(float64))) + "."
			}
		} else if strings.Index(work["type"].(string), "conference") != -1 {
			cited += "[C]. "
		} else if strings.Index(work["type"].(string), "book") != -1 {
			cited += "[M]. "
			if work["host_venue"] != nil {
				if work["host_venue"].(map[string]interface{})["publisher"] != nil {
					cited += work["host_venue"].(map[string]interface{})["publisher"].(string) + ", "
				}
				cited += strconv.Itoa(int(work["publication_year"].(float64))) + "."
			}
		}
	}
	return cited
}
func TransRefs2Cited(refs []interface{}) []map[string]string {
	var newReferencedWorks = make([]map[string]string, 0)
	var ids []string
	for _, v := range refs {
		ids = append(ids, v.(string))
	}
	works, _ := service.GetObjects("works", ids)
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
func TransRefs2Intro(refs []interface{}) []map[string]interface{} {
	var newReferencedWorks = make([]map[string]interface{}, 0)
	var ids []string
	for _, v := range refs {
		ids = append(ids, v.(string))
	}
	works, _ := service.GetObjects("works", ids)
	if works != nil {
		for i, v := range works.Docs {
			if v.Found == true {
				var work = make(map[string]interface{})
				_ = json.Unmarshal(v.Source, &work)
				var host_venue = make(map[string]interface{})
				var newRef = map[string]interface{}{
					"id":               ids[i],
					"title":            work["title"],
					"publication_year": work["publication_year"],
				}
				if work["host_venue"] != nil {
					host_venue = work["host_venue"].(map[string]interface{})
					newRef["host_venue"] = host_venue["display_name"]
				} else {
					newRef["host_venue"] = ""
				}
				newReferencedWorks = append(newReferencedWorks, newRef)
			} else {
				println(ids[i] + " not found")
			}
		}
	}
	return newReferencedWorks
}
func GenAuthorDefaultIntro(a map[string]interface{}) string {
	var intro = fmt.Sprintf("I'm %s. ", a["display_name"].(string))
	if a["last_known_institution"] != nil {
		institution := a["last_known_institution"].(map[string]interface{})
		intro += fmt.Sprintf("I'm currently working at %s. ", institution["display_name"].(string))
	}
	if a["most_cited_work"] != nil {
		intro += fmt.Sprintf("I've post \"%s\", which is my most-cited work. ", a["most_cited_work"].(string))
	}
	if a["x_concepts"] != nil {
		x_concepts := a["x_concepts"].([]interface{})
		intro += fmt.Sprintf("I'm interested in ")
		for i, v := range x_concepts {
			if i > 3 {
				break
			}
			if i != 0 {
				intro += ", "
			}
			intro += v.(map[string]interface{})["display_name"].(string)
		}
		intro += ". "
	}
	intro += "I'm looking for highly motivate students..."
	return intro
}

// GetObject
// @Summary     ??????id???????????? txc
// @Description ??????id????????????????????????author???work???institution,venue,concept W4237558494,W2009180309,W2984203759
// @Tags    esSearch
// @Param       id     query    string true  "??????id"
// @Param       userid query    string false "??????id"
// @Success     200    {string} json   "{"status":200,"res":{}}"
// @Failure     404    {string} json   "{"status":201,"msg":"es get err or not found"}"
// @Failure     400    {string} json   "{"status":400,"msg":"id type error"}"
// @Router      /es/get/ [GET]
func GetObject(c *gin.Context) {
	id := c.Query("id")
	userid := c.Query("userid")
	id = utils.RemovePrefix(id)
	idx, err := utils.TransObjPrefix(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "id type error"})
		return
	}
	res, err := service.GetObject(idx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "es & openalex not found"})
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
	if idx == "works" {
		referenced_works := tmp["referenced_works"].([]interface{})
		tmp["referenced_works"] = TransRefs2Cited(referenced_works)
		related_works := tmp["related_works"].([]interface{})
		tmp["related_works"] = TransRefs2Intro(related_works)
		tmp["cited_string"] = map[string]interface{}{
			"mla": GenMLACited(tmp),
			"apa": GenAPACited(tmp),
			"gb":  GenGBCited(tmp),
		}
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
	if idx == "authors" {
		var info = make(map[string]interface{})
		if userid != "" {
			userid, _ := strconv.ParseUint(userid, 0, 64)
			user, notFound := service.GetUserByID(userid)
			if notFound {
				panic("user not found")
			}
			info["is_mine"] = user.AuthorID == id
			_, notFound = service.GetUserFollow(userid, id)
			info["isfollow"] = notFound == false
		}
		author, notFound := service.GetAuthor(id)
		if notFound {
			info["verified"] = false
			info["headshot"] = "author_default.jpg"
			info["intro"] = GenAuthorDefaultIntro(tmp)
		} else {
			info["verified"] = true
			info["headshot"] = author.HeadShot
			if author.Intro == "" {
				info["intro"] = GenAuthorDefaultIntro(tmp)
			} else {
				info["intro"] = author.Intro
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"data":   tmp,
			"info":   info,
			"status": 200,
		})
		return
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
// @Summary     ???????????? txc
// @Description ???????????????Cond????????????????????????key?????????["type", "author", "institution", "publisher", "venue", "publication_year"]
// @Tags        esSearch
// @Accept      json
// @Produce     json
// @Param       data body     response.BaseSearchQ true "????????????"
// @Success     200  {string} json                 "{"status":200,"res":{obeject}}"
// @Failure     201  {string} json                 "{"status":201,"err":"es search err"}"
// @Router      /es/search/base [POST]
func BaseSearch(c *gin.Context) {
	var d response.BaseSearchQ
	if err := c.ShouldBind(&d); err != nil {
		panic(err)
	}
	var aggs = make(map[string]bool)
	var fields = make([]string, 0)
	boolQuery := elastic.NewBoolQuery()
	if d.Kind == "string" {
		tiQuery := elastic.NewMatchPhraseQuery("full", d.QueryWord)
		b2Query := elastic.NewBoolQuery()
		b2Query.Should(tiQuery)
		boolQuery.Must(b2Query)
		fields = append(fields, "title", "abstract")
	} else if d.Kind == "title" {
		tiQuery := elastic.NewMatchPhraseQuery("title", d.QueryWord)
		b2Query := elastic.NewBoolQuery()
		b2Query.Should(tiQuery)
		boolQuery.Must(b2Query)
		fields = append(fields, "title")
	} else if d.Kind == "abstract" {
		tiQuery := elastic.NewMatchPhraseQuery("abstract", d.QueryWord)
		b2Query := elastic.NewBoolQuery()
		b2Query.Should(tiQuery)
		boolQuery.Must(b2Query)
		fields = append(fields, "abstract")
	} else if d.Kind == "venue" {
		tiQuery := elastic.NewMatchPhraseQuery("host_venue.display_name", d.QueryWord)
		b2Query := elastic.NewBoolQuery()
		b2Query.Should(tiQuery)
		boolQuery.Must(b2Query)
		fields = append(fields, "host_venue.display_name")
	} else if d.Kind == "institution" {
		tiQuery := elastic.NewMatchPhraseQuery("authorships.institutions.display_name", d.QueryWord)
		b2Query := elastic.NewBoolQuery()
		b2Query.Should(tiQuery)
		boolQuery.Must(b2Query)
		fields = append(fields, "authorships.institutions.display_name")
	}
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
	res, err := service.CommonWorkSearch(boolQuery, d.Page, d.Size, d.Sort, d.Asc, aggs, fields)
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
// @Summary     ???????????? txc
// @Description ???????????????query?????????map????????? ??????map??????"content" "field" "logic"
// @Description logic ?????????["and", "or", "not"]
// @Description field ?????????["title", "abstract", "venue", "publisher", "author", "institution", "concept"]
// @Description ???????????????????????????query?????? field???"publication_date" logic?????????and??? ???map??????"begin" "end"????????????????????????
// @Description sort=0??????????????????????????? =1????????????????????? =2???????????????????????????
// @Description asc=0????????? =1?????????
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
	fields := make([]string, 0)
	for _, i := range d.Query {
		fields = append(fields, i["field"])
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
	res, err := service.CommonWorkSearch(boolQuery, d.Page, d.Size, d.Sort, d.Asc, aggs, fields)
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
// @Summary     ??????doi??????work??????????????????????????? txc
// @Description ??????doi??????work???????????????????????????
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
// @Summary     ??????author???id???????????????????????? hr
// @Description ??????author???id????????????????????????, ???????????????Top N???????????????N=10???????????????????????????N???????????????????????????
// @Description
// @Description ????????????????????????1s, ??????????????????????????????????????????????????????????????????????????????
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

// GetStatistics
// @Summary     ?????????????????? txc
// @Description ??????????????????
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

// GetPrefixSuggestions doc
// @Summary     ???????????????????????????????????????results ??????????????? hr
// @description ???????????????????????????????????????results ???????????????
// @Tags        esSearch
// @Param       Field  query    string true "Field ??????????????????????????????"
// @Param       Prefix query    string true "Prefix ?????????????????????????????????"
// @Success     200    {string} string "{"success": true, "msg": "????????????"}"
// @Failure     400    {string} string "{"success": false, "msg": ????????????"}"
// @Failure     402    {string} string "{"success": false, "msg": "es????????????"}"
// @Router      /es/prefix [POST]
func GetPrefixSuggestions(c *gin.Context) {
	var d response.PrefixSuggestionQ
	if err := c.ShouldBind(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": "????????????"})
		panic(err)
	}
	index, field, prefix, topN := "works", d.Field, d.Prefix, 5
	log.Println("index:", index, "field:", field, "prefix:", prefix, "topN:", topN)
	prefixResult, err := service.PrefixSearch(index, field, prefix, topN)
	if err != nil {
		c.JSON(402, gin.H{"success": false, "msg": "es????????????", "err": err})
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "????????????", "res": prefixResult})
}

// AuthorSearch
// @Summary     txc
// @Description ??????????????????????????????,????????????
// @Tags        esSearch
// @Param       query_word query    string true "query_word"
// @Param       page       query    int    true "page"
// @Param       size       query    int    true "size"
// @Param       sort       query    int    true "sort"
// @Param       asc        query    bool   true "asc"
// @Success     200        {string} json   "{"res":{},"msg": "Author Search Success"}"
// @Failure     401        {string} json   "{"msg": "Author Search Error","err":err}"
// @Router      /es/search/author [GET]
func AuthorSearch(c *gin.Context) {
	queryWord := c.Query("query_word")
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))
	sort, _ := strconv.Atoi(c.Query("sort"))
	asc, _ := strconv.ParseBool(c.Query("asc"))
	res, err := service.AuthorSearch(queryWord, page, size, sort, asc)
	if err != nil {
		c.JSON(401, gin.H{
			"msg": "Author Search Error",
			"err": err,
		})
		return
	}
	var data = response.BaseSearchA{}
	data.Hits, data.Works, data.Aggs, _ = utils.NormalizationSearchResult(res)
	c.JSON(200, gin.H{
		"msg": "Author Search Success",
		"res": data,
	})
}

// AuthorSearch2
// @Summary     txc
// @Description ?????????????????????????????? via openalex
// @Tags        esSearch
// @Param       query_word query    string true  "query_word"
// @Param       page       query    int    false "page"
// @Param       size       query    int    false "size"
// @Param       sort       query    string false "sort=cited_by_count|...:_|desc|asc"
// @Success     200        {string} json   "{"res":{},"msg": "Author Search Success"}"
// @Failure     401        {string} json   "{"msg": "openalex Search Error","err":err}"
// @Failure     402        {string} json   "{"msg": "openalex Search Error","err":err}"
// @Router      /es/search/author2 [GET]
func AuthorSearch2(c *gin.Context) {
	//https://api.openalex.org/authors?search=kaiming%20he&page=2&per_page=10&sort=cited_by_count:desc
	//https://api.openalex.org/authors?search=kaiming%20he&group_by=last_known_institution.id
	search := c.Query("query_word")
	page := c.Query("page")
	per_page := c.Query("size")
	sort := c.Query("sort")
	url := "https://api.openalex.org/authors?"
	url += "search=" + search
	if page != "" {
		url += "&page=" + page
	}
	if per_page != "" {
		url += "&per_page=" + per_page
	}
	if sort != "" {
		url += "&sort=" + sort
	}
	data, err := utils.GetByUrl(url)
	if err != nil {
		c.JSON(401, gin.H{
			"msg": "openalex Search Error",
			"err": data,
		})
		return
	}
	if data["meta"] == nil {
		c.JSON(402, gin.H{
			"msg": "openalex Search Error",
			"err": data,
		})
		return
	}
	authors := data["results"].([]interface{})
	for _, v := range authors {
		author := v.(map[string]interface{})
		id := author["id"].(string)
		id = utils.RemovePrefix(id)
		author["id"] = id
		au, notFound := service.GetAuthor(id)
		if notFound {
			author["headshot"] = "author_default.jpg"
		} else {
			author["headshot"] = au.HeadShot
		}
	}
	c.JSON(200, gin.H{
		"msg": "Author Search Success",
		"res": map[string]interface{}{
			"hits":  data["meta"].(map[string]interface{})["count"],
			"works": data["results"],
		},
	})
}
