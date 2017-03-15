package connector

import (
	"bufio"
	"fmt"
	"morego/area"
	"morego/global"
	"math/rand"
	"net"
	"net/http"
	"sync/atomic"
	"time"
	"morego/protocol"
	//"encoding/json"
	"morego/lib/antonholmquist/jason"
	"morego/lib/websocket"
	"morego/golog"

	//"strings"
	"errors"
	//sync"
	"strings"
	"strconv"
)

func WebsocketConnector(ip string, port int) {

	http.Handle("/", websocket.Handler(WebsocketHandle))
	golog.Info("Websocket Connetor bind :", ip, port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		golog.Error("ListenAndServe:", err)
	}

}

/**
 *  处理客户端连接
 */
func WebsocketHandle( wsconn *websocket.Conn ) {

	var err error
	var max_conns int32
	fmt.Println(" websocke client connect:" ,wsconn.RemoteAddr() )
	user_sid:=""
	//remoteAddr :=conn.RemoteAddr()


	atomic.AddInt32(&global.SumConnections, 1)

	max_conns = int32(global.Config.Connector.MaxConections)
	if max_conns > 0 && global.SumConnections > max_conns {
		wsconn.Write([]byte(global.ERROR_MAX_CONNECTIONS + "\n"))
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
	go wsHandleWorkerResponse( wsconn, req_conn )

	// 监听客户端发送的数据

	for {
		var str string

		if err = websocket.Message.Receive(wsconn, &str); err != nil {

			fmt.Println(" websocket.Message.Receive error:" ,user_sid,"  -->", err.Error() )

			wsconn.Close()
			break

		}

		fmt.Println("Client Request: " + str)
		//websocket.Message.Send(ws,str)
		go func(sid string, str string, wsconn *websocket.Conn, req_conn *net.TCPConn) {

			ret, ret_err := wsDspatchMsg(str, wsconn, req_conn)
			if ( ret_err != nil ) {
				if ( ret < 0 ) {
					fmt.Println(ret_err.Error(), str)
					return
				}
				if ( ret == 0 ) {
					fmt.Println(ret_err.Error(), str)
					return
				}
			}

		}(sid, str,wsconn, req_conn)


	}



}




func wsHandleWorkerResponse(conn *websocket.Conn, req_conn *net.TCPConn) {

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

func wsHandleClientMsg(conn *websocket.Conn, req_conn *net.TCPConn, sid string ) {

	// 发包频率判断
	range_count := 1
	limit_date := global.Config.Connector.MaxPacketRate
	var now int64
	var start_time int64
	var range_times int64
	start_time = time.Now().Unix()
	range_times = int64(global.Config.Connector.MaxPacketRateUnit)

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	for {
		if !global.Config.Enable {
			conn.Write([]byte(fmt.Sprintf("%s\r\n", global.DISBALE_RESPONSE)))
			conn.Close()
			break
		}

		// 区间范围内的计数
		if limit_date > 0 {
			now = time.Now().Unix()
			if (now - start_time) <= range_times {
				range_count++
			} else {
				start_time = now
				range_count = 1
			}
			// 判断发包频率是否超过限制
			if range_count > limit_date {
				conn.Write([]byte(global.ERROR_PACKET_RATES + "\n"))
				conn.Close()
				break
			}
		}

		str, err := reader.ReadString('\n')
		//fmt.Println(  "handleConn ReadString: ", string(msg) )
		if err != nil {
			FreeWsConn(conn, sid)
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
		if str == "" {
			continue
		}
		go func(sid string, str string, conn *websocket.Conn, req_conn *net.TCPConn) {

			ret, ret_err := wsDspatchMsg(str, conn, req_conn)
			if ( ret_err != nil ) {
				if ( ret < 0 ) {
					fmt.Println(ret_err.Error(), str)
					return
				}
				if ( ret == 0 ) {
					fmt.Println(ret_err.Error(), str)
					return
				}
			}

		}(sid, str,conn, req_conn)

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
	_type,_ := strconv.Atoi(msg_arr[0])
	cmd := msg_arr[1];
	req_sid := msg_arr[2]
	buf := []byte(str)
	buf = append( buf, '\n')

	//  认证检查
	if ( cmd!="user.getUser" && !CheckSid(req_sid) ) {
		FreeWsConn(wsconn, req_sid)
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
		if ( json_err != nil ) {
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
		if ( json_err != nil ) {
			err = errors.New("broatcast data json format error")
			return -3, err
		}
		area_id, _ := data_json.GetString("area_id")
		to_data, _ := data_json.GetString("data")
		area.Broatcast( from_sid, area_id,to_data )
	}


	return 1, nil
}

