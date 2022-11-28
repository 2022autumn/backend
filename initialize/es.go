package initialize

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/olivere/elastic/v7"

	"IShare/global"
)

func InitElasticSearch() {
	host := global.VP.GetString("es.host")
	ctx := context.Background()
	client, err := elastic.NewClient(
		elastic.SetURL(host),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(global.VP.GetString("es.username"), global.VP.GetString("es.password")),
		elastic.SetHealthcheckInterval(10*time.Second),
		//elastic.SetGzip(true),
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
	)
	if err != nil {
		panic(fmt.Errorf("get esclient error: %s", err))
	}
	info, code, err := client.Ping(host).Do(ctx)
	if err != nil {
		panic(fmt.Errorf("can't ping es"))
	}
	fmt.Printf("ping es code %d, version %s\n", code, info.Version.Number)
	global.ES = client
}
