package worker

import (

	"strconv"
	"bufio"
	"fmt"
	"net"
	"encoding/json"
	"time"
	"reflect"
	"strings"
	"morego/protocol"
	"morego/area"
	"morego/global"
	"morego/golog"
	"morego/worker/golang"
)


// 初始化worker服务
func InitWorkerServer() {

	for _, data := range global.Config.WorkerServer.Servers {

		host, _ := data[0].(string)
		port_str, _ := data[1].(string)
		worker_language, _ := data[2].(string)
		port, _ := strconv.Atoi(port_str)
		//fmt.Println("worker_language:", worker_language)
		if worker_language == "go" {
			go WorkerServer(host, port)
		}
	}
}

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
		//conn.SetDeadline(30*time.Second)
		defer conn.Close()
		//conn.SetNoDelay(false)
		golog.Info("RemoteAddr:", conn.RemoteAddr().String())

		go handleWorker(conn)


	} //end for {
}


func handleWorker(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)

	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if( err.Error()!="EOF"){
				fmt.Println("HandleWork connection error: ", err.Error())
			}

			conn.Write([]byte(protocol.WrapRespErrStr(err.Error())))
			conn.Close()
			break
		}
		if( strings.Replace( string(buf), "\n", "", -1)==""){
			continue
		}
		if string(buf)=="ping" {
			conn.Write([]byte( "pong\n") )
			conn.Close()
			break
		}
		//fmt.Println( "HandleWorkerStr str: ",str)
		go func(buf []byte, conn *net.TCPConn) {

			//msg_err,_type,cmd,req_sid,reqid,req_data := protocol.ParseReqData( str )
			protocolJson := new(protocol.Json)
			protocolJson.Init()
			req_obj,_ := protocolJson.GetReqObj( buf )
			Invoker( conn, req_obj )

		}(buf, conn)
	}
}


func Invoker( conn *net.TCPConn, req_obj protocol.ReqRoot  ) []byte {

	task_obj := new(golang.TaskType).Init( conn, req_obj )

	fmt.Println( task_obj )

	invoker_ret:=InvokeObjectMethod( task_obj,req_obj.Header.Cmd )

	// 判断是否需要响应数据
	if req_obj.Type==protocol.TypeReq  && !req_obj.Header.NoResp {
		protocolJson := new(protocol.Json)
		protocolJson.Init()
		protocolJson.WrapRespObj( req_obj, invoker_ret, 200,"" )
		buf,_ := json.Marshal( protocolJson.ProtocolObj.RespObj )
		buf = append(buf, '\n')
		conn.Write( buf )
	}
	if( global.SingleMode ){
		if  global.IsAuthCmd(req_obj.Header.Cmd)   {
			area.ConnRegister( conn,req_obj.Header.Sid )
		}
	}
	return invoker_ret

}


func InvokeObjectMethod(object interface{}, methodName string, args ...interface{}) []byte{

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	//fmt.Println( "methodName:",methodName )
	ret := reflect.ValueOf(object).MethodByName(methodName).Call(inputs)[0]
	return ret.Bytes()

}
