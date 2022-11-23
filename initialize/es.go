package initialize

import (
	"IShare/global"
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"github.com/olivere/elastic/v7"
)

func InitElasticSearch() {
	host := global.VP.GetString("es.host")
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		//elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(fmt.Errorf("get esclient error"))
	}
	info, code, err := client.Ping(host).Do(ctx)
	if err != nil {
		panic(fmt.Errorf("can't ping es"))
	}
	fmt.Printf("ping es code %d, version %s\n", code, info.Version.Number)
	global.ES = client
}
