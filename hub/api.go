package hub

import (
	"morego/golog"
	"github.com/robfig/cron"
	json_orgin "encoding/json"
	"morego/area"
	"morego/global"
	"os"
	"strings"
	_"morego/lib/websocket"
	"morego/protocol"
	z_type "morego/type"
)

type Api struct {

	Init func()

}


// 获取服务器的根路径
func (api *Api)GetBase() string {

	dir, err:= os.Getwd()
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
		go user_wsconn.Write( []byte(protocol.WrapRespErrStr("kicked")) )
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

	return  area.ChannelAddSid( sid , area_id )

}

func (api *Api)ChannelKickSid(sid string, area_id string) bool {

	area.UnSubscribeChannel( area_id,sid)
	return true

}

func (api *Api)Push( from_sid string ,to_sid string , data  string  ) bool {

	area.Push( to_sid, from_sid, data  )
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

