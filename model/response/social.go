package response

type CommentCreation struct {
	Content string `json:"content" binding:"required"`
	UserID  uint64 `json:"user_id" binding:"required"`
	PaperID string `json:"paper_id" binding:"required"`
}
