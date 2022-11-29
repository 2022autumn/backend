package response

type BaseSearchQ struct {
	Conds     map[string]string
	Page      int
	Size      int
	Kind      string
	QueryWord string
}
type SearchA struct {
}
