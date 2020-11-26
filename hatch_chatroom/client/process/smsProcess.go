package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/client/utils"
	"go_code/chatroom/common/message"
)

type SmsProcess struct {
}

func (this *SmsProcess) SendGoupMes(content string) (err error) {

	// 1. 创建一个mes
	var mes message.Message
	mes.Type = message.SmsMesType

	// 2. 创建一个 SmsMes 实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserID = CurUser.UserID
	smsMes.UserStatus = CurUser.UserStatus

	// 3. 序列化 smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json marhal err = ", err)
		return

	}

	mes.Data = string(data)

	// 4. 对 mes 再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("mes json marshal err = ", err)
		return
	}

	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("send data err = ", err)
		return
	}

	return
}
