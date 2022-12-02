package response

import "encoding/json"

type BaseSearchQ struct {
	Conds     map[string]string
	Page      int
	Size      int
	Kind      string
	Sort      int
	Asc       bool
	QueryWord string
}
type BaseSearchA struct {
	Works []json.RawMessage
	Aggs  map[string]interface{}
	Hits  int64
}
type AdvancedSearchQ struct {
	Query []map[string]string
	Conds map[string]string
	Page  int
	Size  int
	Sort  int
	Asc   bool
}

type DoiSearchQ struct {
	Doi string
}
type GetObjectA struct {
	json.RawMessage
}

type AuthorRelationNet struct {
	Vertex_set []struct {
		id    string
		label string
	}
	Edge_set []struct {
		source string
		target string
		weight int
		works  []struct {
			id    string
			title string
		}
	}
}
