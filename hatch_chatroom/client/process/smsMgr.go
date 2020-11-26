package process

import (
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
)

func outputGroupMes(mes *message.Message) {
	// 1. 反序列化 mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json unmarsha err = ", err)
		return
	}

	// 显示信息
	info := fmt.Sprintf("用户id=%d   对大家说%s\n", smsMes.UserID, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
