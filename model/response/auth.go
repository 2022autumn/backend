package response

type AuthResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
}

type CommonA struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
}

type LoginQ struct {
	Username string `json:"username" binding:"min=3,max=100,required"`
	Password string `json:"password" binding:"gte=6,required"`
}

type RegisterQ struct {
	Name     string `json:"name" binding:"min=3,max=100,required"`
	Password string `json:"password" binding:"gte=6,required"`
}

//
