package response

import "encoding/json"

type BaseSearchQ struct {
	Conds     map[string]string
	Page      int
	Size      int
	Kind      string
	Sort      int
	QueryWord string
}
type BaseSearchA struct {
	Works []json.RawMessage
	Aggs  map[string][]interface{}
	Hits  int64
}
type SearchA struct {
}
