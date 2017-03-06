package connector

import (
	"bufio"
	//"encoding/json"
	"fmt"
	//"math/rand"
	"net"
	//"simple/area"
	"simple/global"
	"simple/golog"
	//"simple/lib/antonholmquist/jason"
	//flatbuffers "github.com/google/flatbuffers/go"
	"simple/protocol"
	//"simple/worker"
	"sync/atomic"
	//"time"
	//"encoding/json"
)

/**
 * 监听客户端连接
 */
func SocketConnector(ip string, port int) {

	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), port, ""})
	if err != nil {
		golog.Error("ListenTCP Exception:", err.Error())
		return
	}
	// 初始化
	golog.Debug("Game Connetor Server :", ip, port)

	listenAcceptTCP(listen)
}

/**
 *  处理客户端连接
 */
func listenAcceptTCP(listen *net.TCPListener) {

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			golog.Error("AcceptTCP Exception::", err.Error())
			continue
		}

		//defer conn.Close()
		atomic.AddInt32(&global.SumConnections, 1)
		conn.SetNoDelay(false)

		// 校验ip地址
		conn.SetKeepAlive(true)
		golog.Info("RemoteAddr:", conn.RemoteAddr().String())

		//remoteAddr :=conn.RemoteAddr()
		// 获取随机worker服务地址
		ip_port := global.GetRandWorkerAddr()

		//fmt.Println("ip_port:", ip_port)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", ip_port)
		if err != nil {
			fmt.Println("req_conn tcpAddr :", err.Error())
			return
		}

		req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer req_conn.Close()
		if err != nil {
			fmt.Println("req_conn net.DialTCP :", err.Error())
			return
		}

		//fmt.Println("RemoteAddr:", conn.RemoteAddr().String(), "sid:", sid, " worker_idf:", "")

		if( global.PackSplitType=="bufferio"){
			go handleClientMsg(conn, req_conn, CreateSid())
		}
		if( global.PackSplitType=="json"){
			go handleConnJson(conn, req_conn,CreateSid())
		}
		go handleWorkerResponse(conn, req_conn)
		//go handleConn(conn, sid, "")

	} //end for {

}

/**
 * 客户端通过json方式封包数据
 */
func handleConnJson(conn *net.TCPConn, req_conn *net.TCPConn, sid string ) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	for {
		if !global.Config.Enable {
			conn.Write([]byte(fmt.Sprintf("%s\r\n", global.DISBALE_RESPONSE)))
			conn.Close()
			break
		}

		buf, err := reader.ReadBytes('\n')
		//fmt.Println(  "handleConn ReadString: ", string(buf) )
		if err != nil {
			FreeConn(conn, sid)
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
		/*
		worker_json, errjson := jason.NewObjectFromBytes(buf)
		checkError(errjson)
		// do some thing
		cmd, _ := worker_json.GetString("cmd")
		// fmt.Printf(" worker_task logic cmd: %s\n", cmd)
		json := fmt.Sprintf(`{"cmd":"%s","data":"%s"}`, cmd, sid)
		golog.Info("handleConnJson json ", json )
		*/
		go reqWorker(buf, req_conn)

	}



}

func handleWorkerResponse(conn *net.TCPConn, req_conn *net.TCPConn) {

	reader := bufio.NewReader(req_conn)
	for {
		msg, err := reader.ReadBytes('\n')
		//fmt.Println("worker_task response:", msg)
		if err != nil {
			fmt.Println("handleWorkerResponse ", "error: ", err.Error())
			req_conn.Close()
			break
		}
		if msg == nil {
			continue
		}

		if string(msg) == "\n" {
			continue
		}
		conn.Write(msg)

	}
}

func handleClientMsg(conn *net.TCPConn, req_conn *net.TCPConn, sid string) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	for {
		if !global.Config.Enable {
			//conn.Write([]byte(fmt.Sprintf("%s\r\n", global.DISBALE_RESPONSE)))
			conn.Close()
			break
		}

		str, err := reader.ReadString('\n')
		//fmt.Println("handleConn ReadString: ", str)
		if err != nil {
			FreeConn(conn, sid)
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
		buf := []byte(str)
		go reqWorker(buf, req_conn)

	}

}

func reqWorker(buf []byte, req_conn *net.TCPConn) {

	req_conn.Write(buf)
	return
	//fmt.Println("worker agent from ", worker_idf, " receive 3:", msg)
	msg := protocol.GetRootAsData(buf, 0)
	//  do some thing
	cmd := string(msg.Cmd())
	data := string(msg.Data())
	req_sid := string(msg.Sid())
	req_id := int64(msg.ReqId())
	golog.Info("HandleConn data:", cmd, data, req_sid, req_id)

}
