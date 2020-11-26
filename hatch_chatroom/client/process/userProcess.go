package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	_ "go_code/chatroom/client/model"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"net"
)

type UserProcess struct {
}

func (this *UserProcess) Register(userID int, userPwd string, userName string) (err error) {
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("客户端 连接服务器失败...", err)
		return
	}
	defer conn.Close()

	// 2. 准备发送消息到服务端
	var mes message.Message            // 声明一个消息的结构体
	mes.Type = message.RegisterMesType // 赋值 消息类型

	var registerMes message.RegisterMes // 声明一个结构体，赋值为 LoginMes 类型
	registerMes.User.UserID = userID
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	data, err := json.Marshal(registerMes) // 因为 结构体不能直接赋值，先转为json，在转为 string
	if err != nil {
		fmt.Println("客户端 registerMes json.marshal err = ", err)
		return
	}

	// 3. data 赋值 mes 的 data， 将 mes 序列化发送
	mes.Data = string(data) // 上面赋值了 消息类型，现在赋值 消息内容的类型

	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("客户端 mes json.marshal err = ", err)
		return
	}

	// 6. 处理服务器返回的数据    创建transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}

	// 5. 发送消息本身
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送消息 err = ", err)
		return
	}

	mes2, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("客户端 readpkg err = ", err)
	}

	// 6. 将数据 反序列化为 正常结构体
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes2.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功")
	} else {
		fmt.Println(registerResMes.Error)
	}

	return
}

func (this *UserProcess) Login(userID int, userPwd string) (err error) {
	// 1. 连接到服务器端
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("客户端 连接服务器失败...", err)
		return
	}
	defer conn.Close()

	// 2. 准备发送消息到服务端
	var mes message.Message         // 声明一个消息的结构体
	mes.Type = message.LoginMesType // 赋值 消息类型

	var loginMes message.LoginMes // 声明一个结构体，赋值为 LoginMes 类型
	loginMes.UserId = userID
	loginMes.UserPwd = userPwd

	data, err := json.Marshal(loginMes) // 因为 结构体不能直接赋值，先转为json，在转为 string
	if err != nil {
		fmt.Println("客户端 loginMes json.marshal err = ", err)
		return
	}

	// 3. data 赋值 mes 的 data， 将 mes 序列化发送
	mes.Data = string(data) // 上面赋值了 消息类型，现在赋值 消息内容的类型
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("客户端 mes json.marshal err = ", err)
		return
	}

	// 4. 发送之前需要先计算 消息长度 发送过去
	var pkgLen uint32
	pkgLen = uint32(len(data)) // len(data)是字符串 -> byte -> unit -> []byte
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	_, err = conn.Write(buf[:4])
	if err != nil {
		fmt.Println("客户端 发送消息长度 err = ", err)
		return
	}

	// 5. 发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("客户端 发送消息本身err = ", err)
		return
	}
	fmt.Printf("客户端 发送消息成功 长度 = %d  内容 = %v \n", len(data), string(data))

	// 6. 处理服务器返回的数据    创建transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes2, err := tf.ReadPkg()
	if err != nil {
		fmt.Println("客户端 readpkg err = ", err)
	}
	// 6. 将数据 反序列化为 正常结构体
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes2.Data), &loginResMes)
	if loginResMes.Code == 200 {

		// 初始化 CurUser
		CurUser.Conn = conn
		CurUser.UserID = userID
		CurUser.UserStatus = message.UserOnline

		fmt.Println("当前在线用户列表如下：")
		for v, id := range loginResMes.UsersId {

			if v == userID {
				continue
			}
			fmt.Printf("用户 %d\n", id)
		}
		go serverProcessMes(conn)
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}

	return
}
