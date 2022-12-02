package response

type LoginQ struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterQ struct {
	Username string `json:"username" binding:"min=3,max=100,required"`
	Password string `json:"password" binding:"gte=6,required"`
}

type ModifyQ struct {
	UserInfo    string `gorm:"size:64;" json:"user_info"`     //个性签名
	Name        string `gorm:"size:32;" json:"name"`          //真实姓名
	Phone       string `gorm:"size:32;" json:"phone"`         //电话号码
	Email       string `gorm:"size:32;" json:"email"`         //邮箱
	Fields      string `gorm:"size:256;" json:"fields"`       //研究领域
	InterestTag string `gorm:"size:256;" json:"interest_tag"` //兴趣词
}

//binding:"required就是gin自带的数据验证，表示数据不为空，为空则返回错误;
