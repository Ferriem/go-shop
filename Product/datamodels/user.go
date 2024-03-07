package datamodels

type User struct {
	ID           int64  `json:"id" form:"ID" sql:"id"`
	NickName     string `json:"nickName" form:"nickName" sql:"nickName"`
	UserName     string `json:"userName" form:"userName", sql:"userName"`
	HashPassword string `json:"-" form:"password" sql:"password"`
}

type UserInfo struct {
	ID       int64  `json:"id"`
	UserName string `json:"userName`
	IP       string `json:"ip"`
}
