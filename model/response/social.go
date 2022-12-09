package response

type CommentCreation struct {
	Content string `json:"content" binding:"required"`
	UserID  uint64 `json:"user_id" binding:"required"`
	PaperID string `json:"paper_id" binding:"required"`
}
type CommentListQuery struct {
	UserID  uint64 `json:"user_id" binding:"required"`
	PaperID string `json:"paper_id" binding:"required"`
}
type CommentUser struct {
	CommentID uint64 `json:"comment_id" binding:"required"`
	UserID    uint64 `json:"user_id" binding:"required"`
}
type FollowAuthorQ struct {
	UserID   uint64 `json:"user_id" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}
type GetUserFollowsQ struct {
	UserID uint64 `json:"user_id" binding:"required"`
}
