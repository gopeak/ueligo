//
//  main
//

package main

import (
	"fmt"
	"morego/global"
	"morego/golog"
	"morego/hub"
	//"strconv"
	//"net"
	//"net/http"
	_ "net/http/pprof"
	"os"
	"os/exec"
	"runtime"
	//"morego/admin"
	"morego/area"
	"morego/connector"
	"morego/lib/syncmap"
	"morego/worker"
	"time"
	//z_type "morego/type"
)

// 启动一个测试的php worker以处理业务流程
func stop_php_worker() {

	c := exec.Command("/bin/sh", "-c", `ps -ef |grep "worker/php/workers.php"  |awk \'{print $2}\' |xargs -i kill -9 {} `)
	d, _ := c.Output()

	golog.Info("Stop_php_worker: ", string(d))

	time.Sleep(time.Second * 1)

}

// 启动一个测试的php worker以处理业务流程
func start_php_worker() {

	stop_php_worker()
	wd, _ := os.Getwd()
	work_num, _ := global.ConfigJson.GetString("worker", "worker_num")
	argv := []string{fmt.Sprintf("%s/worker/php/workers.php", wd), "start", work_num}
	golog.Info("Argv:", argv)
	c := exec.Command("/usr/bin/php", argv...)
	d, _ := c.Output()
	golog.Info("Start_php_worker: ", string(d))

	time.Sleep(time.Second * 1)

}

// 初始化全局变量
func init_global() {

	global.SumConnections = 0
	global.Qps = 0

	// 先在global声明,再使用make函数创建一个非nil的map，nil map不能赋值
	global.Channels = make(map[string]string)

	// global.RpcChannels  =  make(map[string] *z_type.ChannelRpcType )

	global.SyncUserConns = syncmap.New()
	global.SyncUserSessions = syncmap.New()
	global.SyncUserWebsocketConns = syncmap.New()
	global.SyncUserJoinedChannels = syncmap.New()
	global.SyncCrons = syncmap.New()
	global.PackSplitType = global.Config.PackType
	global.SingleMode = global.Config.SingleMode
	global.AuthCcmd = global.Config.Connector.AuthCcmd

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


	golog.Info("Server started!")

	// C:\gopath\mongodb\bin\mongod.exe --dbpath=C:\gopath\mongodb\data
	// D:\soft\MongoDB\bin\mongod.exe --dbpath=D:\soft\MongoDB\data
	select {}

}
