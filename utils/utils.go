package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
)

func BindJsonAndValid(c *gin.Context, model interface{}) interface{} {
	if err := c.ShouldBindJSON(&model); err != nil {
		//_, file, line, _ := runtime.Caller(1)
		//global.LOG.Panic(file + "(line " + strconv.Itoa(line) + "): bind model error")
		panic(err)
	}
	return model
}

func ShouldBindAndValid(c *gin.Context, model interface{}) error {
	if err := c.ShouldBind(&model); err != nil {
		return err
	}
	return nil
}

func GetMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func CloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		return
	}
}

func TransObjPrefix(id string) (ty string, err error) {
	if len(id) == 0 {
		return "", errors.New("id is empty")
	}
	switch id[0] {
	case 'W':
		return "works_v1", nil
	case 'A':
		return "authors", nil
	case 'I':
		return "institutions", nil
	case 'V':
		return "venues", nil
	case 'C':
		return "concepts", nil
	default:
		return "error type", errors.New("error type")
	}
}

// abstract_inverted_index: v
// 检错机制都放到了函数内，外部调用的时候不需要检错。感觉这样写也不是很优雅，传入一个interface{}的参数总感觉怪怪的
// TODO：CodeReview
func TransInvertedIndex2String(v interface{}) (abstract string) {
	if v == nil {
		return ""
	}
	abstract_map := make(map[int]string)
	keys := make([]int, 0)
	// 我们认为数据中abstract_inverted_index一定是一个<string 2 Slice>的map，并且Slice中的元素一定是Float64类型的数值
	if reflect.TypeOf(v).Kind() != reflect.Map {
		log.Println("abstract_inverted_index is not a map")
		log.Println(reflect.TypeOf(v).Kind())
		log.Println(v)
		return ""
	}
	for k1, v1 := range v.(map[string]interface{}) {
		if reflect.TypeOf(v1).Kind() != reflect.Slice {
			log.Println("abstract_inverted_index subelement is not a Slice")
			log.Println(reflect.TypeOf(v1).Kind())
			log.Println(v1)
			return ""
		}
		for _, v2 := range v1.([]interface{}) {
			if reflect.TypeOf(v2).Kind() != reflect.Float64 {
				log.Println("abstract_inverted_index subelement in Slice is not a int")
				log.Println(reflect.TypeOf(v2).Kind())
				log.Println(v2)
				return ""
			}
			keys = append(keys, int(v2.(float64)))
			abstract_map[int(v2.(float64))] = k1
		}
	}
	sort.Ints(keys)
	for _, v := range keys {
		abstract += abstract_map[v] + " "
	}
	return abstract
}

// 规范化es的返回结果
// hits 为es查询结果的总数
// result 为es查询结果的具体内容
// aggs 为es查询结果的聚合结果
// TookInMillis 为es查询耗时
func NormalizationSearchResult(res *elastic.SearchResult) (hits int64, result []json.RawMessage, aggs map[string]interface{}, TookInMillis int64) {
	if res == nil {
		return 0, nil, nil, 0
	}
	TookInMillis = res.TookInMillis
	hits = res.Hits.TotalHits.Value
	result = make([]json.RawMessage, 0)
	if res.Hits.Hits != nil {
		for _, hit := range res.Hits.Hits {
			result = append(result, hit.Source)
		}
	}
	aggs = make(map[string]interface{})
	if res.Aggregations != nil {
		for k, v := range res.Aggregations {
			by, _ := v.MarshalJSON()
			var tmp = make(map[string]interface{})
			_ = json.Unmarshal(by, &tmp)
			aggs[k] = tmp["buckets"].([]interface{})
		}
	}
	return hits, result, aggs, TookInMillis
}

// create works filter map
func InitWorksfilter() map[string]interface{} {
	worksfilter := make(map[string]interface{})

	worksfilter["id"] = true // id 需要修改 "https://openalex.org/W2741809807" -> "W2741809807"

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

	return worksfilter
}

// 保证filter中的key在data中存在
func FilterData(data *map[string]interface{}, filter *map[string]interface{}) {
	for k, v := range *filter {
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
				FilterData(&inner_data, &inner_filter)
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Slice {
			// 如果v为map的数组类型，则遍历data数组，递归
			if (*data)[k] != nil {
				inner_filter := v.([]map[string]interface{})[0]
				for _, value := range (*data)[k].([]interface{}) {
					inner_data := value.(map[string]interface{})
					FilterData(&inner_data, &inner_filter)
				}
			}
		}
	}
}

func RemovePrefix(s string) string {
	return strings.Replace(s, "https://openalex.org/", "", -1)
}
