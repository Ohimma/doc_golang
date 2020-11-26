package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

// 总的消息类型和实体
type Message struct {
	Type string `json: "type"` // 消息类型
	Data string `json: "type"` // 消息内容的类型
}

//用户登录和返回
type LoginMes struct {
	UserId   int    `json: "userid"`   // 用户id
	UserPwd  string `json: "userpwd"`  // 用户密码
	UserName string `json: "username"` // 用户名
}

type LoginResMes struct {
	Code    int    `json: "code"`    // 返回状态码
	UsersId []int  `json: "usersid"` // 返回给用户在线用户列表
	Error   string `json: "error"`   // 返回的状态消息
}

// 注册用户和返回
type RegisterMes struct {
	User User `json: "user"`
}

type RegisterResMes struct {
	Code  int    `json: "code"`  // 返回状态码
	Error string `json: "error"` // 返回的状态消息
}

// 为了配合服务器推送用户状态变化消息
type NotifyUserStatusMes struct {
	UserId int `json: "userid"`
	Status int `json: "status"`
}

// 增加一个发送的消息
type SmsMes struct {
	Content string `json: "content"`
	User
}
