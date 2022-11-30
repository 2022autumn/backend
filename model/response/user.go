package response

type LoginQ struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterQ struct {
	Username string `json:"username" binding:"min=3,max=100,required"`
	Password string `json:"password" binding:"gte=6,required"`
}
