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


// 所有的场景名称列表
var Areas = make([]string, 0, 1000)

// 场景集合
var AreasMap *syncmap.SyncMap

// 一个全局的场景
var GlobalArea   *AreaType


// 所有的用户连接对象
var AllConns *syncmap.SyncMap
var AllWsConns *syncmap.SyncMap


type AreaType struct {

	Id string
	Name string
	// 当前场景包含的socket连接对象
	Conns *syncmap.SyncMap
	// 当前场景包含的websocket连接对象
	WsConns *syncmap.SyncMap

}


// 预创建多个场景
func InitConfigAreas() {

	AreasMap   = syncmap.New()
	AllConns   = syncmap.New()
	AllWsConns = syncmap.New()

	for _, area_id := range global.Config.Area.Init_area {
		CreateChannel(area_id, area_id)
	}
	GlobalArea = new(AreaType)
	GlobalArea.Id = "global"
	GlobalArea.Name = "global"
	GlobalArea.Conns = syncmap.New()
	GlobalArea.WsConns = syncmap.New()
	AreasMap.Set("global",GlobalArea)
}

// 创建一个场景
func CreateChannel(area_id string, name string) {
	golog.Info(area_id, name)
	Areas = append(Areas, area_id)
	area := new(AreaType)
	area.Id = area_id
	area.Name = name
	area.WsConns = syncmap.New()
	area.Conns = syncmap.New()
}

// 删除一个场景
func RemovChannel(id string) {
	golog.Info(id)
	// 1.删除名称
	for index, elem := range Areas {
		if elem==id {
			Areas = append(Areas[:index],Areas[index+1:]...)
			return
		}
	}
	// 删除场景对象
	AreasMap.Delete( id )
}

// 检查是否已经创建了场景
func CheckChannelExist(area_id string) bool {

	if ( AreasMap.Has(area_id) ) {
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
	var area  AreaType
	_item,ok := AreasMap.Get(area_id)
	if( !ok ) {
		golog.Error( "Channel  ",area_id," no exist! "  )
		return
	}else{
		area = _item.(AreaType)
		if( area.Conns.Size()<=0 ){
			area.Conns = syncmap.New()
		}
		if  !area.Conns.Has(sid) {
			area.Conns.Set(sid, conn)
		}
		AreasMap.Set( area_id,area )
	}
}

/**
 *  websocket连接 加入到场景中
 */
func SubscribeWsChannel(area_id string, ws *websocket.Conn, sid string) {

	var area  AreaType
	_item,ok := AreasMap.Get(area_id)
	if( !ok ) {
		golog.Error( "Channel  ",area_id," no exist! "  )
		return
	}else{
		area = _item.(AreaType)
		if( area.WsConns.Size()<=0 ){
			area.WsConns = syncmap.New()
		}
		if  !area.WsConns.Has(sid) {
			area.WsConns.Set(sid, ws)
		}
		AreasMap.Set( area_id,area )
	}
}


func GetSidsByChannel(id string) []string {

	ret := make([]string,0)
	var area  AreaType
	item,ok:= AreasMap.Get(id)
	if( ok ){
		area = item.(AreaType)
		for tmp := range area.Conns.IterItems(){
			ret=append(ret,tmp.Key)
		}
		for tmp := range area.WsConns.IterItems(){
			ret=append(ret,tmp.Key)
		}
	}
	return ret

}



/**
 *  检查用户是否加入到场景中
 */
func CheckUserJoinChannel(area_id string, sid string) bool {

	var area  AreaType
	_item,ok:= AreasMap.Get(area_id)
	if( ok ) {
		area = _item.(AreaType)
		if  area.Conns.Has(sid) {
			return true
		}
		if  area.WsConns.Has(sid) {
			return true
		}
	}
	return false

}


/**
 *  用户退出某个场景
 */
func UnSubscribeChannel(area_id string, sid string) {

	var area  AreaType
	_item,ok:= AreasMap.Get(area_id)
	if( ok ) {
		area = _item.(AreaType)
		area.Conns.Delete( sid )
		area.WsConns.Delete( sid )
		AreasMap.Set( area_id, area )
	}

}

// 用户退出所有场景
func UserUnSubscribeChannel(user_sid string) {

	for index, _ := range Areas {
		UnSubscribeChannel(Areas[index], user_sid)
	}
	UnSubGlobalChannel( user_sid )
}

/**
 *  在场景中广播消息
 */
func Broatcast( sid string,area_id string, msg []byte ) {

	fmt.Println("Broatcast:", sid, area_id, string(msg) )

	var area AreaType
	_item,ok := AreasMap.Get(area_id)
	if( !ok ) {
		return
	}
	area = _item.( AreaType )
	var conn *net.TCPConn
	fmt.Println("广播里有:", area.Conns.Size(),"个连接")
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	// socket部分
	for item := range area.Conns.IterItems() {
		//fmt.Println("key:", item.Key, "value:", item.Value)
		conn = item.Value.(*net.TCPConn)
		//fmt.Println( protocol.WrapBroatcastRespStr(sid,area_id,msg) )

		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		buf,_ := protocolPacket.WrapBroatcastResp( area_id, sid, msg  )
		conn.Write( buf )
	}

	// websocket部分
	fmt.Println("WS广播里有:", area.WsConns.Size(),"个连接")
	var wsconn *websocket.Conn
	for item := range area.WsConns.IterItems() {

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
	fmt.Println("广播里有:", GlobalArea.Conns.Size(),"个conn连接")
	protocolJson := new(protocol.Json)
	protocolJson.Init()
	for item := range GlobalArea.Conns.IterItems() {
		fmt.Println("key:", item.Key, "value:", item.Value)
		conn = item.Value.(*net.TCPConn)
		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		buf,_ := protocolPacket.WrapBroatcastResp( "global", sid, msg  )
		conn.Write( buf )
	}
	fmt.Println("广播里有:", GlobalArea.Conns.Size(),"个ws连接")
	var wsconn *websocket.Conn
	for item := range GlobalArea.WsConns.IterItems() {
		fmt.Println("key:", item.Key, "value:", item.Value)
		wsconn = item.Value.(*websocket.Conn)
		buf, _ := json.Marshal(protocolJson.WrapBroatcastRespObj( "global", sid, msg) )
		go wsconn.Write( buf )
	}
}

func UnSubGlobalChannel( sid string ) {

	GlobalArea.Conns.Delete( sid )
	GlobalArea.WsConns.Delete( sid )


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
	ws:=GetWsConn(to_sid)
	fmt.Println( "push, to_sid:", to_sid , to_data)
	if( ws!=nil ) {
		buf, _ := json.Marshal(protocolJson.WrapPushRespObj( to_sid, from_sid,to_data) )
		_,err:=ws.Write( buf )
		if err!=nil {
			fmt.Println( "wsconn.Write err:",err.Error() )
		}
		return
	}
}




func GetConn(sid string) *net.TCPConn {

	conn, ok := AllConns.Get(sid)
	if !ok {
		return nil
	} else {
		return conn.(*net.TCPConn)
	}
}

func DeleteConn(sid string) {

	AllConns.Delete(sid)

}

func GetWsConn(sid string) *websocket.Conn {
	wsconn, ok := AllWsConns.Get(sid)
	if !ok {
		return nil
	} else {
		return wsconn.(*websocket.Conn)
	}
}

func DeleteWsConn(sid string) {

	AllWsConns.Delete(sid)

}

func DeleteUserssion(sid string) {

	global.SyncUserSessions.Delete(sid)

}

func ConnRegister(conn *net.TCPConn, user_sid string) {

	//SubscribeChannel("area-global", conn, user_sid)

	AllConns.Set( user_sid, conn )

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

	AllWsConns.Set( user_sid, ws )

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


func DeleteSession(sid string) {

	//_, session_exist := global.SyncUserSessions.Get(sid)
	//if session_exist {
		global.SyncUserSessions.Delete(sid)
	//}

}

func DeleteUserChannel(sid string) {

	global.SyncUserJoinedChannels.Delete(sid)

}

func FreeConn(conn *net.TCPConn, sid string) {

	conn.Close()
	golog.Warn("Sid closing:", sid)
	DeleteConn(sid)
	DeleteSession(sid)
	DeleteUserChannel(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	UserUnSubscribeChannel(sid)
	global.SyncUserConns.Delete(sid)

}

func FreeWsConn(ws *websocket.Conn, sid string) {

	//ws.Write([]byte{'E', 'O', 'F'})
	ws.Close()
	golog.Warn("Sid closing:", sid)
	DeleteWsConn(sid)
	DeleteSession(sid)
	DeleteUserChannel(sid)
	atomic.AddInt32(&global.SumConnections, -1)
	UserUnSubscribeChannel(sid)

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
