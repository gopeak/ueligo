//
//  main
//

package main

import (

	"ueligo/global"
	"ueligo/golog"
	"ueligo/hub"
	_ "net/http/pprof"
	"runtime"
	"ueligo/area"
	"ueligo/connector"
	"ueligo/lib/syncmap"
	"ueligo/worker"
)



// 初始化全局变量
func init_global() {

	global.SumConnections = 0
	global.Qps = 0

	// 先在global声明,再使用make函数创建一个非nil的map，nil map不能赋值
	global.AuthCmds = make([]string,0)
	global.UserSessions = syncmap.New()
	global.SingleMode = global.Config.SingleMode
	global.AuthCmds = global.Config.Connector.AuthCcmds
	area.UserJoinedAreas = syncmap.New()
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
	go area.InitConfig()

	// 启动worker
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
