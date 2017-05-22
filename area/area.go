/**
 *  场景管理
 *  创建多个channel,一个channel对应一个publisher,chanel从Hub订阅消息后分发给客户端
 *
 */

package area

import (
	"fmt"
	"net"
	"time"
	//"strings"
	"sync/atomic"
	"math/rand"
	"morego/global"
	"morego/golog"
	"morego/lib/websocket"
	"morego/lib/syncmap"
	"morego/protocol"
	z_type "morego/type"
	"encoding/json"
)

// 预创建多个场景
func InitArea() {
	for _, area_id := range global.Config.Area.Init_area {
		CreateChannel(area_id, area_id)
		global.Channels[area_id] = global.Config.Hub.Hub_host
	}
}

// 创建一个RPC的场景
func CreateChannel(area_id string, name string) {
	golog.Info(area_id, name)
	global.RpcChannels = append(global.RpcChannels, area_id)
	global.SyncRpcChannelConns.Set(area_id,syncmap.New())
	global.SyncRpcChannelWsConns.Set(area_id,syncmap.New())
	//fmt.Println(global.RpcChannels)
}

// 删除一个RPC的场景
func RemovChannel(id string) {
	golog.Info(id)
	// @todo

}

// 检查是否已经创建了场景
func CheckChannelExist(area_id string) bool {

	if ( global.SyncRpcChannelConns.Has(area_id) ) {
		return true
	}
	if ( global.SyncRpcChannelWsConns.Has(area_id) ) {
		return true
	}
	return false

}

func  ChannelAddSid(sid string, area_id string) bool {

	exist := CheckChannelExist(area_id)
	//fmt.Println( area_id," CheckChannelExist:", exist )
	if !exist {
		return false
	}

	// 检查会话用户是否加入过此场景
	//have_joined := area.CheckUserJoinChannel(area_id, sid)
	//fmt.Println( "have_joined:",sid, area_id, have_joined )
	// 如果还没有加入场景,则订阅
	//if !have_joined {
	user_conn := GetConn(sid)
	user_wsconn :=  GetWsConn(sid)
	//fmt.Println( "ChannelAddSid user_wsconn:",user_wsconn )
	// 会话如果属于socket
	if user_conn != nil {
		 SubscribeChannel(area_id, user_conn, sid)
	}
	// 会话如果属于websocket
	if user_wsconn != nil {
		 SubscribeWsChannel(area_id, user_wsconn, sid)
	}
	// 该用户加入过的场景列表
	var userJoinedChannels = make([]string, 0, 1000)
	tmp, ok := global.SyncUserJoinedChannels.Get(sid)
	if ok {
		userJoinedChannels = tmp.([]string)
	}
	userJoinedChannels = append(userJoinedChannels, area_id)
	global.SyncUserJoinedChannels.Set(sid, userJoinedChannels)
	//}
	return true

}

/**
 *  socket连接 加入到场景中
 */
func SubscribeChannel(area_id string, conn *net.TCPConn, sid string) {

	// tcp部分
	var channels *syncmap.SyncMap
	_item,ok := global.SyncRpcChannelConns.Get(area_id)
	if( !ok ) {
		golog.Error( "Channel  ",area_id," no exist! "  )
		return
	}else{
		channels = _item.(*syncmap.SyncMap)
		if( channels.Size()<=0 ){
			channels = syncmap.New()
		}
		if  !channels.Has(sid) {
			channels.Set(sid, conn)
		}
		global.SyncRpcChannelConns.Set( area_id, channels )
		//fmt.Println("Joined  ",area_id," size :", channels.Size() )
	}


}

/**
 *  websocket连接 加入到场景中
 */
func SubscribeWsChannel(area_id string, ws *websocket.Conn, sid string) {

	_item,ok := global.SyncRpcChannelWsConns.Get(area_id)
	if( !ok ) {
		//fmt.Println("Channel  ",area_id," no exist! "  )
		golog.Error( "Channel  ",area_id," no exist! "  )
		return
	}else{
		var channels *syncmap.SyncMap
		channels = _item.(*syncmap.SyncMap)
		if( channels.Size()<=0 ){
			channels = syncmap.New()
		}
		//if  !channels.Has(sid) {
			//fmt.Println("SubscribeWsChannel  ",sid, area_id, ws   )
			channels.Set(sid, ws)
		//}
		global.SyncRpcChannelWsConns.Set( area_id, channels )
		//fmt.Println("Joined  ",area_id," size :", channels.Size() )
	}

}


func GetSidsByChannel(channel_id string) []string {

	ret := make([]string,0)
	if( global.SyncRpcChannelConns.Has( channel_id ) ){
		var channel *syncmap.SyncMap
		item,ok:= global.SyncRpcChannelConns.Get(channel_id)
		if( ok ){
			channel = item.(*syncmap.SyncMap)
			for tmp := range channel.IterItems(){
				ret=append(ret,tmp.Key)
			}

		}
	}
	return ret

}

/**
 *  socket连接 加入到全局场景中
 */
func SubscribeGlobalChannel( conn *net.TCPConn, sid string) {


	if  !global.SyncGlobalChannelConns.Has(sid) {
		global.SyncGlobalChannelConns.Set(sid, conn)
		fmt.Println("global_channel_conns.Set:", sid, conn )
	}else{
		fmt.Println("global_channel_conns exist:", sid )

	}
	fmt.Println("Joined SyncRpcChannelConns size :", global.SyncGlobalChannelConns.Size() )

	//golog.Error(  " sid ", sid, " join ", area_id, global.SyncRpcChannelConns)

}


/**
 *  websocket连接 加入到全局场景中
 */
func SubscribeGlobalChannelWs( ws *websocket.Conn, sid string) {


	if  !global.SyncGlobalChannelWsConns.Has(sid) {
		global.SyncGlobalChannelWsConns.Set(sid, ws)
		fmt.Println("global_channel_conns.Set:", sid, ws )
	}else{
		fmt.Println("global_channel_conns exist:", sid )

	}
	fmt.Println("Joined SyncRpcChannelConns size :", global.SyncGlobalChannelWsConns.Size() )



}

/**
 *  检查用户是否加入到场景中
 */
func CheckUserJoinChannel(area_id string, sid string) bool {

	// tcp部分
	_item,ok := global.SyncRpcChannelConns.Get(area_id)
	if( ok ) {
		var channel_conns *syncmap.SyncMap
		channel_conns = _item.(*syncmap.SyncMap)
		if  channel_conns.Has(sid) {
			return true
		}
	}

	// websocket部分
	_item_ws,okws := global.SyncRpcChannelWsConns.Get(area_id)
	if( okws ) {
		var channel_wsconns *syncmap.SyncMap
		channel_wsconns = _item_ws.(*syncmap.SyncMap)
		if  channel_wsconns.Has(sid) {
			return true
		}
	}
	return false

}


/**
 *  用户推出某个场景
 */
func UnSubscribeChannel(area_id string, sid string) {


	// tcp部分
	_item,ok := global.SyncRpcChannelConns.Get(area_id)
	if( ok ) {
		var channel_conns *syncmap.SyncMap
		channel_conns = _item.(*syncmap.SyncMap)
		channel_conns.Delete(sid)
		global.SyncRpcChannelConns.Set( area_id,channel_conns )

	}

	// websocket部分
	_item_ws,okws := global.SyncRpcChannelWsConns.Get(area_id)
	if( okws ) {
		var channel_wsconns *syncmap.SyncMap
		channel_wsconns = _item_ws.(*syncmap.SyncMap)
		channel_wsconns.Delete(sid)
		global.SyncRpcChannelWsConns.Set( area_id,channel_wsconns )

	}
}

// 用户退出所有场景
func UserUnSubscribeChannel(user_sid string) {

	for index, _ := range global.RpcChannels {
		UnSubscribeChannel(global.RpcChannels[index], user_sid)
	}
	UnSubGlobalChannel( user_sid )
}

/**
 *  在场景中广播消息
 */
func Broatcast( sid string,area_id string, msg []byte ) {

	fmt.Println("Broatcast:", sid, area_id, string(msg) )
	// tcp部分
	var channel_conns *syncmap.SyncMap
	_item,ok := global.SyncRpcChannelConns.Get(area_id)
	if( !ok ) {
		return
	}
	channel_conns = _item.(*syncmap.SyncMap)
	var conn *net.TCPConn
	fmt.Println("广播里有:", channel_conns.Size(),"个连接")
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	for item := range channel_conns.IterItems() {
		//fmt.Println("key:", item.Key, "value:", item.Value)
		conn = item.Value.(*net.TCPConn)
		//fmt.Println( protocol.WrapBroatcastRespStr(sid,area_id,msg) )

		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		buf,_ := protocolPacket.WrapBroatcastResp( area_id, sid, msg  )
		conn.Write( buf )
	}

	// websocket部分
	var channel_wsconns *syncmap.SyncMap
	_item_ws,ok := global.SyncRpcChannelWsConns.Get(area_id)
	if( !ok ) {
		return
	}
	channel_wsconns = _item_ws.(*syncmap.SyncMap)

	fmt.Println("WS广播里有:", channel_wsconns.Size(),"个连接")
	var wsconn *websocket.Conn
	for item := range channel_wsconns.IterItems() {

		wsconn = item.Value.(*websocket.Conn)
		buf, _ := json.Marshal(protocolJson.WrapBroatcastRespObj( area_id, sid, msg) )
		fmt.Println( "WrapBroatcastRespObj:", string(buf) )
		write_len,err:= wsconn.Write( buf )
		if err!=nil {
			fmt.Println("广播 err:", err.Error())
		}
		fmt.Println( "write_len:", write_len )
	}
}

/**
 *  在场景中广播消息
 */
func BroatcastGlobal( sid string, msg []byte ) {

	var conn *net.TCPConn
	fmt.Println("广播里有:", global.SyncGlobalChannelConns.Size(),"个连接")
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	for item := range global.SyncGlobalChannelConns.IterItems() {
		fmt.Println("key:", item.Key, "value:", item.Value)
		conn = item.Value.(*net.TCPConn)
		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		buf,_ := protocolPacket.WrapBroatcastResp( "global", sid, msg  )
		conn.Write( buf )
	}

	var wsconn *websocket.Conn
	for item := range global.SyncGlobalChannelWsConns.IterItems() {
		fmt.Println("key:", item.Key, "value:", item.Value)
		wsconn = item.Value.(*websocket.Conn)
		buf, _ := json.Marshal(protocolJson.WrapBroatcastRespObj( "global", sid, msg) )
		go wsconn.Write( buf )
	}
}

func UnSubGlobalChannel( sid string ) {

	global.SyncGlobalChannelConns.Delete(sid)
	global.SyncGlobalChannelWsConns.Delete(sid)
}

/**
 *  点对点发送消息
 */
func Push(  to_sid string ,from_sid string,to_data string ) {
	conn :=  GetConn(to_sid)
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	if( conn!=nil ) {
		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		buf,err := protocolPacket.WrapPushResp(  to_sid, from_sid,[]byte(to_data) )
		if err!=nil {
			fmt.Println( "protocolPacket.WrapPushResp:",err.Error() )
		}
		_,err =conn.Write( buf )
		if err!=nil {
			fmt.Println( "Push conn.Write err:",err.Error() )
		}
		return
	}
	wsconn:=GetWsConn(to_sid)
	fmt.Println( "push, to_sid:", to_sid , to_data)
	if( wsconn!=nil ) {
		buf, _ := json.Marshal(protocolJson.WrapPushRespObj( to_sid, from_sid,to_data) )
		_,err:=wsconn.Write( buf )
		if err!=nil {
			fmt.Println( "wsconn.Write err:",err.Error() )
		}
		return
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

	//SubscribeChannel("area-global", conn, user_sid)

	//_, ok := global.SyncUserConns.Get(user_sid)
	//if !ok {
		global.SyncUserConns.Set(user_sid, conn)
	//}

	_, ok := global.SyncUserSessions.Get(user_sid)
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
	//SubscribeWsChannel("area-global", ws, user_sid)

	//_, ok := global.SyncUserWebsocketConns.Get(user_sid)
	//if !ok {
		global.SyncUserWebsocketConns.Set(user_sid, ws)
	//}

	_, ok := global.SyncUserSessions.Get(user_sid)
	if !ok {
		data := &z_type.Session{
			ws.RemoteAddr().String(),
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


func CloseWsConn(sid string) {
	_, conn_exist := global.SyncUserWebsocketConns.Get(sid)
	if conn_exist {
		global.SyncUserWebsocketConns.Delete(sid)
	}
}

func CloseConn(sid string) {

	_, conn_exist := global.SyncUserConns.Get(sid)
	if conn_exist {
		global.SyncUserConns.Delete(sid)
	}

}

func CloseSession(sid string) {

	_, session_exist := global.SyncUserSessions.Get(sid)
	if session_exist {
		global.SyncUserSessions.Delete(sid)
	}

}

func CloseUserChannel(sid string) {

	global.SyncUserJoinedChannels.Delete(sid)

}

func FreeConn(conn *net.TCPConn, sid string) {

	conn.Write([]byte{'E', 'O', 'F'})
	conn.Close()
	golog.Warn("Sid closing:", sid)
	CloseConn(sid)
	CloseSession(sid)
	CloseUserChannel(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	UserUnSubscribeChannel(sid)
	global.SyncUserConns.Delete(sid)
	golog.Info("UserConns length:", len(global.UserConns))

}

func FreeWsConn(ws *websocket.Conn, sid string) {

	//ws.Write([]byte{'E', 'O', 'F'})
	ws.Write( []byte{'E', 'O', 'F'} )
	ws.Close()
	golog.Warn("Sid closing:", sid)
	CloseWsConn(sid)
	CloseSession(sid)
	CloseUserChannel(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	UserUnSubscribeChannel(sid)
	golog.Info("UserConns length:", len(global.UserConns))

}

/**
 * 检查
 */
func CheckSid(sid string) bool {

	return true
	_, exist := global.SyncUserSessions.Get(sid)
	return exist
}

func CreateSid() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sid := fmt.Sprintf("%d%d", r.Intn(99999), rand.Intn(999999))
	return sid
}
