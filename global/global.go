/**
 *  定义全局变量
 *
 */

package global

import (
	"fmt"
	"net"
	"github.com/robfig/cron"
	"morego/lib/syncmap"
	z_type "morego/type"
)

const (
	ERROR_PACKET_RATES    = `{"cmd":"error_","data":{"ret":503,"msg":"packet rate limit" }}`        //
	ERROR_MAX_CONNECTIONS = `{"cmd":"error_","data":{"ret":503,"msg":"Max connection limit" }}`     //
	ERROR_RESPONSE        = `{"cmd":"error_","data":{"ret":500,"msg":"RecvMessage error" }}`          //
	DISBALE_RESPONSE      = `{"cmd":"error_","data":{"ret":501,"msg":"Server has been stopped!" }}` //
)


// 服务器当前状态
var AppConfig = &z_type.Appconfig{}

var WorkerServers = make([]string, 0, 1000)

var SumConnections int32

var Qps int64

// 所有的场景名称列表(name:bind)
var Channels = map[string]string{}

var RpcChannels = make([]string, 0, 1000)

var SyncRpcChannelConns  *syncmap.SyncMap

var SyncRpcChannelWsConns   *syncmap.SyncMap

// 全局场景
var SyncGlobalChannelConns *syncmap.SyncMap
var SyncGlobalChannelWsConns *syncmap.SyncMap


// 用户连接对象
var UserConns = map[string]*net.TCPConn{}

// 用户连接对象
var SyncUserConns *syncmap.SyncMap

//  用户会话对象
var  SyncUserSessions *syncmap.SyncMap

var SyncUserJoinedChannels *syncmap.SyncMap
var SyncUserWebsocketConns *syncmap.SyncMap

var SingleMode bool
var PackSplitType string
var AuthCmds []string

var SyncCrons *syncmap.SyncMap
var Crons = map[string]*cron.Cron{}



func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error: %s", err.Error())
	}
}


func IsAuthCmd( cmd string ) bool {

	//fmt.Println( "global.AuthCmds:",AuthCmds )
	for _,c:= range AuthCmds{
		if( c==cmd ){
			return true
		}
	}
	return false

}