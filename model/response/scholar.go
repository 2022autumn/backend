package response

type AddUserConceptQ struct {
	UserID    uint64 `json:"user_id" binding:"required"`
	ConceptID string `json:"concept_id" binding:"required"`
}

type GetPersonalWorksQ struct {
	AuthorID string `json:"author_id" binding:"required"`
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
}

type IgnoreWorkQ struct {
	AuthorID string `json:"author_id" binding:"required"`
	WorkID   string `json:"work_id" binding:"required"`
}

type ModifyPlaceQ struct {
	AuthorID  string `json:"author_id" binding:"required"`
	WorkID    string `json:"work_id" binding:"required"`
	Direction int    `json:"direction" binding:"required"`
}

type TopWorkQ struct {
	AuthorID string `json:"author_id" binding:"required"`
	WorkID   string `json:"work_id" binding:"required"`
}
