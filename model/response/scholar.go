package response

type AddUserConceptQ struct {
	ConceptID string `json:"concept_id" binding:"required"`
}
type ModifyAuthorIntroQ struct {
	AuthorID string `json:"author_id" binding:"required"`
	Intro    string `json:"intro"`
}

type GetPersonalWorksQ struct {
	AuthorID string `json:"author_id" binding:"required"`
	Page     int    `json:"page" binding:"required"`
	PageSize int    `json:"page_size" binding:"required"`
	Display  int    `json:"display"`
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

type GetPaperPDFQ struct {
	// AuthorID string `json:"author_id" binding:"required"`
	WorkID string `json:"work_id" binding:"required"`
}
