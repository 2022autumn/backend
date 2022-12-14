package v2

import (
	v1 "IShare/api/v1"
	"IShare/global"
	"IShare/model/database"
	"IShare/service"
	"IShare/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GetWorkCited(work map[string]interface{}) string {
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
	if work["title"] != nil {
		cited += "\"" + work["title"].(string) + "\""
	} else {
		cited += "\"\""
	}
	if work["host_venue"] != nil {
		if work["host_venue"].(map[string]interface{})["display_name"] != nil {
			cited += "," + work["host_venue"].(map[string]interface{})["display_name"].(string)
		}
	}
	cited += "," + strconv.Itoa(int(work["publication_year"].(float64))) + "."
	return cited
}
func TransRefs2Cited(refs []interface{}) []map[string]interface{} {
	var newRefs = make([]map[string]interface{}, 0)
	var ids []string
	for i, v := range refs {
		if i > 15 {
			break
		}
		ids = append(ids, v.(string))
	}
	works, _ := service.GetObjects2("works", ids)
	if works != nil && works["results"] != nil {
		works := works["results"].([]interface{})
		for _, v := range works {
			work := v.(map[string]interface{})
			newRefs = append(newRefs, map[string]interface{}{
				"id":    work["id"],
				"cited": GetWorkCited(v.(map[string]interface{})),
			})
		}
	}
	return newRefs
}
func TransRefs2Intro(refs []interface{}) []map[string]interface{} {
	var newRefs = make([]map[string]interface{}, 0)
	var ids []string
	for i, v := range refs {
		if i > 10 {
			break
		}
		ids = append(ids, v.(string))
	}
	works, _ := service.GetObjects2("works", ids)
	if works != nil && works["results"] != nil {
		works := works["results"].([]interface{})
		for _, v := range works {
			work := v.(map[string]interface{})
			newRef := map[string]interface{}{
				"id":               work["id"],
				"title":            work["title"],
				"publication_year": work["publication_year"],
			}
			if work["host_venue"] != nil && work["host_venue"].(map[string]interface{})["display_name"] != nil {
				host_venue := work["host_venue"].(map[string]interface{})
				newRef["host_venue"] = host_venue["display_name"]
			} else {
				newRef["host_venue"] = ""
			}
			newRefs = append(newRefs, newRef)
		}
	}
	return newRefs
}

// GetObject2
// @Summary     根据id获取对象 txc
// @Description 根据id获取对象，可以是author，work，institution,venue,concept W4237558494,W2009180309,W2984203759
// @Tags        esSearch
// @Param       id     query    string true  "对象id"
// @Param       userid query    string false "用户id"
// @Success     200    {string} json   "{"status":200,"res":{}}"
// @Failure     404    {string} json   "{"status":201,"msg":"es get err or not found"}"
// @Failure     400    {string} json   "{"status":400,"msg":"id type error"}"
// @Router      /es/get2/ [GET]
func GetObject2(c *gin.Context) {
	id := c.Query("id")
	userid := c.Query("userid")
	id = utils.RemovePrefix(id)
	idx, err := utils.TransObjPrefix(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "id type error"})
		return
	}
	res, err, source := service.GetObject2(idx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "es & openalex not found"})
		return
	}
	if userid != "" {
		userid, _ := strconv.ParseUint(userid, 0, 64)
		ucs, err := service.GetUserConcepts(userid)
		if err != nil {
			c.JSON(405, gin.H{"msg": "get user concepts err"})
			return
		}
		var concepts []interface{}
		if res["concepts"] != nil || res["x_concepts"] != nil {
			if res["concepts"] != nil {
				concepts = res["concepts"].([]interface{})
			}
			if res["x_concepts"] != nil {
				concepts = res["x_concepts"].([]interface{})
			}
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
		if idx == "works" {
			bh := database.BrowseHistory{
				UserID:          userid,
				WorkID:          id,
				Title:           res["title"].(string),
				PublicationYear: strconv.Itoa(int(res["publication_year"].(float64))),
				BrowseTime:      time.Now(),
			}
			if res["host_venue"] != nil {
				host_venue := res["host_venue"].(map[string]interface{})
				if host_venue["display_name"] != nil {
					bh.HostVenue = host_venue["display_name"].(string)
				}
			}
			err := global.DB.Create(&bh).Error
			if err != nil {
				log.Println(err, "create browse history err")
			}
		}
	}
	if idx == "works" {
		if source == 1 {
			if res["abstract_inverted_index"] != nil {
				res["abstract"] = utils.TransInvertedIndex2String(res["abstract_inverted_index"].(map[string]interface{}))
				res["abstract_inverted_index"] = nil
			}
		}
		referenced_works := res["referenced_works"].([]interface{})
		res["referenced_works"] = TransRefs2Cited(referenced_works)
		related_works := res["related_works"].([]interface{})
		res["related_works"] = TransRefs2Intro(related_works)
		res["cited_string"] = map[string]interface{}{
			"mla": v1.GenMLACited(res),
			"apa": v1.GenAPACited(res),
			"gb":  v1.GenGBCited(res),
		}
		authorworks, _ := service.GetWorksByWorkID(id) //添加pdf链接
		res["pdflinks"] = make([]string, 0)
		if res["open_access"] != nil {
			open_access := res["open_access"].(map[string]interface{})
			if open_access["oa_url"] != nil {
				res["pdflinks"] = append(res["pdflinks"].([]string), open_access["oa_url"].(string))
			} else if authorworks != nil && len(authorworks) != 0 {
				open_access["oa_url"] = "http://ishare.horik.cn:8000/api/media/pdf/" + authorworks[0].PDF
			}
		} else if authorworks != nil && len(authorworks) != 0 {
			res["open_access"] = make(map[string]interface{})
			res["open_access"].(map[string]interface{})["oa_url"] = "http://ishare.horik.cn:8000/api/media/pdf/" + authorworks[0].PDF
		}
		for _, v := range authorworks {
			if v.PDF != "" {
				res["pdflinks"] = append(res["pdflinks"].([]string), "http://ishare.horik.cn:8000/api/media/pdf/"+v.PDF)
			}
		}
		wv, notFound := service.GetWorkView(id)
		if notFound {
			wv = database.WorkView{
				WorkID:    id,
				Views:     1,
				WorkTitle: res["title"].(string),
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
		} else {
			info["is_mine"] = false
			info["isfollow"] = false
		}
		author, notFound := service.GetAuthor(id)
		if notFound {
			info["verified"] = false
			info["headshot"] = "author_default.jpg"
			info["intro"] = v1.GenAuthorDefaultIntro(res)
		} else {
			info["verified"] = true
			info["headshot"] = author.HeadShot
			if author.Intro == "" {
				info["intro"] = v1.GenAuthorDefaultIntro(res)
			} else {
				info["intro"] = author.Intro
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"data":   res,
			"info":   info,
			"status": 200,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":   res,
		"status": 200,
	})
}
