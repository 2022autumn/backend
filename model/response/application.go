package response

type CreateApplicationQ struct {
	AuthorName      string `json:"author_name" binding:"required"`
	InstitutionName string `json:"institution_name" binding:"required"`
	WorkEmail       string `json:"work_email" binding:"required"`
	//Field           string `json:"field" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
	UserID   uint64 `json:"user_id" binding:"required"`
}

type HandleApplicationQ struct {
	ApplicationID uint64 `json:"application_id" binding:"required"`
	UserID        uint64 `json:"user_id" binding:"required"`
	HandleRes     string `json:"success" binding:"required"` //是否通过
	HandleContent string `json:"content" binding:"required"` //审批意见
}
