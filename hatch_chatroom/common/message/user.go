package message

type User struct {
	UserID     int    `json: "userid"`
	UserPwd    string `json: "userpwd"`
	UserName   string `json: "username"`
	UserStatus int    `json: "userstatus"`    // 用户状态
}
