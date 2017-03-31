package worker

import (
	"morego/golog"
	"morego/lib/antonholmquist/jason"
	"morego/lib/robfig/cron"
	json_orgin "encoding/json"
	//"fmt"
	"morego/area"
	"morego/global"
	z_type "morego/type"
	"os"
	"path/filepath"
	"strings"
	//"net"
	"fmt"
)

// 获取服务器的根路径
func GetBase() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		golog.Error("GetBase Error ", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)

}

func GetConfig() *jason.Object {

	return global.ConfigJson

}

func GetEnableStatus() bool {
	if global.AppConfig.Enable <= 0 {
		return false
	} else {
		return true
	}

}

func Enable() bool {

	global.AppConfig.Enable = 1
	return true

}

func Disable() bool {

	global.AppConfig.Enable = 0
	return true

}

func AddCron(expression string, exefnc func()) bool {

	if cron, ok := global.Crons[expression]; ok {
		golog.Info("cron exist :", cron)
		return false
	}
	c := cron.New()
	c.AddFunc(expression, exefnc)
	c.Start()
	global.Crons[expression] = c
	return true

}

func RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}

	return true

}

func Get(key string) bool {

	return true

}

func Set(key string, value string) bool {

	return true

}

func GetSession(sid string) *z_type.Session  {
	session, _ := global.SyncUserSessions.Get(sid)
	return session.(*z_type.Session)
}

func GetSessionStr(sid string)  string {

	user_session, exist := global.SyncUserSessions.Get(sid)
	js1 := []byte(`{}`)
	if exist {
		js1, _ = json_orgin.Marshal(user_session)
	}
	return string(js1)

}

func Kick(sid string) bool {

	return true

}

func CreateChannel(id string, name string) bool {

	area.CreateChannel(id, name)
	return true
}

func RemoveChannel(id string) bool {

	return true
}

func GetChannels() bool {

	return true
}

func GetSidsByChannel(channel_id string) bool {

	return true

}

func ChannelAddSid(sid string, area_id string) bool {


	exist := area.CheckChannelExist(area_id)
	if !exist {

		return false
	}

	// 检查会话用户是否加入过此场景
	have_joined := area.CheckUserJoinChannel(area_id, sid)
	fmt.Println( "have_joined:", have_joined )
	// 如果还没有加入场景,则订阅
	if !have_joined {
		user_conn := area.GetConn(sid)
		channel_host := global.Channels[area_id]
		golog.Debug(" join_channel ", user_conn, channel_host, sid)
		user_wsconn := area.GetWsConn(sid)
		fmt.Println( "user_conn:", user_conn )
		// 会话如果属于socket
		if user_conn != nil {
			fmt.Println( "area_id:", area_id )
			area.SubscribeChannel(area_id, user_conn, sid)
		}
		// 会话如果属于websocket
		if user_wsconn != nil {
			area.SubscribeWsChannel(area_id, user_wsconn, sid)
		}
		// 该用户加入过的场景列表
		var userJoinedChannels = make([]string, 0, 1000)
		tmp, ok := global.SyncUserJoinedChannels.Get(sid)
		if ok {
			userJoinedChannels = tmp.([]string)
		}
		userJoinedChannels = append(userJoinedChannels, area_id)
		global.SyncUserJoinedChannels.Set(sid, userJoinedChannels)
	}



	return true

}

func ChannelKickSid(sid string, id string) bool {

	return true

}

func Push(sid string, msg string) bool {

	return true

}

func PushBySids(sids []string, msg string) bool {

	return true

}

func PushAll(msg string) bool {

	return true

}

func Broadcast(msg string) bool {

	return true

}
