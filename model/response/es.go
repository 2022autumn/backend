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
	Vertex_set []map[string]interface{}
	Edge_set   []map[string]interface{}
}

type PrefixSuggestionQ struct {
	Field  string
	Prefix string
}

type PrefixSuggestionA struct {
	Suggestions []string
}
