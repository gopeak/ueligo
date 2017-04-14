package worker

import (

	"strconv"
	"bufio"
	"fmt"
	"net"
	"encoding/json"
	"time"
	"reflect"
	"morego/protocol"
	"morego/area"
	"morego/global"
	"morego/golog"
	"github.com/antonholmquist/jason"
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
	//fmt.Println("global.WorkerServers:", global.WorkerServers)
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

		if global.PackSplitType == "bufferio" {
			go handleWorkerStrSplit(conn)
		}
		if global.PackSplitType == "json" {
			go handleWorkerJson(conn)
		}

	} //end for {
}


func handleWorkerStrSplit(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			if( err.Error()!="EOF"){
				fmt.Println("HandleWork connection error: ", err.Error())
			}

			conn.Write([]byte(protocol.WrapRespErrStr(err.Error())))
			conn.Close()
			break
		}
		//fmt.Println( "HandleWorkerStr str: ",str)
		go func(str string, conn *net.TCPConn) {

			msg_err,_type,cmd,req_sid,reqid,req_data := protocol.ParseReqData( str )
			if( msg_err!=nil ){
				golog.Error(msg_err.Error(),_type,cmd,req_sid,reqid,req_data  )
				return
			}
			if( _type==protocol.TypePing ) {
				conn.Write([]byte(protocol.WrapRespStr("Ping","",0,"")))
				conn.Close()
			}else{


				Invoker( conn,cmd,req_sid,reqid,req_data)
			}



		}(str, conn)
	}
}

func handleWorkerFlatBuffer(conn *net.TCPConn) {

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

func handleWorkerJson(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	d := json.NewDecoder(conn)
	for {

		var msg interface{}

		err := d.Decode(&msg)
		if err != nil {

			conn.Close()
			fmt.Println("d.Decode(&msg) ", err.Error())
			break
		}
		buf, err_encode := json.Marshal(msg)
		if err_encode != nil {
			fmt.Println("json.Marshal error:", err_encode.Error())
			conn.Close()
			break
		}
		msg_json, errjson := jason.NewObjectFromBytes(buf)
		if errjson != nil {
			continue
		}
		cmd, _ := msg_json.GetString("cmd")
		token, _ := msg_json.GetString("token")
		golog.Info("handleWorkerJson:", cmd, token)

		go func(buf []byte, conn *net.TCPConn) {
			conn.Write(append(buf, '\n'))
		}(buf, conn)

	}

}

func Invoker( conn *net.TCPConn,cmd string, req_sid string ,req_id int,req_data string ) string {

	task_obj := new(TaskType).Init( conn, cmd, req_sid,req_id,req_data )
	data:=InvokeObjectMethod( task_obj,cmd )
	//fmt.Println( "Invoker:", data )
	resp_str := protocol.WrapRespStr(cmd, req_sid, req_id, data)
	//fmt.Println( "resp_str:", resp_str )
	conn.Write(  []byte(resp_str) )
	if( global.SingleMode ){
		if cmd==global.AuthCcmd && data=="ok" {
			area.ConnRegister( conn,req_sid)
		}
	}
	return data

}


func InvokeObjectMethod(object interface{}, methodName string, args ...interface{}) string{

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	fmt.Println( "methodName:",methodName )
	ret := reflect.ValueOf(object).MethodByName(methodName).Call(inputs)[0]
	//fmt.Println( "ret:" ,ret)
	data:=""
	value := reflect.ValueOf(&ret)
	value = reflect.Indirect(value)
	switch value.Kind(){      //多选语句switch
	case reflect.String:
		data = fmt.Sprintf("%s",ret)
	case reflect.Int ,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64:
		data = fmt.Sprintf("%d",ret)
	case reflect.Float32,reflect.Float64:
		data = fmt.Sprintf("%f",ret)
	default:
		data = fmt.Sprintf("%s",ret)
	}

	return data
}
