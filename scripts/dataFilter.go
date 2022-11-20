package main

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

//--------------------- datafilter ---------------------
// datafilter文件主要是对openAlex下载下来的json文件进行过滤，只保留需要的字段，去掉不需要的字段，减少文件的大小
// 处理思路：读取json -> 过滤json对象 -> 写入新的json文件 -> 删除旧的json文件
// 起多个协程读取json文件，每个协程读取一个文件。每个协程按行读取文件，每行是一个json对象。然后将json对象转换成map[string]interface{}类型，然后过滤，最后追加写入新的json文件
//
//--------------------- datafilter ---------------------

//------- 1. 读取json文件 超大文件的读取 -------//
/**
* 需求: 读取大小在4G以上的文件，该文件中每一行都是一个json格式的字符串
* 问题：使用ioutil.ReadFile()读取文件时，一次性把所有文件内容读到内存中，导致内存占据过大，程序性能不佳，甚至崩溃（本机16GB还能抗住，但是多线程运行就别想了）鉴于总共需要处理的文件约900G，所以需要一种占据内存小的读取方式，方便后面开多线程处理（8核CPU闲着也是闲着）
* 解决：使用bufio.NewReader()读取文件，每次读取一行，然后再进行json解析
* 参考博客1：https://learnku.com/articles/23559/two-schemes-for-reading-golang-super-large-files
* 参考博客2: https://www.jianshu.com/p/509bb77ec103
* 参考博客3: https://zhuanlan.zhihu.com/p/184937550
 */

/**
* 处理一个json文件，经过过滤后，写入新的json文件，删除旧的json文件
* 新文件的名字是旧文件名字在最前面加上filterred_ eg: authors_data_10.json -> filtered_authors_data_10.json
* @param fileName: json文件绝对路径
* @param filter: 过滤map
 */
func processFile(dir_path string, fileName string, filter map[string]interface{}) {
	log.Println("processFile: ", dir_path+fileName)
	// 获取当前时间
	startTime := time.Now().UnixNano()
	file, err := os.Open(dir_path + fileName)
	if err != nil {
		log.Println("open file error: ", err)
		return
	}
	defer file.Close()

	outfile, outerr := os.OpenFile(dir_path+"filterred_"+fileName, os.O_WRONLY, 0644)
	if outerr != nil {
		os.Create(dir_path + "filterred_" + fileName)
		outfile, outerr = os.OpenFile(dir_path+"filterred_"+fileName, os.O_WRONLY, 0644)
		if outerr != nil {
			log.Println("open outfile error: ", outerr)
			return
		}
	}
	defer outfile.Close()

	reader := bufio.NewReader(file)

	for {
		// 1. 按行读取文件
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("read file error: ", err)
			return
		}
		// 2. json解析
		var data map[string]interface{}
		err = json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Println("json unmarshal error: ", err)
			return
		}
		// 3. 过滤
		if filter != nil {
			// 3.1 过滤
			filterData(&data, &filter)
			// 把过滤后的map转换成json字符串
			jsonStr, err := json.Marshal(data)
			if err != nil {
				log.Println("json marshal error: ", err)
				return
			}
			// 3.2 写入新的json文件
			// 按行把jsonstr写入文件
			_, err = outfile.WriteString(string(jsonStr) + "\n")
			if err != nil {
				log.Println("write file error: ", err)
				return
			}
		} else {
			// 3.1 打印error日志
			log.Printf("filter map is nil, fileName: %s", fileName)
			return
		}
	}
	// 4. 删除旧的json文件
	os.Remove(fileName)
	// 5. 计算时间
	endTime := time.Now().UnixNano()
	log.Printf("processFile: %s, time: %d ms", fileName, (endTime-startTime)/1000000)
}

// 保证filter中的key在data中存在
func filterData(data *map[string]interface{}, filter *map[string]interface{}) {
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
				filterData(&inner_data, &inner_filter)
			}
		} else if reflect.TypeOf(v).Kind() == reflect.Slice {
			// 如果v为map的数组类型，则遍历data数组，递归
			inner_filter := v.([]map[string]interface{})[0]
			for _, value := range (*data)[k].([]interface{}) {
				inner_data := value.(map[string]interface{})
				filterData(&inner_data, &inner_filter)
			}
		}
	}
}

//------- 2. 过滤json文件 filter初始化-------//
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
	conceptsfilter["image_url"] = false
	conceptsfilter["image_thumbnail_url"] = false
	conceptsfilter["international"] = false
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
	institutionsfilter["ids"].(map[string]interface{})["mag"] = false
	institutionsfilter["geo"] = false
	institutionsfilter["international"] = false
	institutionsfilter["associated_institutions"] = make([]map[string]interface{}, 0)
	associated_institution := make(map[string]interface{})
	associated_institution["id"] = true
	associated_institution["ror"] = false
	associated_institution["country_code"] = false
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

	worksfilter["publication_year"] = false

	worksfilter["ids"] = false

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

	authorship["institutions"] = make([]map[string]interface{}, 0) // authorships.institutions 需要修改
	// 建立authorships.institutions数组中的元素map
	institution := make(map[string]interface{})
	institution["id"] = true // authorships.institutions.id 需要修改 "https://openalex.org/I1969205032" -> "I1969205032"
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

	worksfilter["alternate_host_venues"] = false
	worksfilter["referenced_works"] = true
	worksfilter["related_works"] = true

	worksfilter["mesh"] = false
	worksfilter["updated_date"] = false
	worksfilter["created_date"] = false
	return worksfilter
}

// create venues filter map
func initVenuesfilter() map[string]interface{} {
	venuesfilter := make(map[string]interface{})
	venuesfilter["id"] = true
	venuesfilter["issn"] = false
	venuesfilter["is_oa"] = false
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

// 执行test之前需要先make filter获取数据
func test() {
	filter := initFilter()
	processFile("/home/horik/backend/scripts/", "work.json", filter["works"])
	processFile("/home/horik/backend/scripts/", "author.json", filter["authors"])
	processFile("/home/horik/backend/scripts/", "venue.json", filter["venues"])
	processFile("/home/horik/backend/scripts/", "institution.json", filter["institutions"])
	processFile("/home/horik/backend/scripts/", "concept.json", filter["concepts"])
}

func main() {

	// test()

	// 初始化过滤map
	filter := initFilter()

	// 过滤数据部分
	// data_dir_path := []string{"/data/openalex/authors/", "/data/openalex/concepts/", "/data/openalex/institutions/", "/data/openalex/works/", "/data/openalex/venues/"}
	data_dir_path := []string{"/data/openalex/venues/"}
	// 获取每个文件夹下的文件列表
	for _, dir_path := range data_dir_path {
		// 获取文件目录的最后一个目录名
		current_dir_name := path.Base(dir_path)
		current_filter := filter[current_dir_name]
		// fmt.Println(current_dir_name)
		files, err := ioutil.ReadDir(dir_path)
		if err != nil {
			log.Fatal(err)
		}
		// 启动一个处理文件的协程
		for _, file := range files {
			processFile(dir_path, file.Name(), current_filter)
		}
	}
}
