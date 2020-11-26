package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"net"
)

//Transfer 将方法抽取出来，关联到结构体,分析有哪些字段
type Transfer struct {
	Conn net.Conn
	Buf  [8064]byte
}

func (this *Transfer) ReadPkg() (mes *message.Message, err error) {
	fmt.Println("服务端 等待客户端发送数据.....", this.Conn.RemoteAddr())

	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		fmt.Println("服务端 conn.Read err = ", err)
		return
	}
	fmt.Println("服务端 读到的 buf = ", this.Buf[0:4])

	// 根据读到的字节转成 unit32 进行判断
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])
	fmt.Println("服务端 读到的 pkgLen = ", this.Buf[:pkgLen])

	// 根据 pkgLen 读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("服务端 conn.read buf err = ", err)
		return
	}

	// 把 pkgLen 反序列化成 消息类型
	json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("服务端 json marshal err = ", err)
		return
	}

	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	// 1. 先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data)) // len(data)是字符串 -> byte -> unit -> []byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)

	_, err = this.Conn.Write(this.Buf[:4])
	if err != nil {
		fmt.Println("客户端 发送消息长度 err = ", err)
		return
	}

	// 2. 发送数据给对方
	n, err := this.Conn.Write(data)
	if err != nil {
		fmt.Println("客户端 发送消息长度 err = ", err)
		return
	}

	fmt.Printf("客户端 发送消息成功 长度 = %d  内容 = %v n = %v\n", pkgLen, string(data), n)
	return
}
