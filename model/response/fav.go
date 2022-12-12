package response

type TagCreation struct {
	UserID  uint64 `json:"user_id" binding:"required"`
	TagName string `json:"tag_name" binding:"required"`
}
type AddTagToPaper struct {
	UserID  uint64 `json:"user_id" binding:"required"`
	PaperID string `json:"paper_id" binding:"required"`
	TagID   uint64 `json:"tag_id" binding:"required"`
}

type TagPaperListQ struct {
	UserID uint64 `json:"user_id" binding:"required"`
	TagID  uint64 `json:"tag_id" binding:"required"`
}

type UserInfo struct {
	UserID uint64 `json:"user_id" binding:"required"`
}

type RenameTagQ struct {
	UserID     uint64 `json:"user_id" binding:"required"`
	TagID      uint64 `json:"tag_id" binding:"required"`
	NewTagName string `json:"new_tag_name" binding:"required"`
}

type PaperBelongingQ struct {
	UserID  uint64 `json:"user_id" binding:"required"`
	PaperID string `json:"paper_id" binding:"required"`
}
