//
//  main
//

package main

import (

	"morego/global"
	"morego/golog"
	"morego/hub"
	//"strconv"
	//"net"
	//"net/http"
	_ "net/http/pprof"
	"runtime"
	//"morego/admin"
	"morego/area"
	"morego/connector"
	"morego/lib/syncmap"
	"morego/worker"

	//z_type "morego/type"
)



// 初始化全局变量
func init_global() {

	global.SumConnections = 0
	global.Qps = 0

	// 先在global声明,再使用make函数创建一个非nil的map，nil map不能赋值
	global.Channels = make(map[string]string)

	// global.RpcChannels  =  make(map[string] *z_type.ChannelRpcType )

	global.SyncUserConns = syncmap.New()
	global.SyncUserSessions = syncmap.New()
	global.SyncRpcChannelConns =  syncmap.New()
	global.SyncRpcChannelWsConns =  syncmap.New()

	global.SyncUserWebsocketConns = syncmap.New()
	global.SyncUserJoinedChannels = syncmap.New()
	global.SyncGlobalChannelConns = syncmap.New()
	global.SyncGlobalChannelWsConns = syncmap.New()

	global.SyncCrons = syncmap.New()
	global.PackSplitType = global.Config.PackType
	global.SingleMode = global.Config.SingleMode
	global.AuthCcmd = global.Config.Connector.AuthCcmd

	global.InitWorkerAddr()

}

/**
 * zeromore 框架启动
 */
func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	global.InitConfig()

	golog.InitLogger()

	init_global()
	go connector.SocketConnector("", global.Config.Connector.SocketPort)
	go connector.WebsocketConnector("", global.Config.Connector.WebsocketPort)

	// 开启hub服务器
	go hub.HubServer()

	// 预创建多个场景
	go area.InitArea()

	// 启动worker
	//go start_php_worker()
	go worker.InitWorkerServer()

	// 监控
	//go hub.TickWorkerServer()

	// demo应用依赖web服务器
	//go web.HttpServer()

	golog.Info("Server started!")

	// C:\gopath\mongodb\bin\mongod.exe --dbpath=C:\gopath\mongodb\data
	// D:\soft\MongoDB\bin\mongod.exe --dbpath=D:\soft\MongoDB\data
	select {}

}
