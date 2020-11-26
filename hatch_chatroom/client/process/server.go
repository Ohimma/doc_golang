package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
	"net"
	"os"
)

func ShowMenu() {
	fmt.Println("------------- 恭喜xxxx登陆成功 -----------")
	fmt.Println("------------- 1. 显示在线用户列表 -----------")
	fmt.Println("------------- 2. 发送群聊 -----------")
	fmt.Println("------------- 3. 信息列表 -----------")
	fmt.Println("------------- 4. 退出系统 -----------")
	var key int
	fmt.Print("请选择(1-4):")
	fmt.Scanf("%d\n", &key)

	var content string
	smsProcess := &SmsProcess{}
	switch key {
	case 1:
		// fmt.Println("显示在线用户列表........")
		outputOnlineUser()
	case 2:
		// fmt.Println("发送消息........")
		fmt.Println("请输入你想对大家所的话")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGoupMes(content)
	case 3:
		fmt.Println("信息列表........")
	case 4:
		fmt.Println("退出系统........")
		os.Exit(0)
	default:
		fmt.Println("请输入正确得选项(1-3): ")
	}
}

func serverProcessMes(conn net.Conn) {
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端等待服务端发送消息")

		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.Readpkg err = ", err)
			return
		}

		switch mes.Type {
		case message.NotifyUserStatusMesType:
			// 1. 取出 NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			// 2. 把用户的信息，状态，保存到客户的map[int]User 中
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType:
			outputGroupMes(mes)
		default:
			fmt.Println("服务器返回了未知的消息类型")
		}
		// fmt.Printf("mes = %v\n", mes)
	}
}
