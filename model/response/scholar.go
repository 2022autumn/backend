package response

type AddUserConceptQ struct {
	UserID    uint64 `json:"user_id" binding:"required"`
	ConceptID string `json:"concept_id" binding:"required"`
}
