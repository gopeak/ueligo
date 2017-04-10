package hub

import (
	"morego/golog"
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

type Api struct {

}

// 获取服务器的根路径
func (api *Api)GetBase() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		golog.Error("GetBase Error ", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)

}



func (api *Api)GetEnableStatus() bool {
	if global.AppConfig.Enable <= 0 {
		return false
	} else {
		return true
	}

}

func (api *Api)Enable() bool {

	global.AppConfig.Enable = 1
	return true

}

func (api *Api)Disable() bool {

	global.AppConfig.Enable = 0
	return true

}

func (api *Api)AddCron(expression string, exefnc func()) bool {

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

func (api *Api)RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}

	return true

}

func (api *Api)Get(key string) bool {

	return true

}

func (api *Api)Set(key string, value string) bool {

	return true

}

func (api *Api)GetSession(sid string) *z_type.Session  {
	session, _ := global.SyncUserSessions.Get(sid)
	return session.(*z_type.Session)
}

func (api *Api)GetSessionStr(sid string)  string {

	user_session, exist := global.SyncUserSessions.Get(sid)
	js1 := []byte(`{}`)
	if exist {
		js1, _ = json_orgin.Marshal(user_session)
	}
	return string(js1)

}

func (api *Api)Kick(sid string) bool {

	return true

}

func (api *Api)CreateChannel(id string, name string) bool {

	area.CreateChannel(id, name)
	return true
}

func (api *Api)RemoveChannel(id string) bool {

	return true
}

func (api *Api)GetChannels() bool {

	return true
}

func (api *Api)GetSidsByChannel(channel_id string) bool {

	return true

}

func (api *Api)ChannelAddSid(sid string, area_id string) bool {


	exist := area.CheckChannelExist(area_id)
	fmt.Println( area_id," CheckChannelExist:", exist )
	if !exist {
		return false
	}

	// 检查会话用户是否加入过此场景
	have_joined := area.CheckUserJoinChannel(area_id, sid)
	fmt.Println( "have_joined:", have_joined )
	// 如果还没有加入场景,则订阅
	if !have_joined {
		user_conn := area.GetConn(sid)
		user_wsconn := area.GetWsConn(sid)
		// 会话如果属于socket
		if user_conn != nil {
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

func (api *Api)ChannelKickSid(sid string, area_id string) bool {

	area.UnSubscribeChannel( area_id,sid)
	return true

}

func (api *Api)Push(from_sid string, to_sid string, msg string) bool {

	area.Push(to_sid, from_sid, msg)
	return true

}

func (api *Api)PushBySids(from_sid string,to_sids []string, msg string) bool {

	for _,to_sid:=   range to_sids {
		area.Push(to_sid, from_sid, msg)
	}
	return true

}


func (api *Api)BroadcastAll(msg string) bool {

	return true

}
