package worker

import (
	"morego/golog"
	"morego/lib/robfig/cron"
	//"fmt"
	"morego/global"
	z_type "morego/type"
	//"net"
	"fmt"
	"time"
	"net"
	"morego/protocol"
	"morego/hub"
	"bufio"
)


type Sdk struct {


	Connected bool

	HubConn *net.TCPConn

	Cmd string

	Sid string

	Reqid int

	Data string

	/*
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

	*/




}

func (this *Sdk) Init(cmd string,sid string,reqid int,data string) *Sdk{

	this.Cmd = cmd
	this.Sid = sid
	this.Reqid = reqid
	this.Data = data
	this.Connected = false
	return this
}



// 数据连接
func (this *Sdk) connect() bool{

	if this.HubConn!=nil {
		return true
	}
	data :=  global.Config.WorkerServer.ToHub
	hub_host := data[0]
	hub_port_str := data[1]
	ip_port := hub_host + ":" + hub_port_str

	hubconn, err_req := net.DialTimeout("tcp", ip_port, 5 * time.Second)
	if( err_req!=nil ){
		this.HubConn=nil
		return false
	}
	this.HubConn = hubconn
	return true

}

// 向Hub请求数据并监听返回,该请求将会阻塞除非等待返回超时
func (this *Sdk) ReqHub( req_cmd string , data string ) (string,bool) {

	req_str := protocol.WrapReqHubStr( req_cmd,this.Sid,this.Reqid,this.Data)
	this.HubConn.Write([]byte(req_str))
	reader := bufio.NewReader(this.HubConn)

	for {
		buf, err := reader.ReadBytes('\n')
		select {

		case <- time.After(5 * time.Second):
			return "time 5 second",false

		default:
			if err != nil {
				this.HubConn.Close()
				return err.Error(),false

			}
			errcode, _, resp_cmd, _, _, msg_data := protocol.ParseRplyData(string(buf))
			if( errcode==protocol.TypeError ){
				golog.Error( "ReqHub resp err:",msg_data)
				return msg_data,false
			}
			if resp_cmd == req_cmd{
				this.HubConn.Close()
				return msg_data,true
			}
		}
	}

	return "",false
}

func (this *Sdk) PushHub( req_cmd string , data string ) bool {

	req_str := protocol.WrapReqHubStr( req_cmd,this.Sid,this.Reqid,this.Data)
	_,err:=this.HubConn.Write([]byte(req_str))

	if( err!=nil ) {
		return false
	}

	return true
}

// 获取服务器的根路径
func (this *Sdk)  GetBase() string {

	// 单机模式直接返回内存中数据
	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetBase()
	}

	this.connect()

	ret,ok :=this.ReqHub( "GetBase","" )
	if ok {
		return ret
	}
	return ""

}



func (this *Sdk) GetEnableStatus() bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetEnableStatus()
	}

	ret,ok:= this.ReqHub( "GetEnableStatus","" )
	if( !ok ){
		return false
	}
	if( ret=="1" ){
		return true
	}else{
		return false
	}

}

func (this *Sdk) Enable() bool {

	if( global.SingleMode ) {
		global.AppConfig.Enable = 1
		return true
	}
	return this.PushHub( "Enable","")


}

func (this *Sdk) Disable() bool {

	if( global.SingleMode ) {
		global.AppConfig.Enable = 0
		return true
	}
	return this.PushHub( "Disable","")

}

func (this *Sdk) AddCron(expression string, exefnc func()) bool {

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

func (this *Sdk) RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}

	return true

}

func (this *Sdk) Get(key string) string {

	if( global.SingleMode ) {
		str,err:=hub.Get(key)
		if err!=nil {
			golog.Error("Redis Get err:",err.Error())
			return ""
		}
		return str
	}

	ret,ok := this.ReqHub( "Get",key )
	if( !ok ) {
		return ""
	}
	return ret

}

func (this *Sdk) Set(key string, value string,expire int) bool {

	if( global.SingleMode ) {
		ret,err:=hub.Set(key,value,expire)
		if err!=nil {
			golog.Error("Redis Set err:",err.Error())
			return false
		}
		return ret
	}
	json:=fmt.Sprintf(`{"key":"%s","value":"%s","expire":%d}`,key,value,expire)
	ret:= this.PushHub( "Set",json )
	return ret

}
// 该方法仅在单机模式下调用
func (this *Sdk) GetSessionType(sid string) *z_type.Session  {

	session,exist := global.SyncUserSessions.Get(sid)
	if !exist {
		return nil
	}
	return session.(*z_type.Session)
}

func (this *Sdk) GetSession(sid string)  string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetSession( sid )
	}
	ret,ok := this.ReqHub( "GetSession",sid )
	if !ok{
		return ""
	}
	return ret

}

func (this *Sdk) Kick(sid string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Kick( sid )
	}
	return this.PushHub( "Kick",sid)
}

func (this *Sdk) CreateChannel(id string, name string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.CreateChannel( id,name )
	}
	json:=fmt.Sprintf(`{"id":"%s","name":"%s","expire":%d}`,id,name)
	return this.PushHub( "CreateChannel",json)

}

func (this *Sdk) RemoveChannel(id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.RemoveChannel( id )
	}
	return this.PushHub( "RemoveChannel",id)
}

func (this *Sdk) GetChannels() string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetChannels(  )
	}
	ret,ok := this.ReqHub( "GetChannels","" )
	if( !ok ) {
		return "{}"
	}
	return ret
}



func (this *Sdk) GetSidsByChannel(channel_id string) string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetSidsByChannel( channel_id )
	}
	ret,ok :=  this.ReqHub( "GetSidsByChannel",channel_id )
	if( !ok ) {
		return "{}"
	}
	return ret

}

func (this *Sdk) ChannelAddSid(sid string, area_id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.ChannelAddSid( sid, area_id  )
	}
	json:=fmt.Sprintf(`{"sid":"%s","area_id":"%s"}`,sid, area_id )
	return this.PushHub( "ChannelAddSid",json)

}

func (this *Sdk) ChannelKickSid( sid string, area_id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.ChannelKickSid( sid, area_id  )
	}
	json:=fmt.Sprintf(`{"sid":"%s","area_id":"%s"}`,sid, area_id )
	return this.PushHub( "ChannelKickSid",json)

}

func (this *Sdk) Push(from_sid string, to_sid string, msg string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Push( from_sid, to_sid, msg  )
	}
	json:=fmt.Sprintf(`{"from_sid":"%s","to_sid":"%s","msg":"%s"}`,from_sid, to_sid, msg )
	return this.PushHub( "Push",json)

}

func (this *Sdk) PushBySids(from_sid string,to_sids []string, msg string) bool {

	for _,to_sid:=   range to_sids {
		this.Push(from_sid, to_sid, msg)
	}
	return true

}


func (this *Sdk) BroadcastAll(msg string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.BroadcastAll(   msg  )
	}
	return this.PushHub( "BroadcastAll",msg)

}