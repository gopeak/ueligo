package worker

import (
	//"math/rand"
	"morego/global"
	"morego/golog"
	"strconv"
	//"strings"
	//"sync/atomic"
	"morego/protocol"
	//sync"
	"bufio"
	"fmt"
	"net"
	//"os"
	"time"
	"encoding/json"
	"morego/lib/antonholmquist/jason"
)

/**
 * 监听客户端连接
 */
func WorkerServer(host string, port int) {

	fmt.Println("WorkerServer :", host, port)
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(host), (port), ""})
	if err != nil {
		golog.Error("ListenTCP Exception:", err.Error())
		return
	}

	// 处理客户端连接
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			golog.Error("AcceptTCP Exception::", err.Error(), time.Now().UnixNano())
			break
		}
		// 校验ip地址
		conn.SetKeepAlive(true)
		defer conn.Close()
		//conn.SetNoDelay(false)
		golog.Info("RemoteAddr:", conn.RemoteAddr().String())

		if( global.PackSplitType=="bufferio"){
			go handleWorker(conn)
		}
		if( global.PackSplitType=="json"){
			go handleWorkerJson2(conn)
		}


	} //end for {
}

func handleWorker(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)

	for {

		str, err := reader.ReadString('\n')
		//fmt.Println("ReadString: ", str)
		if err != nil {
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
		buf := []byte(str)
		go func(buf []byte, conn *net.TCPConn) {

			msg := protocol.GetRootAsData(buf, 0)
			//  do some thing
			cmd := string(msg.Cmd())
			data := string(msg.Data())
			req_sid := string(msg.Sid())
			req_id := int(msg.ReqId())
			golog.Info("handleWorker  ", cmd, data, req_sid, req_id)
			//fmt.Println("cmd: ", cmd)
			conn.Write(append(buf, '\n'))
		}(buf, conn)

	}

}

func handleWorkerJson2(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)

	for {

		str, err := reader.ReadString('\n')
		// fmt.Println("handleWorkerJson ReadString: ", str)
		if err != nil {
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
		buf := []byte(str)
		go func(buf []byte, conn *net.TCPConn) {

			/*
			msg_json, errjson := jason.NewObjectFromBytes( buf )
			if errjson != nil {
				return
			}
			cmd,  _ := msg_json.GetString("cmd")
			token,  _ := msg_json.GetString("token")
			golog.Error( "handleWorkerJson:", cmd, token )
			*/
			conn.Write(append(buf, '\n'))
		}(buf, conn)

	}

}

func handleWorkerJson(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据

	d := json.NewDecoder(conn)
	for {

		var msg interface{}

		err := d.Decode(&msg)
		if  err != nil {

			conn.Close()
			fmt.Println( "d.Decode(&msg) ", err.Error()  )
			break
		}
		buf,err_encode := json.Marshal( msg )
		if err_encode!=nil {
			fmt.Println( "json.Marshal error:",err_encode.Error() )
			conn.Close()
			break
		}
		msg_json, errjson := jason.NewObjectFromBytes( buf )
		if errjson != nil {
			continue
		}
		cmd,  _ := msg_json.GetString("cmd")
		token,  _ := msg_json.GetString("token")
		golog.Info( "handleWorkerJson:", cmd, token )

		go func(buf []byte, conn *net.TCPConn) {
			conn.Write(append(buf, '\n'))
		}(buf, conn)


	}

}

// 初始化worker服务
func InitWorkerServer() {

	for _, data := range global.Config.WorkerServer.Servers {

		host, _ := data[0].(string)
		port_str, _ := data[1].(string)
		port, _ := strconv.Atoi(port_str)
		global.WorkerServers = append(global.WorkerServers, []string{host, port_str})
		go WorkerServer(host, port)
	}
	//fmt.Println("global.WorkerServers:", global.WorkerServers)
}
