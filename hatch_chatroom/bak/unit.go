package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"go_code/chatroom/common/message"
	"net"
)

func writePkg(conn net.Conn, data []byte) (err error) {
	// 1. 先发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data)) // len(data)是字符串 -> byte -> unit -> []byte
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)

	_, err = conn.Write(buf[:4])
	if err != nil {
		fmt.Println("客户端 发送消息长度 err = ", err)
		return
	}

	// 2. 发送数据给对方
	n, err := conn.Write(data)
	if err != nil {
		fmt.Println("客户端 发送消息长度 err = ", err)
		return
	}

	fmt.Printf("客户端 发送消息成功 长度 = %d  内容 = %v n = %v\n", pkgLen, string(data), n)
	return
}

func readPkg(conn net.Conn) (mes *message.Message, err error) {
	fmt.Println("服务端 等待客户端发送数据.....", conn.RemoteAddr())

	buf := make([]byte, 8096)
	_, err = conn.Read(buf[:])
	if err != nil {
		fmt.Println("服务端 conn.Read err = ", err)
		return
	}
	fmt.Println("服务端 读到的 buf = ", buf[0:4])

	// 根据读到的字节转成 unit32 进行判断
	var pkgLen uint64
	pkgLen = binary.BigEndian.Uint64(buf[0:4])
	fmt.Println("服务端 读到的 pkgLen = ", buf[0:pkgLen])

	// 根据 pkgLen 读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("服务端 conn.read buf err = ", err)
		return
	}

	// 把 pkgLen 反序列化成 消息类型
	json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("服务端 json marshal err = ", err)
		return
	}

	return
}
