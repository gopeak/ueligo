package hub

import (
	"morego/golog"
	"morego/lib/robfig/cron"
	json_orgin "encoding/json"
	//"fmt"
	"morego/area"
	"morego/global"
	"os"
	"path/filepath"
	"strings"
	//"net"
	"fmt"
	"github.com/gorilla/websocket"
	"morego/protocol"
	z_type "morego/type"
)

type Api struct {

	Init func()

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

func (api *Api)GetSession(sid string) string {
	session,exist := global.SyncUserSessions.Get(sid)
	if !exist {
		return "{}"
	}
	str,err := json_orgin.Marshal(session)
	if( err!=nil){
		golog.Error("Api GetSession json Marshal err:",err.Error())
		return "{}"
	}
	return string(str)
}


func (api *Api)Kick(sid string) bool {

	user_conn := area.GetConn(sid)
	if user_conn != nil {
		// 通知消息退出
		user_conn.Write( []byte(protocol.WrapRespErrStr("kicked")))
		area.FreeConn(user_conn,sid )

	}

	user_wsconn := area.GetWsConn(sid)
	if user_wsconn != nil {
		// 通知消息退出
		go user_wsconn.WriteMessage(websocket.TextMessage, []byte(protocol.WrapRespErrStr("kicked")) )
		area.FreeWsConn( user_wsconn,sid)
	}
	area.UserUnSubscribeChannel(sid)
	area.DeleteUserssion(sid)

	return true
}

func (api *Api)CreateChannel(id string, name string) bool {

	area.CreateChannel(id, name)
	return true
}

func (api *Api)RemoveChannel(id string) bool {

	area.RemovChannel(id)
	return true
}

func (api *Api)GetChannels() string {

	buf, err := json_orgin.Marshal(global.Channels)
	if( err!=nil ){
		return "{}"
	}
	return string(buf)

}

func (api *Api)GetSidsByChannel(channel_id string) string {

	buf,err:= json_orgin.Marshal(area.GetSidsByChannel(channel_id))
	if err!=nil {
		return string(buf)
	}else{
		return "[]"
	}

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

func (this *Api) Broadcast( sid string, area_id string, msg string) bool {

	area.Broatcast( sid, area_id ,msg)
	return true

}

func (this *Api) UpdateSession( sid string, data string ) bool {

	tmp, user_session_exist := global.SyncUserSessions.Get(sid)
	var user_session *z_type.Session
	if user_session_exist {
		user_session = tmp.(*z_type.Session)
		user_session.User = data
		global.SyncUserSessions.Set(sid, user_session)
	}
	return true

}



func (api *Api)BroadcastAll(msg string) bool {
	area.BroatcastGlobal("GM",msg)
	return true

}



func (api *Api)GetUserJoinedChannel(sid string) string {

	buf,err:=json_orgin.Marshal(area.GetSidsByChannel(sid))
	if( err!=nil ) {
		return  "[]"
	}
	return  string( buf )

}

func (api *Api)GetAllSession( ) string {

	var UserSessions = map[string]*z_type.Session{}
	for item := range global.SyncUserSessions.IterItems() {
		UserSessions[item.Key] = item.Value.(*z_type.Session)
	}
	ret, _ := json_orgin.Marshal(UserSessions)
	return  string(ret)

}

