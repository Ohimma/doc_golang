package model

import (
	"go_code/chatroom/common/message"
	"net"
)

// 因为在客户端，很多机会使用，做成全局的

type CurUser struct {
	Conn net.Conn
	message.User
}
