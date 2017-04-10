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
	"time"
	"sync"
	"container/list"
	"net"
	"strconv"
)


type Sdk struct {


	GetBase func() (string)

	GetConfig func() (*jason.Object)

	GetEnable func() bool

	Enable func() bool

	Disable func() bool

	AddCron func(expression string, exefnc func()) bool

	RemoveCron func(expression string) bool

	Set func(key string, args ...interface{}) (bool, error)

	Get func(key string) (string, error)

	GetSession func(sid string) *z_type.Session

	GetSessionStr func(sid string) string

	Kick func(sid string) bool

	CreateChannel func(id string, name string) bool

	GetChannels func() map[string]string

	RemoveChannel func(id string) bool

	ChannelAddSid func(sid string, area_id string) bool

	ChannelKickSid func(sid string, area_id string) bool

	Push func(from_sid string, to_sid string, msg string) bool

	PushBySids func(from_sid string,to_sids []string, msg string) bool

	BroadcastAll func(msg string) bool

	Broadcast func( area_id string,msg string) bool

	Connected bool

	HubConn *net.TCPConn



}

func (s *Sdk) connect() bool{

	data :=  global.Config.WorkerServer.ToHub
	hub_host := data[0]
	hub_port_str := data[1]
	ip_port := worker_host + ":" + worker_port_str


	hubconn, err_req := net.DialTimeout("tcp", ip_port, 5 * time.Second)
	if( err_req!=nil ){
		return false
	}
	s.HubConn = hubconn

}

// 获取服务器的根路径
func (s *Sdk)  GetBase() string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		golog.Error("GetBase Error ", err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)

}

func (s *Sdk) GetConfig() *jason.Object {

	return global.ConfigJson

}

func (s *Sdk) GetEnableStatus() bool {
	if global.AppConfig.Enable <= 0 {
		return false
	} else {
		return true
	}

}

func (s *Sdk) Enable() bool {

	global.AppConfig.Enable = 1
	return true

}

func (s *Sdk) Disable() bool {

	global.AppConfig.Enable = 0
	return true

}

func (s *Sdk) AddCron(expression string, exefnc func()) bool {

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

func (s *Sdk) RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}

	return true

}

func (s *Sdk) Get(key string) bool {

	return true

}

func (s *Sdk) Set(key string, value string) bool {

	return true

}

func (s *Sdk) GetSession(sid string) *z_type.Session  {
	session, _ := global.SyncUserSessions.Get(sid)
	return session.(*z_type.Session)
}

func (s *Sdk) GetSessionStr(sid string)  string {


	user_session, exist := global.SyncUserSessions.Get(sid)
	js1 := []byte(`{}`)
	if exist {
		js1, _ = json_orgin.Marshal(user_session)
	}
	return string(js1)

}

func (s *Sdk) Kick(sid string) bool {

	return true

}

func (s *Sdk) CreateChannel(id string, name string) bool {

	area.CreateChannel(id, name)
	return true
}

func (s *Sdk) RemoveChannel(id string) bool {

	return true
}

func (s *Sdk) GetChannels() bool {

	return true
}

func (s *Sdk) GetSidsByChannel(channel_id string) bool {

	return true

}

func (s *Sdk) ChannelAddSid(sid string, area_id string) bool {


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

func (s *Sdk) ChannelKickSid(sid string, area_id string) bool {

	area.UnSubscribeChannel( area_id,sid)
	return true

}

func (s *Sdk) Push(from_sid string, to_sid string, msg string) bool {

	area.Push(to_sid, from_sid, msg)
	return true

}

func (s *Sdk) PushBySids(from_sid string,to_sids []string, msg string) bool {

	for _,to_sid:=   range to_sids {
		area.Push(to_sid, from_sid, msg)
	}
	return true

}


func (s *Sdk) BroadcastAll(msg string) bool {

	return true

}
