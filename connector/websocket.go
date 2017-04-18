package connector

import (
	"bufio"
	"fmt"
	"math/rand"
	"morego/area"
	"morego/global"
	"morego/protocol"
	"net"
	"net/http"
	"sync/atomic"
	"time"
	"morego/golog"
	"github.com/antonholmquist/jason"
	"github.com/gorilla/websocket"
	"errors"
	"flag"
	"log"
	"strconv"
	"strings"
	"os"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
} // use default options

func WebsocketConnector(ip string, port int) {

	golog.Info("Websocket Connetor bind :", ip, port)

	var addr = flag.String("addr", fmt.Sprintf(":%d", port), "http service address")
	http.HandleFunc("/ws", WebsocketHandle)
	wd, _ := os.Getwd()
	http_dir := fmt.Sprintf("%s/web/wwwroot", wd)
	fmt.Println("Http_dir:", http_dir)
	http.Handle("/", http.FileServer(http.Dir(http_dir)))


	log.Fatal(http.ListenAndServe(*addr, nil))

}

/**
 *  处理客户端连接
 */
func WebsocketHandle(writer http.ResponseWriter, request *http.Request) {
	wsconn, err := upgrader.Upgrade(writer, request, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	var max_conns int32
	fmt.Println(" websocke client connect:", wsconn.RemoteAddr())
	user_sid := ""
	//remoteAddr :=conn.RemoteAddr()
	atomic.AddInt32(&global.SumConnections, 1)

	max_conns = int32(global.Config.Connector.MaxConections)
	if max_conns > 0 && global.SumConnections > max_conns {

		wsconn.WriteMessage(websocket.TextMessage, []byte(global.ERROR_MAX_CONNECTIONS))
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sid := fmt.Sprintf("%d%d", r.Intn(99999), rand.Intn(999999))

	configAddr := global.GetRandWorkerAddr()
	fmt.Println("ip_port:", configAddr)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", configAddr)
	checkError(err)
	req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
	//defer req_conn.Close()
	checkError(err)
	go wsHandleWorkerResponse(wsconn, req_conn)

	// 监听客户端发送的数据

	for {
		var str string
		_, buf, err := wsconn.ReadMessage()
		if err != nil {
			fmt.Println(" websocket.Message.Receive error:", user_sid, "  -->", err.Error())
			wsconn.Close()
			break
		}
		str = string(buf)
		//fmt.Println("Client Request: " + str)

		go func(sid string, str string, wsconn *websocket.Conn, req_conn *net.TCPConn) {

			ret, ret_err := wsDspatchMsg(str, wsconn, req_conn)
			if ret_err != nil {
				if ret < 0 {
					fmt.Println(ret_err.Error(), str)
					return
				}
				if ret == 0 {
					fmt.Println(ret_err.Error(), str)
					return
				}
			}

		}(sid, str, wsconn, req_conn)

	}

}

func wsHandleWorkerResponse(wsconn *websocket.Conn, req_conn *net.TCPConn) {

	reader := bufio.NewReader(req_conn)
	var l *sync.RWMutex
	l = new(sync.RWMutex)

	for {
		buf, err := reader.ReadBytes('\n')
		//fmt.Println("worker_task response:", msg)
		if err != nil {
			fmt.Println("handleWorkerResponse ", "error: ", err.Error())
			req_conn.Close()
			break
		}

		if( strings.Replace(string(buf), "\n", "", -1)==""){
			continue
		}
		_,_,cmd,_,_,msg_data := protocol.ParseRplyData(string(buf))

		if cmd==global.AuthCcmd  {
			fmt.Println( "AuthCcmd:",string(buf) )
			data_json ,err_json:= jason.NewObjectFromBytes( []byte(msg_data ) )
			if( err_json!=nil ) {
				golog.Error("auth  json err:",err_json.Error())
				continue
			}
			auth_ret,_ := data_json.GetString("ret")
			if( auth_ret=="ok"){
				sid,_ := data_json.GetString("id")
				area.WsConnRegister( wsconn,sid)
			}
		}

		l.Lock()
		wsconn.WriteMessage(websocket.TextMessage, buf)
		l.Unlock()
	}
}

/**
 * 根据消息类型分发处理
 */
func wsDspatchMsg(str string, wsconn *websocket.Conn, req_conn *net.TCPConn) (int, error) {

	var err error
	msg_arr := strings.Split(str, "||")
	if len(msg_arr) < 5 {
		wsconn.Close()
		err = errors.New("request data length error")
		return -1, err
	}
	_type, _ := strconv.Atoi(msg_arr[0])
	cmd := msg_arr[1]
	req_sid := msg_arr[2]
	buf := []byte(str)
	buf = append(buf, '\n')

	//  认证检查
	if cmd != "user.getUser" && !area.CheckSid(req_sid) {
		area.FreeWsConn(wsconn, req_sid)
		err = errors.New("认证失败")
		return 0, err
	}
	// 请求
	if _type == protocol.TypeReq {
		go req_conn.Write(buf)
	}
	if _type == protocol.TypePush {
		from_sid := msg_arr[2]
		data_json, json_err := jason.NewObjectFromBytes([]byte(msg_arr[4]))
		if json_err != nil {
			err = errors.New("push data json format error")
			return -2, err
		}
		to_sid, _ := data_json.GetString("sid")
		to_data, _ := data_json.GetString("data")
		area.Push(to_sid, from_sid, to_data)
	}
	if _type == protocol.TypeBroadcast {
		//from_sid := msg_arr[2]
		from_sid := msg_arr[2]
		data_json, json_err := jason.NewObjectFromBytes([]byte(msg_arr[4]))
		if json_err != nil {
			err = errors.New("broatcast data json format error")
			return -3, err
		}
		area_id, _ := data_json.GetString("area_id")
		to_data, _ := data_json.GetString("data")
		area.Broatcast(from_sid, area_id, to_data)
	}

	return 1, nil
}
