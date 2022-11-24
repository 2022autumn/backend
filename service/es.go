package service

import (
	"IShare/global"
	"context"

	"github.com/olivere/elastic/v7"
)

var LIMITCOUNT = 10000000

func GetWork(boolQuery *elastic.BoolQuery) (res *elastic.SearchResult, err error) {
	return global.ES.Search().Index("works").Query(boolQuery).Do(context.Background())
}
func GetObject(index string, id string) (res *elastic.SearchResult, err error) {
	termQuery := elastic.NewMatchQuery("id", id)
	return global.ES.Search().Index(index).Query(termQuery).Do(context.Background())
}
func CommonWorkSearch(page int, size int,
	boolQuery *elastic.BoolQuery, sortType int, ascending bool) (
	*elastic.SearchResult, error) {
	typesAgg := elastic.NewTermsAggregation().Field("type.keyword")
	institutionsAgg := elastic.NewTermsAggregation().Field("authorships.institutions.keyword")
	publishersAgg := elastic.NewTermsAggregation().Field("host_venue.publisher.keyword")
	authorsAgg := elastic.NewTermsAggregation().Field("authorships.author.display_name.keyword")
	// minDateAgg, maxYearAgg := elastic.NewMinAggregation().Field("publication_date"), elastic.NewMaxAggregation().Field("publication_date")
	service := global.ES.Search().Query(boolQuery).Size(size).TerminateAfter(LIMITCOUNT).
		Aggregation("type", typesAgg).
		Aggregation("institution", institutionsAgg).
		Aggregation("publishers", publishersAgg).
		Aggregation("authors", authorsAgg)
	var res *elastic.SearchResult
	var err error
	if sortType == 0 {
		res, err = service.From((page - 1) * size).Do(context.Background())
	} else if sortType == 1 {
		res, err = service.Sort("cited_by_count", ascending).From((page - 1) * size).Do(context.Background())
	} else if sortType == 2 {
		res, err = service.Sort("publication_date", ascending).From((page - 1) * size).Do(context.Background())
	}
	return res, err
}
