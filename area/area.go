/**
 *  场景管理
 *  创建多个channel,一个channel对应一个publisher,chanel从Hub订阅消息后分发给客户端
 *
 */

package area

import (
	"simple/global"
	"simple/lib/websocket"

	"simple/golog"
	//"simple/locks"
	"fmt"
	"net"
	"simple/lib/syncmap"
	z_type "simple/type"
	"time"

	//"strings"
)

// 预创建多个场景
func InitArea() {
	for _, area_id := range global.Config.Area.Init_area {
		CreateChannel(area_id, area_id)
		global.Channels[area_id] = global.Config.Hub.Hub_host
	}
}

// 创建一个RPC的场景
func CreateChannel(id string, name string) {
	golog.Info(id, name)
	global.RpcChannels = append(global.RpcChannels, id)
	global.SyncRpcChannelConns = append(global.SyncRpcChannelConns, syncmap.New())
	global.SyncRpcChannelWsConns = append(global.SyncRpcChannelWsConns, syncmap.New())
	//fmt.Println(global.RpcChannels)

}

// 删除一个RPC的场景
func RemovChannel(id string) {
	golog.Info(id)
	// @todo

}

// 检查是否已经创建了场景
func CheckChannelExist(name string) bool {
	//请使用快速查找法
	exist := 0
	for exist = range global.RpcChannels {
		if global.RpcChannels[exist] == name {
			break
		}
	}
	if exist == 0 {
		return false
	} else {
		return true
	}

}

/**
 *  socket连接 加入到场景中
 */
func SubscribeChannel(id string, conn *net.TCPConn, sid string) {

	index := 0
	for index = range global.RpcChannels {
		if global.RpcChannels[index] == id {
			break
		}
	}

	channel_conns := global.SyncRpcChannelConns[index]
	_, ok := channel_conns.Get(sid)
	if !ok {
		channel_conns.Set(sid, conn)
	}
	golog.Info("sid ", sid, "join ", id, global.SyncRpcChannelConns)

}

/**
 *  socket连接 加入到场景中
 */
func CheckUserJoinChannel(id string, sid string) bool {

	index := 0
	for index = range global.RpcChannels {
		if global.RpcChannels[index] == id {
			break
		}
	}

	channel_conns := global.SyncRpcChannelConns[index]
	_, ok := channel_conns.Get(sid)
	if ok {
		return true
	}

	channel_wss := global.SyncRpcChannelWsConns[index]
	_, ok = channel_wss.Get(sid)
	if ok {
		return true
	}
	return false
}

/**
 *  socket连接 加入到场景中
 */
func SubscribeWsChannel(id string, ws *websocket.Conn, sid string) {

	golog.Info(id, ws, sid)
	index := 0
	for index = range global.RpcChannels {
		if global.RpcChannels[index] == id {
			break
		}
	}
	channel_wss := global.SyncRpcChannelWsConns[index]
	_, ok := channel_wss.Get(sid)
	if !ok {
		channel_wss.Set(sid, ws)
	}

}

/**
 *  离开到场景
 */
func UnSubscribeChannel(id string, sid string) {

	golog.Info(id, sid)
	index := 0
	for index = range global.RpcChannels {
		if global.RpcChannels[index] == id {
			break
		}
	}

	channel_conns := global.SyncRpcChannelConns[index]
	_, ok := channel_conns.Get(sid)
	if !ok {
		channel_conns.Delete(sid)
	}

	channel_wss := global.SyncRpcChannelWsConns[index]
	_, ok = channel_wss.Get(sid)
	if !ok {
		channel_wss.Delete(sid)
	}

}

func UserUnSubscribeChannel(user_sid string) {

	for index, _ := range global.RpcChannels {
		UnSubscribeChannel(global.RpcChannels[index], user_sid)
	}
}

/**
 *  在场景中广播消息
 */
func Broatcast(id string, msg string) {

	golog.Info(id, msg)
	index := 0
	for index = range global.RpcChannels {
		if global.RpcChannels[index] == id {
			break
		}
	}
	channel_conns := global.SyncRpcChannelConns[index]
	var conn *net.TCPConn
	for item := range channel_conns.IterItems() {
		// fmt.Println("key:", item.Key, "value:", item.Value)
		conn = item.Value.(*net.TCPConn)
		conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
	}

	channel_wss := global.SyncRpcChannelWsConns[index]
	var wsconn *websocket.Conn
	for item := range channel_wss.IterItems() {
		// fmt.Println("key:", item.Key, "value:", item.Value)
		wsconn = item.Value.(*websocket.Conn)
		go websocket.Message.Send(wsconn, msg)
	}

}

func GetConn(sid string) *net.TCPConn {

	conn, ok := global.SyncUserConns.Get(sid)
	if !ok {
		return nil
	} else {
		return conn.(*net.TCPConn)
	}
}

func DeleteConn(sid string) {

	global.SyncUserConns.Delete(sid)

}

func GetWsConn(sid string) *websocket.Conn {
	wsconn, ok := global.SyncUserWebsocketConns.Get(sid)
	if !ok {
		return nil
	} else {
		return wsconn.(*websocket.Conn)
	}
}

func DeleteWsConn(sid string) {

	global.SyncUserWebsocketConns.Delete(sid)

}

func DeleteUserssion(sid string) {

	global.SyncUserSessions.Delete(sid)

}

func ConnRegister(conn *net.TCPConn, user_sid string) {

	SubscribeChannel("area-global", conn, user_sid)

	_, ok := global.SyncUserConns.Get(user_sid)
	if !ok {
		global.SyncUserConns.Set(user_sid, conn)
	}

	_, ok = global.SyncUserSessions.Get(user_sid)
	if !ok {
		data := &z_type.Session{
			conn.RemoteAddr().String(),
			"{}",
			true,  // 登录成功
			false, // 是否被踢出
			user_sid,
			time.Now().Unix(), //加入时间
			time.Now().Unix(),
		}
		global.SyncUserSessions.Set(user_sid, data)
	}

}

func WsConnRegister(ws *websocket.Conn, user_sid string) {

	golog.Debug("user_sid: ", user_sid)
	SubscribeWsChannel("area-global", ws, user_sid)

	_, ok := global.SyncUserWebsocketConns.Get(user_sid)
	if !ok {
		global.SyncUserWebsocketConns.Set(user_sid, ws)
	}

	_, ok = global.SyncUserSessions.Get(user_sid)
	if !ok {
		data := &z_type.Session{
			"", // @todo websocket ip
			"{}",
			true,  // 登录成功
			false, // 是否被踢出
			user_sid,
			time.Now().Unix(), //加入时间
			time.Now().Unix(),
		}
		global.SyncUserSessions.Set(user_sid, data)
	}

}
