package service

import (
	"IShare/global"
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

	"github.com/olivere/elastic/v7"
)

var LIMITCOUNT = 10000000

func GetWork(boolQuery *elastic.BoolQuery) (res *elastic.SearchResult, err error) {
	return global.ES.Search().Index("works").Query(boolQuery).Do(context.Background())
}
func GetObject(index string, id string) (res *elastic.SearchResult, err error) {
	termQuery := elastic.NewMatchQuery("id", id)
	return global.ES.Search().Index(index).Query(termQuery).Do(context.Background())
}
func CommonWorkSearch(boolQuery *elastic.BoolQuery, page int, size int,
	sortType int, ascending bool, aggs map[string]bool) (
	*elastic.SearchResult, error) {
	//typesAgg := elastic.NewTermsAggregation().Field("type.keyword")
	//institutionsAgg := elastic.NewTermsAggregation().Field("authorships.institutions.display_name.keyword")
	//publishersAgg := elastic.NewTermsAggregation().Field("host_venue.publisher.keyword")
	//venuesAgg := elastic.NewTermsAggregation().Field("host_venue.display_name.keyword")
	//authorsAgg := elastic.NewTermsAggregation().Field("authorships.author.display_name.keyword").Size(30)
	//minDateAgg, maxYearAgg := elastic.NewMinAggregation().Field("publication_year"), elastic.NewMaxAggregation().Field("publication_year")
	//publicationYearAgg := elastic.NewTermsAggregation().Field("publication_year")
	service := global.ES.Search().Index("works").Query(boolQuery).Size(size).TerminateAfter(LIMITCOUNT).Timeout("2s")
	addAggToSearch(service, aggs)
	//Aggregation("types", typesAgg).
	//Aggregation("institutions", institutionsAgg).
	//Aggregation("venues", venuesAgg).
	//Aggregation("publishers", publishersAgg).
	//Aggregation("authors", authorsAgg).
	//Aggregation("publication_years", publicationYearAgg)
	//Aggregation("min_year", minDateAgg).
	//Aggregation("max_year", maxYearAgg)
	var res *elastic.SearchResult
	var err error
	if sortType == 0 {
		res, err = service.From((page - 1) * size).Do(context.Background())
	} else if sortType == 1 {
		res, err = service.Sort("cited_by_count", ascending).From((page - 1) * size).Do(context.Background())
	} else if sortType == 2 {
		res, err = service.Sort("publication_date", ascending).From((page - 1) * size).Do(context.Background())
	}
	return res, err
}
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
			elastic.NewTermsAggregation().Field("publication_year"))
	}
	return service
}

// 计算学者关系网络
func ComputeAuthorRelationNet(author_id string) (Vertex_set []map[string]interface{}, Edge_set []map[string]interface{}, err error) {
	// 1. 判断author_id类型
	ty, err := utils.TransObjPrefix(author_id)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}
	if ty != "authors" {
		log.Println("author_id is not an author id")
		return nil, nil, errors.New("author_id is not an author id")
	}
	// 2. 获取author_id对应的author
	res, err := GetObject(ty, author_id)
	if err != nil {
		log.Println("GetObject err: ", err)
		return nil, nil, err
	}
	// 2.1 判断author_id对应的author是否存在
	hits, authors, _, _ := utils.NormalizationSearchResult(res)
	if hits == 0 {
		log.Println("falut author id mapping to nil")
		return nil, nil, errors.New("falut author id mapping to nil")
	}
	author := authors[0]
	// 3. 反序列化author实体，获取实体中的display_name\ works_api_url
	var author_map map[string]interface{}
	_ = json.Unmarshal(author, &author_map)
	display_name := author_map["display_name"].(string)
	works_api_url := author_map["works_api_url"].(string)

	// 4. 获取author_id对应的author的所有作品
	works := make([]map[string]interface{}, 0)
	getAllWorksByUrl(works_api_url, &works)

	Vertex_set = make([]map[string]interface{}, 0)
	Edge_set = make([]map[string]interface{}, 0)
	Vertex_set = append(Vertex_set, map[string]interface{}{
		"id":    author_id,
		"label": display_name,
	})
	for _, work := range works {
		// work_id := work["id"].(string)
		// work_display_name := work["display_name"].(string)
		work_authorships := work["authorships"].([]interface{})
		for _, work_authorship := range work_authorships {
			work_authorship_map := work_authorship.(map[string]interface{})
			work_author_id := work_authorship_map["author"].(map[string]interface{})["id"].(string)
			exist := false
			for _, Vertex := range Vertex_set {
				if Vertex["id"] == work_author_id {
					exist = true
					break
				}
			}
			if !exist {
				work_author_display_name := work_authorship_map["author"].(map[string]interface{})["display_name"].(string)
				Vertex_set = append(Vertex_set, map[string]interface{}{
					"id":    work_author_id,
					"label": work_author_display_name,
				})
			}
			if work_author_id != author_id {
				exist := false
				for _, Edge := range Edge_set {
					if Edge["source"] == author_id && Edge["target"] == work_author_id {
						exist = true
						Edge["weight"] = Edge["weight"].(int) + 1
						break
					}
				}
				if !exist {
					Edge_set = append(Edge_set, map[string]interface{}{
						"source": author_id,
						"target": work_author_id,
						"weight": 1,
					})
				}
			}
		}
	}
	TopVertex_set := make([]map[string]interface{}, 0)
	TopEdge_set := make([]map[string]interface{}, 0)
	GetTopN(&Vertex_set, &Edge_set, &TopVertex_set, &TopEdge_set, 10)
	return TopVertex_set, TopEdge_set, nil
}

func GetTopN(Vertex_set *[]map[string]interface{}, Edge_set *[]map[string]interface{}, TopVertex_set *[]map[string]interface{}, TopEdge_set *[]map[string]interface{}, n int) {
	// 1. 获取Top N 的Edge
	sort.Slice(*Edge_set, func(i, j int) bool {
		return (*Edge_set)[i]["weight"].(int) > (*Edge_set)[j]["weight"].(int)
	})
	for i := 0; i < n && i < len(*Edge_set); i++ {
		*TopEdge_set = append(*TopEdge_set, (*Edge_set)[i])
	}
	for _, Edge := range *TopEdge_set {
		source := Edge["source"].(string)
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
		target := Edge["target"].(string)
		for _, Vertex := range *Vertex_set {
			if Vertex["id"] == target {
				*TopVertex_set = append(*TopVertex_set, Vertex)
				break
			}
		}
	}
}

func GetWorksByUrl(works_api_url string, page int, works *[]map[string]interface{}) (total_pages int, err error) {
	request_url := works_api_url + "&page=" + strconv.Itoa(page)
	resp, err := http.Get(request_url)
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
	return
}

func getAllWorksByUrl(works_api_url string, works *[]map[string]interface{}) (err error) {
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
	return nil
}

func GetAuthorRelationNet(authorid string) (Vertex_set []map[string]interface{}, Edge_set []map[string]interface{}, err error) {
	Vertex_set, Edge_set, err = ComputeAuthorRelationNet(authorid)
	return Vertex_set, Edge_set, err
}
