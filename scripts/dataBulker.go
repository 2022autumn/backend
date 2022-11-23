package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/olivere/elastic"
	"golang.org/x/sync/semaphore"
)

var logdata = true
var DEBUG = false
var BULK_SIZE = 10000
var BATCH_SIZE = 10
var DATA_PATHS = []string{"/data/openalex/authors/", "/data/openalex/institutions/", "/data/openalex/concepts/", "/data/openalex/venues/", "/data/openalex/works/"}
var TEST_DATA_PATHS = []string{"/data/testdata/authors/", "/data/testdata/institutions/", "/data/testdata/concepts/", "/data/testdata/venues/", "/data/testdata/works/"}
var FILE_PATHS_AUTHORS = []string{}
var FILE_PATHS_INSTITUTIONS = []string{}
var FILE_PATHS_CONCEPTS = []string{}
var FILE_PATHS_VENUES = []string{}
var FILE_PATHS_WORKS = []string{}
var client *elastic.Client

// 遍历DATA_PATHS，讲每个文件夹下的文件名读取出来，各自放到一个数组里面
func get_filepath() {
	if DEBUG {
		log.Print("debug mode")
		DATA_PATHS = TEST_DATA_PATHS
	}
	for _, path := range DATA_PATHS {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			switch path {
			case "/data/openalex/authors/", "/data/testdata/authors/":
				FILE_PATHS_AUTHORS = append(FILE_PATHS_AUTHORS, path+file.Name())
			case "/data/openalex/institutions/", "/data/testdata/institutions/":
				FILE_PATHS_INSTITUTIONS = append(FILE_PATHS_INSTITUTIONS, path+file.Name())
			case "/data/openalex/concepts/", "/data/testdata/concepts/":
				FILE_PATHS_CONCEPTS = append(FILE_PATHS_CONCEPTS, path+file.Name())
			case "/data/openalex/venues/", "/data/testdata/venues/":
				FILE_PATHS_VENUES = append(FILE_PATHS_VENUES, path+file.Name())
			case "/data/openalex/works/", "/data/testdata/works/":
				FILE_PATHS_WORKS = append(FILE_PATHS_WORKS, path+file.Name())
			}
		}
	}
}

const (
	Limit  = 50 // 同时运行的goroutine上限
	Weight = 1  // 信号量的权重
)

var sem = semaphore.NewWeighted(Limit)

func proc_files() {
	var wg sync.WaitGroup

	// bulk authors
	for _, file := range FILE_PATHS_AUTHORS {
		// 每个文件开启一个协程
		wg.Add(1)
		go bulk_file(file, &wg, "authors")
	}
	for _, file := range FILE_PATHS_INSTITUTIONS {
		wg.Add(1)
		go bulk_file(file, &wg, "institutions")
	}
	for _, file := range FILE_PATHS_CONCEPTS {
		wg.Add(1)
		go bulk_file(file, &wg, "concepts")
	}
	for _, file := range FILE_PATHS_VENUES {
		wg.Add(1)
		go bulk_file(file, &wg, "venues")
	}
	// for _, file := range FILE_PATHS_WORKS {
	// 	wg.Add(1)
	// 	go bulk_file(file, &wg, "works")
	// }
	wg.Wait()
}

func bulk_file(file string, wg *sync.WaitGroup, index string) {
	sem.Acquire(context.Background(), Weight)
	defer wg.Done()
	defer sem.Release(Weight)
	fd, err := os.Open(file)
	defer fd.Close()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("start bulk file: ", file)
	reader := bufio.NewReader(fd)
	bulkRequest := client.Bulk()
	for {
		// 1. 按行读取文件
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println("read file error: ", err, file)
			return
		}
		// 2. json解析
		var data map[string]interface{}
		err = json.Unmarshal([]byte(line), &data)
		if err != nil {
			log.Println("json unmarshal error: ", err)
			return
		}
		// 3. 构造bulk request
		req := elastic.NewBulkIndexRequest().Index(index).Id(data["id"].(string)).Doc(data)
		bulkRequest = bulkRequest.Add(req)
		// 4. 每BULK_SIZE条数据提交一次
		if bulkRequest.NumberOfActions() >= BULK_SIZE {
			_, err := bulkRequest.Do(context.Background())
			if err != nil {
				log.Println("bulk error: ", err, " error file: ", file)
				return
			}
			bulkRequest = client.Bulk()
		}
	}
	// 5. 提交剩余的数据
	if bulkRequest.NumberOfActions() > 0 {
		_, err := bulkRequest.Do(context.Background())
		if err != nil {
			log.Println("bulk error: ", err, " error file: ", file)
			return
		}
	}
	log.Println("bulk file success, remove it: ", file)
	os.Remove(file)
}

func newClient() {
	var err error
	client, err = elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		//elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	ctx := context.Background()
	info, code, err := client.Ping("http://localhost:9200").Do(ctx)
	if err != nil {
		panic(fmt.Errorf("can't ping es"))
	}
	log.Printf("ping es code %d, version %s\n", code, info.Version.Number)
	if err != nil {
		panic(err)
	}
}

func log_file() {
	for _, file := range FILE_PATHS_AUTHORS {
		f, err := os.OpenFile("/data/pred_authors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("open file error: ", err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(file + "\n"); err != nil {
			log.Println("write file error: ", err)
			return
		}
		// 将文件名写入到一个文件里面，作为日志
		f.Write([]byte(file + "\n"))
	}
	for _, file := range FILE_PATHS_INSTITUTIONS {
		f, err := os.OpenFile("/data/pred_institutions.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("open file error: ", err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(file + "\n"); err != nil {
			log.Println("write file error: ", err)
			return
		}
		// 将文件名写入到一个文件里面，作为日志
		f.Write([]byte(file + "\n"))
	}
	for _, file := range FILE_PATHS_WORKS {
		f, err := os.OpenFile("/data/pred_works.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("open file error: ", err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(file + "\n"); err != nil {
			log.Println("write file error: ", err)
			return
		}
		// 将文件名写入到一个文件里面，作为日志
		f.Write([]byte(file + "\n"))
	}
	for _, file := range FILE_PATHS_VENUES {
		f, err := os.OpenFile("/data/pred_venues.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("open file error: ", err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(file + "\n"); err != nil {
			log.Println("write file error: ", err)
			return
		}
		// 将文件名写入到一个文件里面，作为日志
		f.Write([]byte(file + "\n"))
	}
	for _, file := range FILE_PATHS_CONCEPTS {
		f, err := os.OpenFile("/data/pred_concepts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("open file error: ", err)
			return
		}
		defer f.Close()
		if _, err := f.WriteString(file + "\n"); err != nil {
			log.Println("write file error: ", err)
			return
		}
		// 将文件名写入到一个文件里面，作为日志
		f.Write([]byte(file + "\n"))
	}
}

func main() {
	newClient()
	ctx := context.Background()
	client.Ping("http://localhost:9200").Do(ctx)
	get_filepath()
	if logdata {
		log_file()
	}
	start_time := time.Now()
	proc_files()
	end_time := time.Now()
	log.Println("total time: ", end_time.Sub(start_time))
}
