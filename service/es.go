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
func GetObjects(index string, ids []string) (res *elastic.MgetResponse, err error) {
	mgetService := global.ES.MultiGet()
	for _, id := range ids {
		mgetService.Add(elastic.NewMultiGetItem().Index(index).Id(id))
	}
	return mgetService.Do(context.Background())
}
func GetObject(index string, id string) (res *elastic.GetResult, err error) {
	//termQuery := elastic.NewMatchQuery("id", id)
	//return global.ES.Search().Index(index).Query(termQuery).Do(context.Background())
	return global.ES.Get().Index(index).Id(id).Do(context.Background())
}
func CommonWorkSearch(boolQuery *elastic.BoolQuery, page int, size int,
	sortType int, ascending bool, aggs map[string]bool) (
	*elastic.SearchResult, error) {
	//typesAgg := elastic.NewTermsAggregation().Field("type.keyword")
	//institutionsAgg := elastic.NewTermsAggregation().Field("authorships.institutions.display_name.keyword")
	//publishersAgg := elastic.NewTermsAggregation().Field("host_venue.publisher.keyword")
	//venuesAgg := elastic.NewTermsAggregation().Field("host_venue.display_name.keyword")
	//authorsAgg := elastic.NewTermsAggregation().Field("authorships.author.display_name.keyword").Size(30)
	//minDateAgg, maxYearAgg := elastic.NewMinAggregation().Field("publication_year"), elastic.NewMaxAggregation().Field("publication_year")
	//publicationYearAgg := elastic.NewTermsAggregation().Field("publication_year")
	service := global.ES.Search().Index("works").Query(boolQuery).Size(size).TerminateAfter(LIMITCOUNT).Timeout("2s")
	addAggToSearch(service, aggs)
	//Aggregation("types", typesAgg).
	//Aggregation("institutions", institutionsAgg).
	//Aggregation("venues", venuesAgg).
	//Aggregation("publishers", publishersAgg).
	//Aggregation("authors", authorsAgg).
	//Aggregation("publication_years", publicationYearAgg)
	//Aggregation("min_year", minDateAgg).
	//Aggregation("max_year", maxYearAgg)
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
func addAggToSearch(service *elastic.SearchService, aggNames map[string]bool) *elastic.SearchService {
	if aggNames["types"] {
		service = service.Aggregation("types",
			elastic.NewTermsAggregation().Field("type.keyword"))
	}
	if aggNames["institutions"] {
		service = service.Aggregation("institutions",
			elastic.NewTermsAggregation().Field("authorships.institutions.display_name.keyword"))
	}
	if aggNames["venues"] {
		service = service.Aggregation("venues",
			elastic.NewTermsAggregation().Field("host_venue.display_name.keyword"))
	}
	if aggNames["publishers"] {
		service = service.Aggregation("publishers",
			elastic.NewTermsAggregation().Field("host_venue.publisher.keyword"))
	}
	if aggNames["authors"] {
		service = service.Aggregation("authors",
			elastic.NewTermsAggregation().Field("authorships.author.display_name.keyword").
				Size(30))
	}
	if aggNames["publication_years"] {
		service = service.Aggregation("publication_years",
			elastic.NewTermsAggregation().Field("publication_year"))
	}
	return service
}
