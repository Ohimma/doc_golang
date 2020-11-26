package model

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go_code/chatroom/common/message"
)

var (
	MyUserDao *UserDao
)

// 定义 UserDao结构体
// 完成对 User 得各种操作

type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 1. 根据用户 id 返回 用户信息和 err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	res, err := redis.String(conn.Do("hget", "users", id))
	if err != nil {
		if err == redis.ErrNil { // 表示users哈希中，没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &User{}

	// json 反序列化为 User 实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json unmarshal err = ", err)
		return
	}

	return
}

// 2. 完成对登录得验证，
func (this *UserDao) Login(userID int, userPwd string) (user *User, err error) {
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userID)
	if err != nil {
		return
	}

	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

func (this *UserDao) Register(user *message.User) (err error) {
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserID)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	// 这时说明 id 在 redis 还没有，可以完成注册
	data, err := json.Marshal(user)
	if err != nil {
		return
	}

	// 注册入库
	_, err = conn.Do("HSET", "users", user.UserID, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err = ", err)
		return
	}
	return
}
