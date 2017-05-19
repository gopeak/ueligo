package worker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
	"morego/area"
	"morego/global"
	"morego/golog"
	"morego/protocol"
	"morego/worker/golang"
	"github.com/BurntSushi/toml"

)



type WorkerConfigType struct {

	Loglevel     string		`toml:"loglevel"`
	RpcType      string		`toml:"rpc_type"`
	SingleMode   bool	  	`toml:"single_mode"`
	Servers [][]string       	`toml:"servers"`
	ToHub []string  		`toml:"to_hub"`
	//Mysql MysqlConfigType 		`toml:"mysql"`

}




var 	WorkerConfig   WorkerConfigType

// 初始化worker服务
func InitWorkerServer() {

	var err error
	_, err = toml.DecodeFile("worker/worker.toml", &WorkerConfig )
	if  err != nil {
		fmt.Println("toml.DecodeFile error:", err.Error())
		return
	}
	for _, data := range WorkerConfig.Servers {

		if len( data )<=2 {
			fmt.Println("worker.toml servers length err:" ,data )
			continue
		}
		host := data[0]
		port_str := data[1]
		worker_language  := data[2]
		port, _ := strconv.Atoi(port_str)
		if worker_language == "go" {
			go WorkerServer(host, port)
		}
	}
	time.Sleep( 1*time.Second)
	golang.InitReqHubPool()
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
			if err.Error() != "EOF" {
				fmt.Println("HandleWork connection error: ", err.Error())
			}

			conn.Write([]byte(protocol.WrapRespErrStr(err.Error())))
			conn.Close()
			break
		}
		if strings.Replace(string(buf), "\n", "", -1) == "" {
			continue
		}
		if string(buf) == "ping" {
			conn.Write([]byte("pong\n"))
			conn.Close()
			break
		}
		//fmt.Println( "HandleWorkerStr str: ",str)
		go func(buf []byte, conn *net.TCPConn) {

			protocolJson := new(protocol.Json)
			protocolJson.Init()
			req_obj, _ := protocolJson.GetReqObj(buf)
			Invoker(conn, req_obj)

		}(buf, conn)
	}
}

func Invoker(conn *net.TCPConn, req_obj *protocol.ReqRoot) interface{} {

	task_obj := new(golang.TaskType).Init(conn, req_obj)

	invoker_ret := InvokeObjectMethod(task_obj, req_obj.Header.Cmd)
	//fmt.Println("invoker_ret", invoker_ret)
	// 判断是否需要响应数据
	if req_obj.Type == "req" && !req_obj.Header.NoResp {
		protocolJson := new(protocol.Json)
		protocolJson.Init()
		protocolJson.WrapRespObj(req_obj, invoker_ret, 200 )
		buf, _ := json.Marshal( protocolJson.ProtocolObj.RespObj )
		buf = append(buf, '\n')
		conn.Write(buf)
	}
	if global.SingleMode {
		if global.IsAuthCmd(req_obj.Header.Cmd) {
			area.ConnRegister(conn, req_obj.Header.Sid)
		}
	}
	return invoker_ret
}

func InvokeObjectMethod(object interface{}, methodName string, args ...interface{}) interface{} {

	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	//fmt.Println("methodName:", methodName)
	ret := reflect.ValueOf(object).MethodByName(methodName).Call(inputs)[0]

	switch vtype := ret.Interface().(type) {

	case nil:
		return nil
	case bool:
		return ret.Interface().(bool)

	case float32:
		return ret.Interface().(float32)
	case float64:
		return ret.Interface().(float32)
	case int:
		return ret.Interface().(int)
	case uint8:
		return ret.Interface().(uint8)
	case uint16:
		return ret.Interface().(uint16)
	case uint32:
		return ret.Interface().(uint32)
	case uint64:
		return ret.Interface().(uint64)
	case int8:
		return ret.Interface().(int8)
	case int16:
		return ret.Interface().(int16)
	case int32:
		return ret.Interface().(int32)
	case int64:
		return ret.Interface().(int64)
	case []byte:
		return  ret.Interface().([]byte)
	case string:
		return  ret.Interface().(string)
	case []string:
		return ret.Interface().([]string)
	case map[string]string:
		return ret.Interface().(map[string]string)
	case map[string]interface{}:
		return ret.Interface().(map[string]interface{})
	case golang.ReturnType:
		return ret.Interface().(golang.ReturnType)
	default:
		fmt.Println("vtype:", vtype)
		golog.Error( "返回的类型无法处理:",vtype)
	}
	return ""

}
