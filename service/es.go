package service

import (
	"IShare/global"
	"IShare/model/database"
	"IShare/utils"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/olivere/elastic/v7"
)

var LIMITCOUNT = 10000000

func GetWork(boolQuery *elastic.BoolQuery) (res *elastic.SearchResult, err error) {
	return global.ES.Search().Index("works").Query(boolQuery).Do(context.Background())
}

// 通过ids批量获取对象，multiGet
func GetObjects(index string, ids []string) (res *elastic.MgetResponse, err error) {
	mgetService := global.ES.MultiGet()
	for _, id := range ids {
		mgetService.Add(elastic.NewMultiGetItem().Index(index).Id(id))
	}
	return mgetService.Do(context.Background())
}

// 通过id获取对象，get
func GetObject(index string, id string) (res *elastic.GetResult, err error) {
	//termQuery := elastic.NewMatchQuery("id", id)
	//return global.ES.Search().Index(index).Query(termQuery).Do(context.Background())
	return global.ES.Get().Index(index).Id(id).Do(context.Background())
}
func GetObject2(index string, id string) (data map[string]interface{}, err error, source int) {
	res, err := global.ES.Get().Index(index).Id(id).Do(context.Background())
	if err != nil {
		log.Println("id", id)
		//https://api.openalex.org/works/W2741809807
		data, err := utils.GetByUrl("https://api.openalex.org/" + index + "/" + id)
		if err != nil {
			log.Println("<ERROR in GetObject2> GetByUrl error: ", err)
			return nil, err, 1
		}
		if index == "works" {
			workfilter := utils.InitWorksfilter()
			utils.FilterData(&data, &workfilter)
		} else if index == "authors" {
			authorfilter := utils.InitAuthorsfilter()
			utils.FilterData(&data, &authorfilter)
		}
		return data, err, 1
	}
	var data2 map[string]interface{}
	err = json.Unmarshal(res.Source, &data2)
	if err != nil {
		log.Println("GetObjects2 Unmarshal error: ", err)
		return nil, err, 0
	}
	return data2, err, 0
}

// mget对象 不保证顺序
func GetObjects2(index string, ids []string) (res map[string]interface{}, err error) {
	//https://api.openalex.org/works?filter=openalex_id:W1966354963|W1968537269|W1980250440|W2000332333|W2015137959|W2028830377
	/*
		W1966354963",
		      "W1968537269",
		      "W1980250440",
		      "W2000332333",
		      "W2015137959",
		      "W2028830377",
		      "W2033020845",
		      "W2040155070",
		      "W2045114325",
		      "W2050314443",
		      "W2057721183",
		      "W2068745974",
		      "W2069713827",
		      "W2073672062",
		      "W2073891001",
		      "W2081536057",
		      "W2126581974",
		      "W2131691431",
		      "W2148347826",
		      "W2796700885"
	*/
	if len(ids) == 0 {
		return map[string]interface{}{
			"results": []interface{}{},
		}, nil
	}
	url := "https://api.openalex.org/" + index + "?filter=openalex_id:"
	for i, id := range ids {
		if i != 0 {
			url += "|"
		}
		url += id
	}
	url += "&per_page=200"
	data, err := utils.GetByUrl(url)
	if err != nil {
		return nil, err
	}
	if index == "works" {
		workfilter := utils.InitWorksfilter()
		for _, v := range data["results"].([]interface{}) {
			work := v.(map[string]interface{})
			utils.FilterData(&work, &workfilter)
		}
	} else if index == "authors" {
		authorfilter := utils.InitAuthorsfilter()
		for _, v := range data["results"].([]interface{}) {
			author := v.(map[string]interface{})
			utils.FilterData(&author, &authorfilter)
		}
	}
	return data, nil
}

// 通用搜索，针对works
func CommonWorkSearch(boolQuery *elastic.BoolQuery, page int, size int,
	sortType int, ascending bool, aggs map[string]bool, fields []string) (res *elastic.SearchResult, err error) {
	timeout := global.VP.GetString("es.timeout")
	workIndex := "works"
	service := global.ES.Search().Index(workIndex).Query(boolQuery).Size(size).TerminateAfter(LIMITCOUNT).Timeout(timeout)
	addAggToSearch(service, aggs)
	// addHighlightToSearch(service, fields)
	if sortType == 0 {
		res, err = service.Sort("_score", ascending).From((page - 1) * size).Do(context.Background())
	} else if sortType == 1 {
		res, err = service.Sort("cited_by_count", ascending).From((page - 1) * size).Do(context.Background())
	} else if sortType == 2 {
		res, err = service.Sort("publication_date", ascending).From((page - 1) * size).Do(context.Background())
	}
	return res, err
}

// 为搜索添加高亮
func addHighlightToSearch(service *elastic.SearchService, fields []string) *elastic.SearchService {
	// 定义highlight
	highlight := elastic.NewHighlight()
	// 指定需要高亮的字段
	for _, field := range fields {
		highlight = highlight.Fields(elastic.NewHighlighterField(field))
	}
	// 指定高亮的返回逻辑 <div style='color: red;'>...msg...</div>
	highlight = highlight.PreTags("<div style='color: yellow;'>").PostTags("</div>")
	service = service.Highlight(highlight)
	return service
}

// 为搜索添加添加聚合
func addAggToSearch(service *elastic.SearchService, aggNames map[string]bool) *elastic.SearchService {
	if aggNames["types"] {
		service = service.Aggregation("types",
			elastic.NewTermsAggregation().Field("type.keyword"))
	}
	if aggNames["institutions"] {
		service = service.Aggregation("institutions",
			elastic.NewTermsAggregation().Field("authorships.institutions.display_name.keyword"))
	}
	if aggNames["venues"] {
		service = service.Aggregation("venues",
			elastic.NewTermsAggregation().Field("host_venue.display_name.keyword"))
	}
	if aggNames["publishers"] {
		service = service.Aggregation("publishers",
			elastic.NewTermsAggregation().Field("host_venue.publisher.keyword"))
	}
	if aggNames["authors"] {
		service = service.Aggregation("authors",
			elastic.NewTermsAggregation().Field("authorships.author.display_name.keyword").
				Size(30))
	}
	if aggNames["publication_years"] {
		service = service.Aggregation("publication_years",
			elastic.NewTermsAggregation().Field("publication_year").Order("_key", false))
	}
	return service
}

// 计算学者关系网络
func ComputeAuthorRelationNet(author_id string) (Vertex_set []map[string]interface{}, Edge_set []map[string]interface{}, err error) {
	author_id = utils.RemovePrefix(author_id)
	author_id = utils.RemovePrefix(author_id)
	log.Println("author_id: ", author_id)
	author_map, err, _ := GetObject2("authors", author_id)
	if err != nil {
		log.Println("GetObject err: ", err)
		return nil, nil, err
	}
	display_name := author_map["display_name"].(string)
	works_api_url := author_map["works_api_url"].(string)

	// 4. 获取author_id对应的author的所有作品
	works := make([]map[string]interface{}, 0)
	GetAllWorksByUrl(works_api_url, &works)

	Vertex_set = make([]map[string]interface{}, 0)
	Edge_set = make([]map[string]interface{}, 0)
	// label 只取前5个字符
	Vertex_set = append(Vertex_set, map[string]interface{}{
		"id":    author_id,
		"label": display_name[:5],
		"full":  display_name,
	})
	for _, work := range works {
		// work_id := work["id"].(string)
		// work_display_name := work["display_name"].(string)
		work_authorships := work["authorships"].([]interface{})
		for _, work_authorship := range work_authorships {
			work_authorship_map := work_authorship.(map[string]interface{})
			work_author_id := work_authorship_map["author"].(map[string]interface{})["id"].(string)
			exist := false
			// log.Println("work_author_id: ", work_author_id)
			for _, Vertex := range Vertex_set {
				// 判断是否已经存在, 如果存在则不添加 通过id string判断
				if Vertex["id"] == work_author_id {
					// log.Println("exist: ", Vertex["id"])
					exist = true
					break
				}
			}
			if !exist {
				// log.Println("not exist: ", work_author_id)
				work_author_display_name := work_authorship_map["author"].(map[string]interface{})["display_name"].(string)
				Vertex_set = append(Vertex_set, map[string]interface{}{
					"id":    work_author_id,
					"label": work_author_display_name[:5],
					"full":  work_author_display_name,
				})
			}
			if work_author_id != author_id {
				exist := false
				for _, Edge := range Edge_set {
					if Edge["from"] == author_id && Edge["to"] == work_author_id {
						exist = true
						Edge["weight"] = Edge["weight"].(int) + 1
						Edge["width"] = Edge["width"].(int) + 1
						dispaly_work := make(map[string]interface{})
						dispaly_work["id"] = work["id"].(string)
						dispaly_work["title"] = work["title"].(string)
						Edge["works"] = append(Edge["works"].([]interface{}), dispaly_work)
						break
					}
				}
				if !exist {
					Edge_set = append(Edge_set, map[string]interface{}{
						"from":   author_id,
						"to":     work_author_id,
						"weight": 1,
						"width":  1,
						"works":  []interface{}{},
					})
					dispaly_work := make(map[string]interface{})
					dispaly_work["id"] = work["id"].(string)
					dispaly_work["title"] = work["title"].(string)
					Edge_set[len(Edge_set)-1]["works"] = append(Edge_set[len(Edge_set)-1]["works"].([]interface{}), dispaly_work)
				}
			}
		}
	}
	// log.Println("Vertex_set: ", Vertex_set)
	// log.Println("Edge_set: ", Edge_set)
	TopVertex_set := make([]map[string]interface{}, 0)
	TopEdge_set := make([]map[string]interface{}, 0)
	GetTopN(&Vertex_set, &Edge_set, &TopVertex_set, &TopEdge_set, 10)
	return TopVertex_set, TopEdge_set, nil
}

// GetTopN 获取Top N 的Vertex和Edge
func GetTopN(Vertex_set *[]map[string]interface{}, Edge_set *[]map[string]interface{}, TopVertex_set *[]map[string]interface{}, TopEdge_set *[]map[string]interface{}, n int) {
	// 1. 获取Top N 的Edge
	sort.Slice(*Edge_set, func(i, j int) bool {
		return (*Edge_set)[i]["weight"].(int) > (*Edge_set)[j]["weight"].(int)
	})
	for i := 0; i < n && i < len(*Edge_set); i++ {
		*TopEdge_set = append(*TopEdge_set, (*Edge_set)[i])
	}
	for _, Edge := range *TopEdge_set {
		source := Edge["from"].(string)
		for _, Vertex := range *Vertex_set {
			if Vertex["id"] == source {
				*TopVertex_set = append(*TopVertex_set, Vertex)
				goto exit
			}
		}
	}
exit:
	// 2. 获取Top N Edge target 对应的Vertex
	for _, Edge := range *TopEdge_set {
		target := Edge["to"].(string)
		for _, Vertex := range *Vertex_set {
			if Vertex["id"] == target {
				*TopVertex_set = append(*TopVertex_set, Vertex)
				break
			}
		}
	}
}

// GetWorksByUrl 获取作者的作品，分页获取并返回总页数
func GetWorksByUrl(works_api_url string, page int, works *[]map[string]interface{}) (total_pages int, err error) {
	start_time := time.Now()
	request_url := works_api_url + "&page=" + strconv.Itoa(page)
	req_st_time := time.Now()
	resp, err := http.Get(request_url)
	log.Println("- single get works_api_url time: ", time.Since(req_st_time))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		log.Println("get works_api_url fail: ", string(body))
		return 0, errors.New("get works_api_url fail: " + string(body))
	}
	res := make(map[string]interface{})
	_ = json.Unmarshal(body, &res)
	count := int(res["meta"].(map[string]interface{})["count"].(float64))
	total_pages = int(math.Ceil(float64(count) / 25))
	works_list := res["results"].([]interface{})
	for _, work := range works_list {
		work_min := make(map[string]interface{})
		work_min["id"] = work.(map[string]interface{})["id"]
		work_min["title"] = work.(map[string]interface{})["title"]
		work_min["authorships"] = work.(map[string]interface{})["authorships"]
		*works = append(*works, work_min)
	}
	log.Println("-- single GetWorksByUrl time: ", time.Since(start_time))
	return
}

// GetAllPersonalWorksByUrl 获取作者的所有作品的列表
func GetAllPersonalWorksByUrl(works_api_url string, works *[]database.PersonalWorks, author_id string) (err error) {
	start_time := time.Now()
	data := make([]map[string]interface{}, 0)
	total_pages, err := GetWorksByUrl(works_api_url, 1, &data)
	if err != nil {
		log.Println("GetWorksByUrl err: ", err)
		return err
	}
	for i := 2; i <= total_pages; i++ {
		_, err := GetWorksByUrl(works_api_url, i, &data)
		if err != nil {
			log.Println("GetWorksByUrl err: ", err)
			return err
		}
	}
	for i, work := range data {
		presonal_work := database.PersonalWorks{
			AuthorID: author_id,
			WorkID:   utils.RemovePrefix(work["id"].(string)),
			Place:    i,
		}
		(*works) = append(*works, presonal_work)
	}
	log.Println("Total: GetAllWorksByUrl time: ", time.Since(start_time))
	return nil
}

// GetAllWorksByUrl 获取作者的所有作品的列表
func GetAllWorksByUrl(works_api_url string, works *[]map[string]interface{}) (err error) {
	start_time := time.Now()
	total_pages, err := GetWorksByUrl(works_api_url, 1, works)
	if err != nil {
		log.Println("GetWorksByUrl err: ", err)
		return err
	}
	for i := 2; i <= total_pages; i++ {
		_, err := GetWorksByUrl(works_api_url, i, works)
		if err != nil {
			log.Println("GetWorksByUrl err: ", err)
			return err
		}
	}
	filter := utils.InitWorksfilter()
	for _, work := range *works {
		utils.FilterData(&work, &filter)
	}
	log.Println("Total: GetAllWorksByUrl time: ", time.Since(start_time))
	return nil
}

// GetAuthorRelationNet 获取作者的关系网络，向上层提供接口
func GetAuthorRelationNet(authorid string) (Vertex_set []map[string]interface{}, Edge_set []map[string]interface{}, err error) {
	Vertex_set, Edge_set, err = ComputeAuthorRelationNet(authorid)
	return Vertex_set, Edge_set, err
}

var DocDict = []string{"authors", "works", "institutions", "venues", "concepts"}

// GetStatistics 获取学术数据的数目统计信息
func GetStatistics() (map[string]int64, error) {
	statistics := make(map[string]int64)
	for _, doc := range DocDict {
		count, err := global.ES.Count(doc).Do(context.Background())
		if err != nil {
			log.Println("GetStatistics of "+doc+" err: ", err)
			return nil, err
		}
		statistics[doc] = count
	}
	return statistics, nil
}

// elasticsearch前缀查询,给出前缀prefix，返回前topN个搜索提示结果
func PrefixSearch(index string, field string, prefix string, topN int) ([]string, error) {
	query := elastic.NewMatchPhrasePrefixQuery(field, prefix)
	// slop 参数告诉 match_phrase 查询词条相隔多远时仍然能将文档视为匹配 什么是相隔多远？ 意思是说为了让查询和文档匹配你需要移动词条多少次？
	query = query.MaxExpansions(topN * 2).Slop(2) //最多扩展到topN*2个结果
	searchResult, err := global.ES.Search().Index(index).Query(query).Size(topN).Do(context.Background())
	if err != nil {
		log.Println("PrefixSearch err: ", err)
		return nil, err
	}
	results := make([]string, 0)
	for _, hit := range searchResult.Hits.Hits {
		item := make(map[string]interface{})
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			panic(err)
		}
		results = append(results, item[field].(string))
	}
	log.Println(results)
	return results, nil
}

func AuthorSearch(queryWord string, page int, size int,
	sortType int, ascending bool) (res *elastic.SearchResult, err error) {
	//fuzzyQuery := elastic.NewFuzzyQuery("display_name", queryWord).Fuzziness(2)
	query := elastic.NewMatchPhraseQuery("display_name", queryWord).Slop(3)
	service := global.ES.Search().Index("authors").Query(query).Size(size).TerminateAfter(LIMITCOUNT)
	service = service.Aggregation("institution",
		elastic.NewTermsAggregation().Field("last_known_institution.display_name.keyword").Size(15))
	if sortType == 0 {
		res, err = service.Sort("_score", ascending).From((page - 1) * size).Do(context.Background())
	} else if sortType == 1 {
		res, err = service.Sort("cited_by_count", ascending).From((page - 1) * size).Do(context.Background())
	} else if sortType == 2 {
		res, err = service.Sort("works_count", ascending).From((page - 1) * size).Do(context.Background())
	}
	return res, err
}
