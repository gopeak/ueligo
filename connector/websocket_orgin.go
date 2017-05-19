package connector

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"morego/area"
	"morego/global"
	"morego/golog"
	"morego/lib/websocket"
	"morego/protocol"
	"morego/worker/golang"
	"net"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"encoding/json"
)

func WebsocketConnector(ip string, port int) {

	golog.Info("Websocket Connetor bind :", ip, port)

	var addr = flag.String("addr", fmt.Sprintf(":%d", port), "http service address")

	http.Handle("/ws", websocket.Handler(WebsocketHandleClient))

	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/web/wwwroot", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))
	// 初始化群组
	golang.InitGlobalGroup()
	// http请求处理
	golang.InitHandler()

	log.Fatal(http.ListenAndServe(*addr, nil))

}

/**
 *  处理客户端连接
 */
func WebsocketHandleClient(wsconn *websocket.Conn) {

	var max_conns int32
	fmt.Println(" websocke client connect:", wsconn.RemoteAddr())
	//remoteAddr :=conn.RemoteAddr()
	atomic.AddInt32(&global.SumConnections, 1)

	max_conns = int32(global.Config.Connector.MaxConections)
	if max_conns > 0 && global.SumConnections > max_conns {
		wsconn.Write([]byte(global.ERROR_MAX_CONNECTIONS))
		return

	}
	configAddr := global.GetRandWorkerAddr()
	fmt.Println("ip_port:", configAddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", configAddr)
	checkError(err)
	req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//defer req_conn.Close()
	checkError(err)
	go wsHandleWorkerResponse(wsconn, req_conn)
	last_sid := ""
	// 监听客户端发送的数据

	for {
		var buf []byte
		if err = websocket.Message.Receive(wsconn, &buf); err != nil {
			fmt.Println(" websocket.Message.Receive error:", last_sid, "  -->", err.Error())
			area.FreeWsConn(wsconn, last_sid)
			break
		}
		protocolJson := new(protocol.Json)
		protocolJson.Init()
		req_obj, err := protocolJson.GetReqObj(buf)
		if err != nil {
			golog.Error("1.WebsocketHandle protocolJson.GetReqObj err : " + string(buf) + err.Error())
			continue
		}
		last_sid = req_obj.Header.Sid
		fmt.Println("req_obj.Header.Cmd: " +  req_obj.Header.Cmd)
		//_,_,_,last_sid,_,_ = protocol.ParseReqData( str )

		go func(req_obj *protocol.ReqRoot, wsconn *websocket.Conn, req_conn *net.TCPConn) {

			ret, ret_err := wsDspatchMsg(req_obj, wsconn, req_conn)
			if ret_err != nil {
				if ret < 0 {
					fmt.Println(ret_err.Error())
					return
				}
				if ret == 0 {
					fmt.Println(ret_err.Error())
					return
				}
			}

		}(req_obj, wsconn, req_conn)

	}

}

func wsHandleWorkerResponse(wsconn *websocket.Conn, req_conn *net.TCPConn) {

	reader := bufio.NewReader(req_conn)

	for {
		str, err := reader.ReadString('\n')
		fmt.Println("worker_task response:", str )
		if err != nil {
			fmt.Println("handleWorkerResponse ", "error: ", err.Error())
			req_conn.Close()
			break
		}
		str = strings.Replace( str, "\n", "", -1)
		if str== "" {
			continue
		}
		WorkerResponseProcess(wsconn, nil, []byte(str) )
		//var l *sync.RWMutex
		//l = new(sync.RWMutex)
		//l.Lock()
		go wsconn.Write([]byte(str))
		//l.Unlock()
	}
}

func WorkerResponseProcess(wsconn *websocket.Conn, conn *net.TCPConn, buf []byte) {

	protocolJson := new(protocol.Json)
	protocolJson.Init()
	resp_obj, _ := protocolJson.GetRespObj(buf)
	//fmt.Println("handleWorkerResponse resp_obj: ", resp_obj )

	if global.IsAuthCmd(resp_obj.Header.Cmd) {
		data := resp_obj.Data.(map[string]interface{})
		auth_ret := data["ret"].(string)
		_sid := data["sid"].(string)
		if auth_ret == "ok" {
			if conn != nil {
				area.ConnRegister(conn, _sid)
			}
			if wsconn != nil {
				area.WsConnRegister(wsconn, _sid)
			}
			fmt.Println("handleWorkerResponse AuthCmd sid: ", resp_obj.Header.Cmd, _sid )
		}
	}
}

/**
 * 根据消息类型分发处理
 */
func wsDspatchMsg(req_obj *protocol.ReqRoot, wsconn *websocket.Conn, req_conn *net.TCPConn) (int, error) {

	var err error

	buf, _ := json.Marshal(req_obj)
	buf = append(buf, '\n')

	//  认证检查, @todo 通过sid和worker判断非认证接口不能提交到worker中
	if !global.IsAuthCmd(req_obj.Header.Cmd) && !area.CheckSid(req_obj.Header.Sid) {
		area.FreeWsConn(wsconn, req_obj.Header.Sid)
		err = errors.New("认证失败")
		return 0, err
	}
	// 提交给worker  @todo判断单机模式下不需要请求worker
	if req_conn != nil {
		go req_conn.Write(buf)
	}

	return 1, nil
}
