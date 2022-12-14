package main

import (
	"IShare/global"
	"IShare/initialize"
	"IShare/utils"
	"bufio"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/olivere/elastic/v7"
	"golang.org/x/sync/semaphore"
)

var DEBUG = false

// 一般建议是1000-5000个文档，如果你的文档很大，可以适当减少队列，大小建议是5-15MB，默认不能超过100M，
var BULK_SIZE = 5000 // 一个文档约为1K，5000个文档约为5M

const (
	Limit  = 5 // 同时运行的goroutine上限
	Weight = 1 // 信号量的权重
)

var sem = semaphore.NewWeighted(Limit)

func main() {
	start_time := time.Now()
	initialize.InitViper()
	initialize.InitElasticSearch()
	filter := initFilter()
	if DEBUG {
		return
	}
	data_dir_path := []string{"/data/openalex/authors/", "/data/openalex/concepts/", "/data/openalex/institutions/", "/data/openalex/works/", "/data/openalex/venues/"}
	if DEBUG {
		data_dir_path = []string{"/data/openalex/testdata/authors/", "/data/openalex/testdata/concepts/", "/data/openalex/testdata/institutions/", "/data/openalex/testdata/works/", "/data/openalex/testdata/venues/"}
		// single test for works
		if true {
			data_dir_path = []string{"/data/openalex/testdata/works/"}
		}
	}
	var wg sync.WaitGroup
	for _, dir_path := range data_dir_path {
		current_dir_name := path.Base(dir_path)
		index := current_dir_name
		if index == "works" {
			index = "works_v1"
		}
		current_filter := filter[current_dir_name]
		files, err := ioutil.ReadDir(dir_path)
		if err != nil {
			log.Println("read dir: ", current_dir_name, " error: ", err)
			return
		}
		// 启动一个处理文件的协程
		for _, file := range files {
			wg.Add(1)
			if DEBUG {
				index = current_dir_name + "_test"
			}
			go processFile(dir_path, file.Name(), current_filter, &wg, index)
		}
	}
	wg.Wait()
	end_time := time.Now()
	log.Println("total time: ", end_time.Sub(start_time))
}

// create filter map
func initFilter() map[string]map[string]interface{} {
	filter := make(map[string]map[string]interface{})
	// 建立works json的过滤map
	filter["authors"] = initAuthorsfilter()
	filter["concepts"] = initConceptsfilter()
	filter["institutions"] = initInstitutionsfilter()
	filter["works"] = initWorksfilter()
	filter["venues"] = initVenuesfilter()
	return filter
}

// create authors filter map
func initAuthorsfilter() map[string]interface{} {
	authorsfilter := make(map[string]interface{})
	authorsfilter["id"] = true
	authorsfilter["orcid"] = false
	authorsfilter["display_name_alternatives"] = false
	authorsfilter["ids"] = make(map[string]interface{})
	authorsfilter["ids"].(map[string]interface{})["openalex"] = false
	authorsfilter["ids"].(map[string]interface{})["mag"] = false
	authorsfilter["last_known_institution"] = make(map[string]interface{})
	authorsfilter["last_known_institution"].(map[string]interface{})["id"] = true
	authorsfilter["last_known_institution"].(map[string]interface{})["ror"] = false
	authorsfilter["last_known_institution"].(map[string]interface{})["country_code"] = false
	authorsfilter["last_known_institution"].(map[string]interface{})["type"] = false
	authorsfilter["x_concepts"] = make([]map[string]interface{}, 0)
	x_concept := make(map[string]interface{})
	x_concept["id"] = true
	authorsfilter["x_concepts"] = append(authorsfilter["x_concepts"].([]map[string]interface{}), x_concept)
	authorsfilter["updated_date"] = false
	authorsfilter["created_date"] = false
	return authorsfilter
}

// create concepts filter map
func initConceptsfilter() map[string]interface{} {
	conceptsfilter := make(map[string]interface{})
	conceptsfilter["id"] = true
	conceptsfilter["ids"] = make(map[string]interface{})
	conceptsfilter["ids"].(map[string]interface{})["openalex"] = false
	conceptsfilter["ids"].(map[string]interface{})["mag"] = false
	conceptsfilter["international"] = false
	conceptsfilter["ancestors"] = make([]map[string]interface{}, 0)
	ancestor := make(map[string]interface{})
	ancestor["id"] = true
	conceptsfilter["ancestors"] = append(conceptsfilter["ancestors"].([]map[string]interface{}), ancestor)
	conceptsfilter["related_concepts"] = make([]map[string]interface{}, 0)
	related_concept := make(map[string]interface{})
	related_concept["id"] = true
	conceptsfilter["related_concepts"] = append(conceptsfilter["related_concepts"].([]map[string]interface{}), related_concept)
	conceptsfilter["updated_date"] = false
	conceptsfilter["created_date"] = false
	return conceptsfilter
}

// create institutions filter map
func initInstitutionsfilter() map[string]interface{} {
	institutionsfilter := make(map[string]interface{})
	institutionsfilter["id"] = true
	institutionsfilter["country_code"] = false
	institutionsfilter["ids"] = make(map[string]interface{})
	institutionsfilter["ids"].(map[string]interface{})["openalex"] = false
	institutionsfilter["ids"].(map[string]interface{})["ror"] = false
	institutionsfilter["ids"].(map[string]interface{})["mag"] = false
	institutionsfilter["geo"] = make(map[string]interface{})
	institutionsfilter["geo"].(map[string]interface{})["geonames_city_id"] = false
	institutionsfilter["geo"].(map[string]interface{})["country_code"] = false
	institutionsfilter["geo"].(map[string]interface{})["latitude"] = false
	institutionsfilter["geo"].(map[string]interface{})["longitude"] = false
	institutionsfilter["international"] = false
	institutionsfilter["associated_institutions"] = make([]map[string]interface{}, 0)
	associated_institution := make(map[string]interface{})
	associated_institution["id"] = true
	associated_institution["ror"] = false
	associated_institution["country_code"] = false
	associated_institution["type"] = false
	institutionsfilter["associated_institutions"] = append(institutionsfilter["associated_institutions"].([]map[string]interface{}), associated_institution)
	institutionsfilter["x_concepts"] = make([]map[string]interface{}, 0)
	x_concept := make(map[string]interface{})
	x_concept["id"] = true
	institutionsfilter["x_concepts"] = append(institutionsfilter["x_concepts"].([]map[string]interface{}), x_concept)
	institutionsfilter["updated_date"] = false
	institutionsfilter["created_date"] = false
	return institutionsfilter
}

// create works filter map
func initWorksfilter() map[string]interface{} {
	worksfilter := make(map[string]interface{})

	worksfilter["id"] = true // id 需要修改 "https://openalex.org/W2741809807" -> "W2741809807"

	worksfilter["display_name"] = false

	// worksfilter["publication_year"] = false

	worksfilter["ids"] = make(map[string]interface{})
	worksfilter["ids"].(map[string]interface{})["openalex"] = false
	worksfilter["ids"].(map[string]interface{})["mag"] = false
	worksfilter["ids"].(map[string]interface{})["doi"] = false

	worksfilter["host_venue"] = make(map[string]interface{}) // host_venue 需要修改
	worksfilter["host_venue"].(map[string]interface{})["id"] = true
	worksfilter["host_venue"].(map[string]interface{})["issn"] = false
	worksfilter["host_venue"].(map[string]interface{})["is_oa"] = false
	worksfilter["host_venue"].(map[string]interface{})["version"] = false
	worksfilter["host_venue"].(map[string]interface{})["license"] = false

	// 建立authorships数组
	worksfilter["authorships"] = make([]map[string]interface{}, 0) // authorships 需要修改
	// 建立authorships数组中的元素map
	authorship := make(map[string]interface{})
	authorship["author"] = make(map[string]interface{})
	authorship["author"].(map[string]interface{})["id"] = true // authorships.author.id 需要修改 "https://openalex.org/A1969205032" -> "A1969205032"
	authorship["author"].(map[string]interface{})["orcid"] = false
	authorship["raw_affiliation_string"] = false
	authorship["institutions"] = make([]map[string]interface{}, 0) // authorships.institutions 需要修改
	// 建立authorships.institutions数组中的元素map
	institution := make(map[string]interface{})
	institution["id"] = true // authorships.institutions.id 需要修改 "https://openalex.org/I1969205032" -> "I1969205032"
	institution["ror"] = false
	institution["country_code"] = false
	institution["type"] = false
	// 向authorships.institutions数组中添加元素map
	authorship["institutions"] = append(authorship["institutions"].([]map[string]interface{}), institution)
	// 向worksfilter.authorships中添加元素map
	worksfilter["authorships"] = append(worksfilter["authorships"].([]map[string]interface{}), authorship)

	worksfilter["biblio"] = false
	worksfilter["is_retracted"] = false
	worksfilter["is_paratext"] = false

	worksfilter["concepts"] = make([]map[string]interface{}, 0) // concepts 需要修改
	concept := make(map[string]interface{})
	concept["id"] = true // concepts.id 需要修改 "https://openalex.org/C1969205032" -> "C1969205032"
	worksfilter["concepts"] = append(worksfilter["concepts"].([]map[string]interface{}), concept)

	worksfilter["mesh"] = false
	worksfilter["alternate_host_venues"] = false

	worksfilter["abstract_inverted_index"] = "trans"

	worksfilter["referenced_works"] = true
	worksfilter["related_works"] = true

	worksfilter["ngrams_url"] = false
	worksfilter["updated_date"] = false
	worksfilter["created_date"] = false
	return worksfilter
}

// create venues filter map
func initVenuesfilter() map[string]interface{} {
	venuesfilter := make(map[string]interface{})
	venuesfilter["id"] = true
	venuesfilter["issn"] = false
	venuesfilter["is_in_doaj"] = false
	venuesfilter["ids"] = false
	venuesfilter["x_concepts"] = make([]map[string]interface{}, 0)
	x_concept := make(map[string]interface{})
	x_concept["id"] = true
	venuesfilter["x_concepts"] = append(venuesfilter["x_concepts"].([]map[string]interface{}), x_concept)
	venuesfilter["updated_date"] = false
	venuesfilter["created_date"] = false
	return venuesfilter
}

/**
* 处理一个json文件，经过过滤后，写入新的json文件，删除旧的json文件
* 新文件的名字是旧文件名字在最前面加上filterred_ eg: authors_data_10.json -> filtered_authors_data_10.json
* @param fileName: json文件绝对路径
* @param filter: 过滤map
 */
func processFile(dir_path string, fileName string, filter map[string]interface{}, wg *sync.WaitGroup, index string) {
	sem.Acquire(context.Background(), Weight)
	defer wg.Done()
	defer sem.Release(Weight)
	startTime := time.Now()

	log.Println("processFile: ", dir_path+fileName)
	file, err := os.Open(dir_path + fileName)
	if err != nil {
		log.Println("open file error: ", err)
		return
	}
	defer file.Close()
	// 文件大小为0，直接跳过
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("get file info error: ", err)
		return
	}
	if fileInfo.Size() == 0 {
		log.Println("file size is 0, skip: ", dir_path+fileName)
		return
	}
	reader := bufio.NewReader(file)
	client := global.ES
	bulkRequest := client.Bulk()
	for {
		// 1. 按行读取文件
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("read file error: ", err, " in file: ", fileName)
			return
		}
		// 2. json解析
		var data map[string]interface{}
		err = json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Println("json unmarshal error: ", err, " in file: ", fileName)
			return
		}
		filterData(&data, &filter)
		req := elastic.NewBulkIndexRequest().Index(index).Id(data["id"].(string)).Doc(data)
		bulkRequest = bulkRequest.Add(req)
		// 每BULK_SIZE条数据提交一次
		if bulkRequest.NumberOfActions() >= BULK_SIZE {
			// 插入时，遇到相同id的数据，会更新原数据
			_, err := bulkRequest.Do(context.Background())
			if err != nil {
				log.Println("bulk error: ", err, " error file: ", file)
				return
			}
			bulkRequest = client.Bulk()
		}
	}
	if bulkRequest.NumberOfActions() > 0 {
		_, err := bulkRequest.Do(context.Background())
		if err != nil {
			log.Println("bulk error: ", err, " error file: ", file)
			return
		}
	}
	os.Remove(dir_path + fileName)
	log.Println("processFile: ", dir_path+fileName, " done, cost time: ", time.Since(startTime))
}

// 保证filter中的key在data中存在
func filterData(data *map[string]interface{}, filter *map[string]interface{}) {
	for k, v := range *filter {
		if k == "abstract_inverted_index" {
			abstract := utils.TransInvertedIndex2String((*data)[k])
			// 删去abstract_inverted_index
			delete(*data, "abstract_inverted_index")
			// 添加abstract字段
			(*data)["abstract"] = abstract
		}
		// 如果v为bool类型，若为true则修改，若为false则删除
		if reflect.TypeOf(v).Kind() == reflect.Bool {
			if v.(bool) {
				// 修改规则类似："https://openalex.org/W2741809807" -> "W2741809807"
				if (*data)[k] != nil {
					// data[k] 为string类型，需要修改
					if reflect.TypeOf((*data)[k]).Kind() == reflect.String {
						(*data)[k] = strings.Replace((*data)[k].(string), "https://openalex.org/", "", -1)
					}
					// data[k] 为数组类型 每个元素都是string类型，都需要修改
					if reflect.TypeOf((*data)[k]).Kind() == reflect.Slice {
						for i, v := range (*data)[k].([]interface{}) {
							(*data)[k].([]interface{})[i] = strings.Replace(v.(string), "https://openalex.org/", "", -1)
						}
					}
				}
			} else {
				delete(*data, k)
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Map {
			// 如果v为map类型，则递归
			if (*data)[k] != nil {
				inner_data := (*data)[k].(map[string]interface{})
				inner_filter := v.(map[string]interface{})
				filterData(&inner_data, &inner_filter)
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Slice {
			// 如果v为map的数组类型，则遍历data数组，递归
			if (*data)[k] != nil {
				inner_filter := v.([]map[string]interface{})[0]
				for _, value := range (*data)[k].([]interface{}) {
					inner_data := value.(map[string]interface{})
					filterData(&inner_data, &inner_filter)
				}
			}
		}
	}
}
