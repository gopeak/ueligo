/**
 *  定义全局变量
 *
 */

package global

import (
	"fmt"
	"net"
	"github.com/antonholmquist/jason"
	"morego/lib/robfig/cron"
	"morego/lib/syncmap"
	"github.com/gorilla/websocket"
	z_type "morego/type"
	"strings"
)

const (
	ERROR_PACKET_RATES    = `{"cmd":"error_","data":{"ret":503,"msg":"packet rate limit" }}`        //
	ERROR_MAX_CONNECTIONS = `{"cmd":"error_","data":{"ret":503,"msg":"Max connection limit" }}`     //
	ERROR_RESPONSE        = `{"cmd":"error_","data":{"ret":500,"msg":"RecvMessage error" }}`        //
	DISBALE_RESPONSE      = `{"cmd":"error_","data":{"ret":501,"msg":"Server has been stopped!" }}` //
)


const (
	DATA_REQ_CONNECT    = `1`
	DATA_REQ_MSG        = `2`
	DATA_WORKER_CONNECT = `3`
	DATA_WORKER_REPLY   = `4`
)


// 服务器当前状态
var AppConfig = &z_type.Appconfig{}

// worker连接数组
var WorkerNbrs []string

var WorkerServers = make([]string, 0, 1000)

var SumConnections int32

var Qps int64

// 所有的场景名称列表(name:bind)
var Channels = map[string]string{}

var RpcChannels = make([]string, 0, 1000)

var SyncRpcChannelConns  *syncmap.SyncMap

var SyncRpcChannelWsConns   *syncmap.SyncMap

var SyncRpcChannelSids = make([][]string, 0, 10090)

// 全局场景
var SyncGlobalChannelConns *syncmap.SyncMap
var SyncGlobalChannelWsConns *syncmap.SyncMap

// 会话用户的加入的场景列表
var UserChannels = map[string][]string{}

// 会话用户订阅的场景列表
var ConfigJson *jason.Object

// 用户连接对象
var UserConns = map[string]*net.TCPConn{}

// 用户连接对象
var SyncUserConns *syncmap.SyncMap

// 用户会话对象
var UserSessions = map[string]*z_type.Session{}

// 安全的用户会话对象
var SyncUserSessions *syncmap.SyncMap
var SyncUserJoinedChannels *syncmap.SyncMap
var UserWebsocketConns = map[string]*websocket.Conn{}
var SyncUserWebsocketConns *syncmap.SyncMap
var RpcType string

var SingleMode bool
var PackSplitType string
var AuthCmds []string

var SyncCrons *syncmap.SyncMap
var Crons = map[string]*cron.Cron{}

//var ReqAgentConns *syncmap.SyncMap

const (
	Splitstr = "||"
)

//  转义json字符串
func EncodeJsonStr(str string) string {
	str = strings.Replace(str, `"`, `\"`, -1)
	return str
}

// 反解json字符串
func DecodeJsonStr(str string) string {
	str = strings.Replace(str, `\"`, `"`, -1)
	return str
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Fatal error: %s", err.Error())
	}
}


func IsAuthCmd( cmd string ) bool {

	fmt.Println( "global.AuthCmds:",AuthCmds )
	for _,c:= range AuthCmds{
		if( c==cmd ){
			return true
		}
	}
	return false

}