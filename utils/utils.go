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
	switch id[0] {
	case 'W':
		return "works", nil
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

// 计算学者关系网络
func ComputeAuthorRelationNet(authot_id string) {
	ty, err := TransObjPrefix(authot_id)
	if err != nil {
		log.Println(err)
		return
	}
	if ty != "authors" {
		log.Println("authot_id is not an author id")
		return
	}
	// res, err := service.GetObject(ty, authot_id)
	// if err != nil {
	// 	log.Println("GetObject err: ", err)
	// 	return
	// }
	// // 1. 获取学者的所有作品

}
