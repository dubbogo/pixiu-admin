package entity

type UserInfo struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Role string `json:"role"`
	//Enable bool `json:"enable"`
	//Createdate time.Time
	//Updatedate time.Time
}
