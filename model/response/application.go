package response

type CreateApplicationQ struct {
	RealName    string `json:"author_name" binding:"required"`
	Institution string `json:"institution" binding:"required"`
	Email       string `json:"email" binding:"required"`
	VerifyCode  string `json:"verify_code" binding:"required"`
	Content     string `json:"content"`
	AuthorID    string `json:"author_id" binding:"required"`
	UserID      uint64 `json:"user_id" binding:"required"`
}

type HandleApplicationQ struct {
	ApplicationID uint64 `json:"application_id" binding:"required"`
	UserID        uint64 `json:"user_id" binding:"required"` //谁审核的
	Status        int    `json:"status" binding:"required"`  //是否通过
	Content       string `json:"content"`                    //审批意见
}
type GetVerifyCodeQ struct {
	Email  string `json:"email" binding:"required"`
	UserID uint64 `json:"user_id" binding:"required"`
}
