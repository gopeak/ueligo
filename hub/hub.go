//
//  Hub server
//
//

package hub

import (
	"bufio"
	json2 "encoding/json"
	"fmt"
	"net"
	"morego/area"
	"morego/global"
	"morego/golog"
	"morego/lib/antonholmquist/jason"
	"morego/lib/websocket"
	z_type "morego/type"
	"strconv"
	"time"
)

/**
 * 监听客户端连接
 */
func HubServer() {

	hub_host := global.Config.Hub.Hub_host
	hub_port, _ := strconv.Atoi(global.Config.Hub.Hub_port)
	fmt.Println("Hub  Server :", hub_host, hub_port)
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(hub_host), hub_port, ""})
	if err != nil {
		golog.Error("Hub listenTCP Exception:", err.Error())
		return
	}

	hubListen(listen)
}

/**
 *  处理客户端连接
 */
func hubListen(listen *net.TCPListener) {

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			golog.Error("AcceptTCP Exception::", err.Error(), time.Now().UnixNano())
			break
		}
		// 校验ip地址
		conn.SetKeepAlive(true)
		///defer conn.Close()
		conn.SetNoDelay(false)

		//go handleWorkerWithJson( conn  )
		go handleHubConnWithBufferio(conn)

	} //end for {

}

func handleHubConnWithBufferio(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)

	for {

		msg, err := reader.ReadBytes('\n')

		if err != nil {
			//fmt.Println( "Hub handleWorker connection error: ", err.Error())
			// 超时处理
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {

			}
			closeHubConn(conn)
			break

		}
		if msg == nil {
			continue
		}
		//fmt.Println("handleHub  from :" , msg)

		go hubWorkeDispath(msg, conn)

	}

}

func closeHubConn(conn *net.TCPConn) {

	conn.Write([]byte{'E', 'O', 'F'})
	conn.Close()

}

//  Worker using REQ socket to do load-balancing
//
func hubWorkeDispath(msg []byte, conn *net.TCPConn) {

	//  Process messages as they arrive

	//fmt.Println( "hubWorkeDispath:", string(msg) )

	// 角色登录成功后
	ret_json, _ := jason.NewObjectFromBytes(msg)
	cmd, _ := ret_json.GetString("cmd")

	if cmd == "get_enable" {

		conn.Write([]byte(string(global.AppConfig.Enable)))
		conn.Close()
	}

	if cmd == "enabled" {
		global.AppConfig.Enable = 1
		conn.Write([]byte(string(global.AppConfig.Enable)))
	}

	if cmd == "disabled" {
		global.AppConfig.Enable = 0
		conn.Write([]byte(string(global.AppConfig.Enable)))
	}

	// 创建场景
	if cmd == "create_channel" {

		name, _ := ret_json.GetString("name")
		golog.Debug(" join_channel ", name)

		go area.CreateChannel(name, name)
		global.Channels[name] = ""

		conn.Write([]byte(`ok`))
	}

	// 销毁场景
	if cmd == "remove_channel" {

		id, _ := ret_json.GetString("name")
		golog.Debug("remove_channel:", id)
		area.RemovChannel(id)
		conn.Write([]byte(`ok`))

	}
	// 加入到一个场景中
	if cmd == "join_channel" {

		sid, _ := ret_json.GetString("sid")
		name, _ := ret_json.GetString("name")
		golog.Debug("join_channel:", sid, name)

		// 如果场景不存在，则返回错误
		exist := area.CheckChannelExist(name)
		if !exist {

			conn.Write([]byte(`error,channel not exist`))

		} else {

			// 检查会话用户是否加入过此场景
			have_joined := area.CheckUserJoinChannel(name, sid)

			// 如果还没有加入场景,则订阅
			if !have_joined {

				user_conn := area.GetConn(sid)
				channel_host := global.Channels[name]
				golog.Debug(" join_channel ", user_conn, channel_host, sid)
				user_wsconn := area.GetWsConn(sid)

				// 会话如果属于socket
				if user_conn != nil {
					go area.SubscribeChannel(name, user_conn, sid)
				}
				// 会话如果属于websocket
				if user_wsconn != nil {
					go area.SubscribeWsChannel(name, user_wsconn, sid)
				}
				var userJoinedChannels = make([]string, 0, 1000)
				tmp, ok := global.SyncUserJoinedChannels.Get(sid)
				if ok {
					userJoinedChannels = tmp.([]string)
				}
				userJoinedChannels = append(userJoinedChannels, name)
				global.SyncUserJoinedChannels.Set(sid, userJoinedChannels)
			}

			conn.Write([]byte(`ok`))
		}

	}

	if cmd == "leave_channel" {

		sid, _ := ret_json.GetString("sid")
		name, _ := ret_json.GetString("name")
		golog.Debug("remove_channel:", sid, name)
		// 离开场景则关闭此订阅
		go area.UnSubscribeChannel(name, sid)

		user_channels, exist := global.UserChannels[sid]
		if exist {
			for i := 0; i < len(user_channels); i++ {
				if user_channels[i] == name {
					user_channels = append(user_channels[:i], user_channels[i+1:]...)
					global.UserChannels[sid] = user_channels
					break
				}
			}
		}

		golog.Debug("userChannels's ", sid, ":", global.UserChannels[sid])
		golog.Debug("hub_worker leave_channel:", sid, name)
		conn.Write([]byte(`ok`))
	}

	if cmd == "kick" {
		sid, _ := ret_json.GetString("sid")
		user_conn := area.GetConn(sid)
		if user_conn != nil {
			// 发送消息退出
			user_conn.Write([]byte(`{"cmd":"error_","data":{"ret":0,"msg":"Server kicked " }}`))
			user_conn.Close()
			area.DeleteConn(sid)
		}

		user_wsconn := area.GetWsConn(sid)
		if user_wsconn != nil {
			// 发送消息退出
			websocket.Message.Send(user_wsconn, `{"cmd":"error_","data":{"ret":0,"msg":"Server kicked " }}`)
			area.DeleteWsConn(sid)
		}
		area.UserUnSubscribeChannel(sid)
		area.DeleteUserssion(sid)
		conn.Write([]byte(`ok`))
	}
	if cmd == "push" {

		sid, _ := ret_json.GetString("sid")
		push_data, _ := ret_json.GetString("data")

		user_conn := area.GetConn(sid)
		if user_conn != nil {
			user_conn.Write([]byte(fmt.Sprintf("%s\r\n", push_data)))
		}
		user_wsconn := area.GetWsConn(sid)
		if user_wsconn != nil {
			websocket.Message.Send(user_wsconn, fmt.Sprintf("%s\r\n", push_data))
		}

		golog.Debug("hub_worker push to  --------------->:", sid, push_data)
		conn.Write([]byte(`ok`))
	}

	if cmd == "broatcast" {
		sid, _ := ret_json.GetString("sid")
		channel, _ := ret_json.GetString("id")
		data, _ := ret_json.GetString("data")
		area.Broatcast( sid,channel, data)
		golog.Debug("hub_worker broadcast to :", channel, "   ", data)
		conn.Write([]byte(`ok`))
	}

	if cmd == "get_channels" {

		js1, _ := json2.Marshal(global.Channels)
		fmt.Println("(global.Channels:", (global.Channels))
		conn.Write(js1)
		conn.Close()
	}

	if cmd == "get_session" {

		sid, _ := ret_json.GetString("sid")
		user_session, exist := global.SyncUserSessions.Get(sid)
		js1 := []byte(`{}`)
		if exist {
			js1, _ = json2.Marshal(user_session)
		}
		conn.Write(js1)
		conn.Close()
	}

	if cmd == "get_all_session" {

		var UserSessions = map[string]*z_type.Session{}
		for item := range global.SyncUserSessions.IterItems() {
			UserSessions[item.Key] = item.Value.(*z_type.Session)
		}
		js1, _ := json2.Marshal(UserSessions)
		conn.Write(js1)
		conn.Close()

	}

	if cmd == "update_session" {

		sid, _ := ret_json.GetString("sid")
		data, _ := ret_json.GetString("data")
		tmp, user_session_exist := global.SyncUserSessions.Get(sid)
		var user_session *z_type.Session
		if user_session_exist {
			user_session = tmp.(*z_type.Session)
			user_session.User = data
			global.SyncUserSessions.Set(sid, user_session)
		}

		golog.Info("User Session  :", sid, user_session)
		conn.Write([]byte(`ok`))
	}

	if cmd == "get_user_join_channels" {

		sid, _ := ret_json.GetString("sid")
		js1 := []byte(`[]`)

		tmp, ok := global.SyncUserJoinedChannels.Get(sid)
		if ok {
			userJoinedChannels := tmp.([]string)
			js1, _ = json2.Marshal(userJoinedChannels)

		}
		// 发送消息退出
		conn.Write(js1)
		conn.Close()

	}

	if cmd == "set" {
		key, _ := ret_json.GetString("key")
		data, _ := ret_json.GetString("data")
		_, err := Set(key, data)
		if err == nil {
			conn.Write([]byte(`ok`))
		} else {
			conn.Write([]byte(`data server error!`))
		}
	}

	if cmd == "get" {
		key, _ := ret_json.GetString("key")
		reply, err := Get(key)
		if err == nil {
			conn.Write([]byte(reply))
		} else {
			conn.Write([]byte(`error`))
		}
	}
	if cmd == "delete" {
		key, _ := ret_json.GetString("key")
		_, err := Delete(key)
		if err == nil {
			conn.Write([]byte(`ok`))
		} else {
			conn.Write([]byte(`data server error!`))
		}
	}

}
