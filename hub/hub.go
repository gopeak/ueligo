//
//  Hub server
//
//

package hub

import (
	"bufio"
	"fmt"
	"net"
	"morego/global"
	"morego/golog"
	"github.com/antonholmquist/jason"
	"morego/protocol"
	"strconv"
	"time"
	"strings"
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
		//fmt.Println("handleHub  from :" , msg)
		if( strings.Replace(string(msg), "\n", "", -1)==""){
			continue
		}
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
	fmt.Println( "hubDispath str:", string(msg))
	msg_err,cmd,sid,reqid,data := protocol.ParseHubReqData(string(msg))
	if( msg_err!=nil ){
		fmt.Println( "hubDispath err:",msg_err.Error(),cmd,sid,reqid,data )
		return
	}
	api := new(Api)
	fmt.Println( "hubWorkeDispath cmd:", cmd )

	if cmd == "GetBase" {
		conn.Write([]byte(string(global.AppConfig.Enable)))
		conn.Close()
		return
	}
	if cmd == "GetEnableStatus" {
		conn.Write([]byte(string(global.AppConfig.Enable)))
	}
	if cmd == "Enable" {
		global.AppConfig.Enable = 1
		conn.Write([]byte(string(global.AppConfig.Enable)))
	}
	if cmd == "Disable" {
		global.AppConfig.Enable = 0
		conn.Write([]byte(string(`1`)))
	}
	if cmd == "Get" {
		str,err:=Get(data)
		if( err!=nil ) {
			conn.Write([]byte(protocol.WrapHubRespErrStr(err.Error(),cmd)))
			return
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "Set" {
		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			golog.Error("Hub Set json err:",err_json.Error())
			return
		}
		key,err_key := data_json.GetString("key")
		value,err_v := data_json.GetString("value")
		expire,err_e := data_json.GetInt64("expire")
		if( err_key!=nil || err_v!=nil || err_e!=nil ){
			golog.Error("Hub Set json err:",err_key.Error()+err_v.Error()+err_e.Error())
			return
		}
		_,err:=Set(key,value,expire)
		if( err!=nil ) {
			golog.Error("Hub Set err:",err.Error())
		}
		return

	}

	if cmd == "GetSession" {
		str :=api.GetSession(data)
		fmt.Println( "api.GetSession:",str)
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "Kick" {
		ret :=api.Kick(data)
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "CreateChannel" {
		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			golog.Error("Hub Set json err:",err_json.Error())
			return
		}
		id,err1 := data_json.GetString("id")
		name,err2 := data_json.GetString("name")
		if( err1!=nil || err2!=nil )  {
			golog.Error("Hub Set json err:",err1.Error()+err2.Error() )
			return
		}
		ret:=api.CreateChannel( id, name )
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return

	}

	if cmd == "RemoveChannel" {
		ret :=api.RemoveChannel(data)
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "GetChannels" {
		ret :=api.GetChannels()
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,string(ret))))
		return
	}

	if cmd == "GetSidsByChannel" {
		ret :=api.GetSidsByChannel( data )
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,string(ret))))
		return
	}

	if cmd == "ChannelAddSid" {
		fmt.Println("ChannelKickSid", data )
		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			golog.Error("Hub ChannelAddSid json err:",err_json.Error())
			return
		}
		sid,err1 := data_json.GetString("sid")
		area_id,err2 := data_json.GetString("area_id")
		if( err1!=nil || err2!=nil )  {
			golog.Error("Hub ChannelAddSid json err:",err1.Error()+err2.Error() )
			return
		}
		ret :=api.ChannelAddSid(sid, area_id )
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}
	if cmd == "ChannelKickSid" {

		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			golog.Error("Hub ChannelKickSid json err:",err_json.Error())
			return
		}
		sid,err1 := data_json.GetString("sid")
		area_id,err2 := data_json.GetString("area_id")
		if( err1!=nil || err2!=nil )  {
			golog.Error("Hub ChannelKickSid json err:",err1.Error()+err2.Error() )
			return
		}
		ret :=api.ChannelKickSid(sid, area_id )
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "Push" {
		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			golog.Error("Hub Push json err:",err_json.Error())
			return
		}
		from_sid,err1 := data_json.GetString("from_sid")
		to_sid,err2 := data_json.GetString("to_sid")
		to_data,err2 := data_json.GetString("msg")
		if( err1!=nil || err2!=nil )  {
			golog.Error("Hub Push json err:",err1.Error()+err2.Error() )
			return
		}
		ret :=api.Push(from_sid, to_sid ,to_data )
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "BroadcastAll" {
		ret :=api.BroadcastAll(data)
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}


	if cmd == "Broatcast" {
		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			golog.Error("Hub Broatcast json err:",err_json.Error())
			return
		}
		sid,err1 := data_json.GetString("sid")
		area_sid,err2 := data_json.GetString("area_sid")
		to_data,err2 := data_json.GetString("msg")
		if( err1!=nil || err2!=nil )  {
			golog.Error("Hub data_json json err:",err1.Error()+err2.Error() )
			return
		}
		ret :=api.Broadcast(sid, area_sid ,to_data )
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "UpdateSession" {

		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			golog.Error("Hub UpdateSession json err:",err_json.Error())
			return
		}
		sid,err1 := data_json.GetString("sid")
		to_data,err2 := data_json.GetString("data")
		if( err1!=nil || err2!=nil )  {
			golog.Error("Hub UpdateSession json err:",err1.Error()+err2.Error() )
			return
		}
		ret :=api.UpdateSession(sid, to_data )
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,str)))
		return
	}

	if cmd == "GetUserJoinedChannel" {
		data_json ,err_json:= jason.NewObjectFromBytes( []byte(data) )
		if( err_json!=nil ) {
			err_str :="Hub UpdateSession json err:"+err_json.Error()
			golog.Error( err_str )
			conn.Write([]byte(protocol.WrapRespErrStr(err_json.Error())))
			return
		}
		sid, _ := data_json.GetString("sid")
		ret :=api.GetUserJoinedChannel(sid )

		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,ret)))
		return

	}

	if cmd == "GetAllSession" {

		ret :=api.GetAllSession()
		conn.Write([]byte(protocol.WrapHubRespStr(cmd, sid, reqid,string(ret))))
		return

	}


}


