// package main

// import (
// 	"fmt"
// )

// // 定义用户名和密码
// var userID int
// var userPwd string

// func main() {
// 	// 加两个变量用于客户输入选项，和循环菜单控制
// 	var key int
// 	var loop = true

// 	for {
// 		fmt.Println("-----------------  欢迎登陆多人在线聊条系统 ----------------")
// 		fmt.Println("\t\t\t 1. 登录聊天室")
// 		fmt.Println("\t\t\t 2. 注册用户")
// 		fmt.Println("\t\t\t 3. 退出系统")
// 		fmt.Println("\t\t\t 请选择(1-3): ")

// 		fmt.Scanf("%d\n", &key)

// 		switch key {
// 		case 1:
// 			fmt.Println("登陆聊天室........")
// 			loop = false
// 		case 2:
// 			fmt.Println("注册用户........")
// 			loop = false
// 		case 3:
// 			fmt.Println("退出系统........")
// 			loop = false
// 		default:
// 			fmt.Println("请输入正确得选项(1-3): ")
// 		}

// 		if loop == false {
// 			break
// 		}
// 	}

// 	if key == 1 {
// 		fmt.Print("请输入用户id: ")
// 		fmt.Scanf("%d\n", &userID)
// 		fmt.Print("请输入密码: ")
// 		fmt.Scanf("%s\n", &userPwd)

// 		login(userID, userPwd)
// 	} else if key == 2 {
// 		fmt.Println("开始注册用户 ....")
// 	}
// }
