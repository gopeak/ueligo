package golang

import (
	"morego/golog"
	"github.com/robfig/cron"
	"morego/global"
	z_type "morego/type"
	"fmt"
	"time"
	"net"
	"morego/protocol"
	"morego/hub"
	"bufio"
	"morego/lib/syncmap"
	"morego/util"
	"strconv"
)


type Sdk struct {

	Connected bool

	HubConn *net.TCPConn

	Cmd string

	Sid string

	Reqid int

	Data []byte

}

type PushReqHub struct {

	Sid bool
	Msg string
	Info map[string]string

}

type AfterWorkCallback func(   resp_buf string ) (string)

var ReqSeqCallbacks *syncmap.SyncMap



var ReqHubConns  =  make( []*net.TCPConn, 0 )
var InitialCap  int

func (sdk *Sdk) Init(cmd string,sid string,reqid int,data []byte) *Sdk{

	sdk.Cmd = cmd
	sdk.Sid = sid
	sdk.Reqid = reqid
	sdk.Data = data
	sdk.Connected = false
	return sdk
}

// 数据连接
func (sdk *Sdk) connect() bool{

	if sdk.HubConn!=nil {
		return true
	}
	data :=  global.Config.WorkerServer.ToHub
	hub_host := data[0]
	hub_port_str := data[1]
	ip_port := hub_host + ":" + hub_port_str

	tcpAddr, _ := net.ResolveTCPAddr("tcp4", ip_port)
	hubconn, err_req := net.DialTCP("tcp", nil, tcpAddr)
	if( err_req!=nil ){
		sdk.HubConn=nil
		return false
	}
	sdk.HubConn = hubconn
	return true

}


func   InitReqHubPool() {

	// create a factory() to be used with channel based pool
	ReqSeqCallbacks = syncmap.New()

	InitialCap  = 10
	fmt.Println( "global.Config.WorkerServer",global.Config.WorkerServer.Servers)

	factory    := func() (*net.TCPConn, error) {
		data :=  global.Config.WorkerServer.ToHub
		hub_host := data[0]
		hub_port_str := data[1]
		ip_port := hub_host + ":" + hub_port_str

		tcpAddr, _ := net.ResolveTCPAddr("tcp4", ip_port)
		hubconn, err_req := net.DialTCP("tcp", nil, tcpAddr)
		//fmt.Println( "InitConnectionHubPool hubconn ", hubconn )

		return hubconn,err_req
	}
	for i := 0; i < InitialCap; i++ {

		var err_req error
		conn, err_req:= factory()
		if( err_req!=nil ) {
			golog.Error( "InitConnectionHubPool hubconn  err:", err_req.Error() )
			continue
		}
		ReqHubConns = append( ReqHubConns, conn )
		go handleReqHubResponse( conn )
	}
}


// 侦听Hub server返回的数据，然后回调worker的函数
func  handleReqHubResponse(conn *net.TCPConn) {

	time.Sleep( 2*time.Second)
	reader := bufio.NewReader(conn)
	defer func() {
		err := recover()
		if err != nil {
			conn.Close()
			fmt.Println( "ReadHubResp err :", err)
		}
	}()
	for {
		buf ,err := protocol.Unpack( reader)
		if err != nil {
			fmt.Println( "handleReqHubResponse protocol.Unpack error: ", err.Error())
			conn.Close()
			break
		}
		resp_cmd,req_id,resp_err,msg_data := protocol.ReadHubResp(buf)
		if resp_err!=""{
			golog.Error( "handleReqHubResponse protocol.ReadHubResp err:",resp_err )
			continue
		}
		callback_key:=resp_cmd + req_id
		if req_id=="0"{
			ReqSeqCallbacks.Delete( callback_key )
			continue

		}

		fmt.Println( "callback_key:", callback_key )
		fmt.Println( "callback data:", string(msg_data)  )
		_item,ok := ReqSeqCallbacks.Get( callback_key )
		if( ok ) {
			callback := _item.( AfterWorkCallback )
			fmt.Println( "callback func :", callback  )
			callback( string(msg_data) )
			ReqSeqCallbacks.Delete( callback_key )
		}
	}

}


// 向Hub请求数据并监听返回,该请求将会阻塞除非等待返回超时
func (sdk *Sdk) ReqHubAsync( req_cmd string , data string ,handler AfterWorkCallback  ) (string,bool) {

	req_id := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)
	req_buf := protocol.MakeHubReq( req_cmd, sdk.Sid, req_id, data )
	req_buf,_ = protocol.Packet( req_buf )

	index := util.RandInt64(0, int64(len(ReqHubConns)))
	req_hub_conn  := ReqHubConns[index]

	//req_hub_conn,err := HubConnsPool.Get()
	//fmt.Println( "ReqHubConns:", ReqHubConns )
	if( req_hub_conn==nil  ){
		golog.Error( "req_hub_conn is nil "  )
		return "", false
	}
	callback_key:=req_cmd+ req_id
	ReqSeqCallbacks.Set( callback_key, handler )
	//fmt.Println( "ReqHubAsync:", callback_key )
	_,err := req_hub_conn.Write( req_buf )
	if err!=nil {
		golog.Error( "req_hub_conn.Write err:" , err.Error() )
	}
	return "",false
}


// 向Hub请求数据并监听返回,该请求将会阻塞除非等待返回超时
func (sdk *Sdk) ReqHub( req_cmd string , data string ) (string,bool) {

	req_id := strconv.FormatInt( time.Now().UTC().UnixNano(), 10)
	req_buf := protocol.MakeHubReq( req_cmd, sdk.Sid, req_id, data )
	req_buf,_ = protocol.Packet( req_buf )
	//fmt.Println( "req_str:", string(req_buf) )

	sdk.connect()
	_,err:=sdk.HubConn.Write( req_buf )
	if( err!=nil ) {
		return "sdk.HubConn.Write err",false
	}
	reader := bufio.NewReader(sdk.HubConn)
	defer func() {
		err := recover()
		if err != nil {
			sdk.HubConn.Close()
			fmt.Println( "ReqHub err :", err)
		}
	}()
	for {
		buf ,err := protocol.Unpack( reader)
		select {

		case <- time.After(5 * time.Second):
			return "time 5 second",false

		default:
			if err != nil {
				fmt.Println( "handleReqHubResponse protocol.Unpack error: ", err.Error())
				sdk.HubConn.Close()
				return err.Error(),false
			}
			resp_cmd,resp_req_id,resp_err,ret_buf := protocol.ReadHubResp(buf)

			if resp_cmd+resp_req_id == req_cmd+req_id{
				// 如果服务返回错误
				if( resp_err!=""  ){
					golog.Error( "ReqHub resp err:",resp_err)
					return resp_err,false
				}
				sdk.HubConn.Close()
				return string(ret_buf),true
			}
		}
	}

	return "",false
}

func (sdk *Sdk) PushHub( req_cmd string , data string ) bool {

	req_buf := protocol.MakeHubReq( req_cmd,sdk.Sid, strconv.Itoa( int(sdk.Reqid) ), data )
	req_buf,_ = protocol.Packet( req_buf )
	sdk.connect()
	_,err:=sdk.HubConn.Write( req_buf )
	if( err!=nil ) {
		return false
	}

	return true
}



// 获取服务器的根路径
func (sdk *Sdk)  GetBase() string {

	// 单机模式直接返回内存中数据
	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetBase()
	}

	ret,ok :=sdk.ReqHub( "GetBase","" )
	if ok {
		return ret
	}
	return ""

}

// 获取服务启用状态
func (sdk *Sdk) GetEnableStatus() bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetEnableStatus()
	}

	ret,ok:= sdk.ReqHub( "GetEnableStatus","" )
	if( !ok ){
		return false
	}
	if( ret=="1" ){
		return true
	}else{
		return false
	}

}

func (sdk *Sdk) Enable() bool {

	if( global.SingleMode ) {
		global.AppConfig.Enable = 1
		return true
	}
	return sdk.PushHub( "Enable","")


}

func (sdk *Sdk) Disable() bool {

	if( global.SingleMode ) {
		global.AppConfig.Enable = 0
		return true
	}
	return sdk.PushHub( "Disable","")

}

func (sdk *Sdk) AddCron(expression string, exefnc func()) bool {

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

func (sdk *Sdk) RemoveCron(expression string) bool {

	if cron, ok := global.Crons[expression]; ok {
		delete(global.Crons, expression)
		cron.Stop()
	} else {
		return false
	}

	return true

}

func (sdk *Sdk) Get(key string) string {

	if( global.SingleMode ) {
		str,err:=hub.Get(key)
		if err!=nil {
			golog.Error("Redis Get err:",err.Error())
			return ""
		}
		return str
	}

	ret,ok := sdk.ReqHub( "Get",key )
	if( !ok ) {
		return ""
	}
	return ret

}

func (sdk *Sdk) Set(key string, value string,expire int) bool {

	if( global.SingleMode ) {
		ret,err:=hub.Set(key,value,expire)
		if err!=nil {
			golog.Error("Redis Set err:",err.Error())
			return false
		}
		return ret
	}
	json:=fmt.Sprintf(`{"key":"%s","value":"%s","expire":%d}`,key,value,expire)
	ret:= sdk.PushHub( "Set",json )
	return ret
}

// 该方法仅在单机模式下调用
func (sdk *Sdk) GetSessionType(sid string) *z_type.Session  {

	session,exist := global.SyncUserSessions.Get(sid)
	if !exist {
		return nil
	}
	return session.(*z_type.Session)
}

func (sdk *Sdk) GetSession(sid string)  string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetSession( sid )
	}
	ret,ok := sdk.ReqHub( "GetSession",sid  )
	if !ok{
		return ""
	}
	return ret

}

func (sdk *Sdk) Kick(sid string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Kick( sid )
	}
	return sdk.PushHub( "Kick",sid)
}

func (sdk *Sdk) CreateChannel(id string, name string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.CreateChannel( id,name )
	}
	json:=fmt.Sprintf(`{"id":"%s","name":"%s","expire":%d}`,id,name)
	return sdk.PushHub( "CreateChannel",json)

}

func (sdk *Sdk) RemoveChannel(id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.RemoveChannel( id )
	}
	return sdk.PushHub( "RemoveChannel",id)
}

func (sdk *Sdk) GetChannels() string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetChannels(  )
	}
	ret,ok := sdk.ReqHub( "GetChannels","" )
	if( !ok ) {
		return "{}"
	}
	return ret
}



func (sdk *Sdk) GetSidsByChannel(channel_id string) string {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetSidsByChannel( channel_id )
	}
	ret,ok :=  sdk.ReqHub( "GetSidsByChannel",channel_id )
	if( !ok ) {
		return "{}"
	}
	return ret

}

func (sdk *Sdk) ChannelAddSid(sid string, area_id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.ChannelAddSid( sid, area_id  )
	}
	json:=fmt.Sprintf(`{"sid":"%s","area_id":"%s"}`,sid, area_id )
	return sdk.PushHub( "ChannelAddSid",json)

}

func (sdk *Sdk) ChannelKickSid( sid string, area_id string) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.ChannelKickSid( sid, area_id  )
	}
	json:=fmt.Sprintf(`{"sid":"%s","area_id":"%s"}`,sid, area_id )
	return sdk.PushHub( "ChannelKickSid",json)

}

func (sdk *Sdk) Push( from_sid string ,to_sid string , data  []byte ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Push ( from_sid,to_sid, string(data)  )
	}

	return sdk.PushHub( "Push",string(data) )

}

func (sdk *Sdk) PushBySids(from_sid string,to_sids []string, data []byte) bool {

	for _,to_sid:=   range to_sids {
		sdk.Push(from_sid, to_sid, data )
	}
	return true

}

func (sdk *Sdk) Broatcast(sid string ,area_id string,  data []byte ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.Broadcast( sid,area_id, data  )
	}
	return sdk.PushHub( "Broatcast",string(data) )

}


func (sdk *Sdk) BroadcastAll( msg []byte ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.BroadcastAll( msg )
	}
	return sdk.PushHub( "BroadcastAll", string(msg) )

}


func (sdk *Sdk) UpdateSession( sid string, data string ) bool {

	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.UpdateSession( sid, data )
	}
	json:=fmt.Sprintf(`{"sid":"%s","data":"%s"}`,sid, data )
	return sdk.PushHub( "UpdateSession",json)

}

func (sdk *Sdk)GetUserJoinedChannel(sid string) string {

	// 单机模式直接返回内存中数据
	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetUserJoinedChannel(sid)
	}

	ret,ok :=sdk.ReqHub( "GetUserJoinedChannel",sid)
	if ok {
		return ret
	}
	return ""

}

func (sdk *Sdk)GetAllSession( ) string {

	// 单机模式直接返回内存中数据
	if( global.SingleMode ) {
		api := new(hub.Api)
		return api.GetAllSession()
	}

	ret,ok :=sdk.ReqHub( "GetAllSession","")
	if ok {
		return ret
	}
	return ""

}
