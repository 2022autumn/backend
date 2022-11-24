package utils

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"

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

func TransInvertedIndex2String(map[string]interface{}) string {

	return ""
}
