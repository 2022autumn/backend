package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/gin-gonic/gin"
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
	abstract_map := make(map[int]string)
	keys := make([]int, 0)
	// 我们认为数据中abstract_inverted_index一定是一个<string 2 Slice>的map，并且Slice中的元素一定是Float64类型的数值
	if reflect.TypeOf(v).Kind() != reflect.Map {
		fmt.Println("abstract_inverted_index is not a map")
		fmt.Println(reflect.TypeOf(v).Kind())
		fmt.Println(v)
		return ""
	}
	for k1, v1 := range v.(map[string]interface{}) {
		if reflect.TypeOf(v1).Kind() != reflect.Slice {
			fmt.Println("abstract_inverted_index subelement is not a Slice")
			fmt.Println(reflect.TypeOf(v1).Kind())
			fmt.Println(v1)
			return ""
		}
		for _, v2 := range v1.([]interface{}) {
			if reflect.TypeOf(v2).Kind() != reflect.Float64 {
				fmt.Println("abstract_inverted_index subelement in Slice is not a int")
				fmt.Println(reflect.TypeOf(v2).Kind())
				fmt.Println(v2)
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
