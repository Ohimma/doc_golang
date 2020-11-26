package process2

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"go_code/chatroom/server/model"
	"go_code/chatroom/server/utils"
	"net"
)

// 处理用户 登录 注册 注销 用户列表显示

type UserProcess struct {
	Conn net.Conn
	// 增加一个字段，表示该conn是哪个用户得
	UserID int
}

// 通知所有在线用户的方法，userid 会通知其他人说我上线了
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	// 遍历 onlineUser ，然后一个个发送xxx上线了
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		// 开始通知，单独写一个方法
		up.NotifyMeOnlineUser(userId)

	}
}

func (this *UserProcess) NotifyMeOnlineUser(userId int) {
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	// 将 notifyUserStatusMes 序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json marshal err = ", err)
		return
	}
	// 将序列化后的值 在赋给 mes.Data
	mes.Data = string(data)

	// 对mes进行序列化，然后发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json marshal err = ", err)
		return
	}

	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyOnline err = ", err)
		return
	}
}

func (this *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {

	// 1. 先从读到的结果取出 data，反序列化成 RegisterMes 结构体
	var RegisterMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &RegisterMes)
	if err != nil {
		fmt.Println("json unmarshal err = ", err)
		return
	}

	// 2. 声明一个 resMes 结构体
	var resMes message.Message
	resMes.Type = message.RegisterResMesType

	// 3. 声明一个返回的消息类型 loginResMes
	var RegisterResMes message.RegisterResMes

	// 取redis验证
	err = model.MyUserDao.Register(&RegisterMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			RegisterResMes.Code = 204
			RegisterResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			RegisterResMes.Code = 505
			RegisterResMes.Error = "注册发生未知错误"
		}
	} else {
		RegisterResMes.Code = 200
		fmt.Println("注册成功")
	}

	// 4. 将 loginResMes 序列化后输出
	data, err := json.Marshal(RegisterResMes)
	if err != nil {
		fmt.Println("json marshal err = ", err)
		return
	}

	// 5. 把 data 赋值给 resmes，然后对 resmes 序列化发送
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json marshal fail = ", err)
		return
	}
	// 6. 发送data

	// 1. 已经使用了 mvc 模式
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)

	return

}

func (this *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 1. 先从读到的结果取出 data，反序列化成 LoginMes 结构体
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json unmarshal err = ", err)
		return
	}

	// 2. 声明一个 resMes 结构体
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	// 3. 声明一个返回的消息类型 loginResMes
	var loginResMes message.LoginResMes

	// 取redis验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = 200
		// 将登陆成功的用户id，赋给this，this传给userMgr
		this.UserID = loginMes.UserId
		userMgr.AddOnlineUser(this)
		this.NotifyOthersOnlineUser(loginMes.UserId)
		// 将登陆成功的 userid，返回给新登录用户
		for id, _ := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}
		fmt.Println("登陆成功", user)
	}

	// if loginMes.UserID == 100 && loginMes.UserPwd == "123456" {
	// 	fmt.Println("用户合法")
	// 	loginResMes.Code = 200
	// } else {
	// 	fmt.Println("用户不合法")
	// 	loginResMes.Code = 500
	// 	loginResMes.Error = "用户不合法，请先注册"
	// }

	// 4. 将 loginResMes 序列化后输出
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json marshal err = ", err)
		return
	}

	// 5. 把 data 赋值给 resmes，然后对 resmes 序列化发送
	resMes.Data = string(data)
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json marshal fail = ", err)
		return
	}
	// 6. 发送data

	// 1. 已经使用了 mvc 模式
	tf := &utils.Transfer{
		Conn: this.Conn,
	}

	err = tf.WritePkg(data)

	return
}
