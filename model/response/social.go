package response

type CommentCreation struct {
	Content string `json:"content" binding:"required"`
}
