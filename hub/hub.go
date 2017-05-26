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
	"encoding/json"
	"morego/util"
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
		//buf, err := reader.ReadBytes('\n')
		buf ,err := protocol.Unpack( reader )
		if err != nil {
			//fmt.Println( "Hub handleWorker connection error: ", err.Error())
			// 超时处理
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {

			}
			closeHubConn(conn)
			break
		}
		go hubWorkeDispath( buf , conn)

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

	cmd,from_sid,reqid,data_buf :=protocol.ReadHubReq( msg )
	data := string( data_buf )

	api := new(Api)
	fmt.Println( "hubWorkeDispath cmd:", cmd )

	if cmd == "GetBase" {
		flatbuf := protocol.MakeHubResp(cmd,reqid,"",api.GetBase())
		wrote_buf,_:=protocol.Packet( flatbuf )
		n,errw := conn.Write( wrote_buf )
		if errw!=nil {
			fmt.Println( "hubWorkeDispath err:", errw.Error() )
		}
		fmt.Println( "hubWorkeDispath GetBase write:",n, api.GetBase()  )
		//conn.Close()
		return
	}
	if cmd == "GetEnableStatus" {
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",string(global.AppConfig.Enable))  )
		return
	}
	if cmd == "Enable" {
		global.AppConfig.Enable = 1
		conn.Write( protocol.MakeHubResp(cmd,reqid,"","1")  )
		return
	}
	if cmd == "Disable" {
		global.AppConfig.Enable = 0
		conn.Write( protocol.MakeHubResp(cmd,reqid,"","1")  )
		return
	}
	if cmd == "Get" {
		str,err:=Get(data)
		if( err!=nil ) {
			conn.Write([]byte(protocol.MakeHubResp(cmd,reqid,err.Error(),"")))
			return
		}
		conn.Write([]byte(protocol.MakeHubResp(cmd,  reqid, "", str)))
		return
	}

	if cmd == "Set" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
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
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}

	if cmd == "Kick" {
		ret :=api.Kick(data)
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}

	if cmd == "CreateChannel" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
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
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return

	}

	if cmd == "RemoveChannel" {
		ret :=api.RemoveChannel(data)
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}

	if cmd == "GetChannels" {
		ret :=api.GetChannels()
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",string(ret) )  )
		return
	}

	if cmd == "GetSidsByChannel" {
		ret :=api.GetSidsByChannel( data )
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",string(ret) )  )
		return
	}

	if cmd == "ChannelAddSid" {
		fmt.Println("ChannelKickSid", data )
		data_buf = util.TrimX001( data_buf )
		var map_data map[string]string
		err_json := json.Unmarshal( data_buf ,&map_data )
		//data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			golog.Error("Hub ChannelAddSid json Unmarshal err:",err_json.Error())
			return
		}
		sid ,_ok1:= map_data["sid"]
		area_id ,_ok2:= map_data["area_id"]
		if( !_ok1 )  {
			golog.Error("Hub ChannelAddSid json sid no found" )
			return
		}
		if( !_ok2 )  {
			golog.Error("Hub ChannelAddSid json area_id no found"  )
			return
		}
		ret :=api.ChannelAddSid(sid, area_id )
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}
	if cmd == "ChannelKickSid" {

		data_buf = util.TrimX001( data_buf )
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
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
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}

	if cmd == "Push" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			golog.Error("Hub Push json err:",err_json.Error())
			return
		}
		to_sid,err2 := data_json.GetString("sid")
		if err2!=nil    {
			golog.Error("Hub Push json err:",err2.Error())
			return
		}

		ret := api.Push( from_sid, to_sid ,string(data_buf) )

		str :="0"
		if ret {
			str = "1"
		}
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}

	if cmd == "BroadcastAll" {
		ret :=api.BroadcastAll(data_buf)
		str :="0"
		if ret{
			str = "1"
		}
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}


	if cmd == "Broatcast" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			golog.Error("Hub Broatcast json err:",err_json.Error())
			return
		}
		area_id,err2 := data_json.GetString("area_id")
		if(  err2!=nil )  {
			golog.Error("Hub data_json json err:",err2.Error() )
			return
		}
		ret := api.Broadcast( from_sid, area_id ,data_buf )
		str := "0"
		if ret{
			str = "1"
		}
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}

	if cmd == "UpdateSession" {

		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
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
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",str )  )
		return
	}

	if cmd == "GetUserJoinedChannel" {
		data_json ,err_json:= jason.NewObjectFromBytes( data_buf )
		if( err_json!=nil ) {
			err_str :="Hub UpdateSession json err:"+err_json.Error()
			golog.Error( err_str )
			conn.Write( protocol.MakeHubResp(cmd,reqid,err_json.Error(),"" )  )
			return
		}
		sid, _ := data_json.GetString("sid")
		ret :=api.GetUserJoinedChannel(sid )

		conn.Write( protocol.MakeHubResp(cmd,reqid,"",string(ret) )  )
		return

	}

	if cmd == "GetAllSession" {

		ret :=api.GetAllSession()
		conn.Write( protocol.MakeHubResp(cmd,reqid,"",string(ret) )  )
		return

	}


}


