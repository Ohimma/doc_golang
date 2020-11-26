package process2

import (
	"fmt"
)

// UserMgr 实例 在服务器只有一个，但是用得到地方有很多，所以定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成对userMgr 初始化
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// 完成对 onlineMgr 得添加
func (this *UserMgr) AddOnlineUser(up *UserProcess) {
	this.onlineUsers[up.UserID] = up
}

// 完成对 onlineMgr 的删除
func (this *UserMgr) DelOnlineUser(userID int) {
	delete(this.onlineUsers, userID)
}

// 完成对 onlineMgr 的查询
func (this *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return this.onlineUsers
}

// 根据 id 返回对应的值 点对点
func (this *UserMgr) GetOnlineUserById(userID int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[userID]
	if !ok {
		err = fmt.Errorf("用户 %d 不存在", userID)
		return
	}
	return
}
