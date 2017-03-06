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
	//"morego/protocol"
	//"encoding/json"
	"morego/lib/websocket"
	"morego/golog"

	//"strings"
	//"io"
	//sync"
)

func WebsocketConnector(ip string, port int) {

	http.Handle("/", websocket.Handler(WebsocketHandler))
	golog.Info("Websocket Connetor bind :", ip, port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		golog.Error("ListenAndServe:", err)
	}

}

/**
 *  处理客户端连接
 */
func WebsocketHandler(ws *websocket.Conn) {

	var err error
	var max_conns int32

	atomic.AddInt32(&global.SumConnections, 1)
	user_sid := ""

	for {
		var str string
		atomic.AddInt32(&global.SumConnections, 1)

		max_conns = int32(global.Config.Connector.MaxConections)
		if max_conns > 0 && global.SumConnections > max_conns {
			ws.Write([]byte(global.ERROR_MAX_CONNECTIONS + "\n"))
			continue
		}
		if err = websocket.Message.Receive(ws, &str); err != nil {

			golog.Error(" websocket.Message.Receive error:", user_sid, "  -->", err.Error())
			if err.Error() == "EOF" {

				FreeWsConn(ws, user_sid)
			}

		}

		//remoteAddr :=conn.RemoteAddr()
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		sid := fmt.Sprintf("%d%d", r.Intn(99999), rand.Intn(999999))


		configAddr := global.GetRandWorkerAddr()
		fmt.Println("ip_port:", configAddr)
		tcpAddr, err := net.ResolveTCPAddr("tcp4", configAddr)
		checkError(err)
		req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
		//defer req_conn.Close()
		checkError(err)
		worker_idf := ""

		//fmt.Println("RemoteAddr:", conn.RemoteAddr().String() , "sid:", sid ," worker_idf:" ,worker_idf )

		// 接收worker返回的数据
		go WsReqWorkerAgentWithBufferio(ws, req_conn, sid, worker_idf)

		go WsHandleConnWithBufferio(ws, req_conn, sid, worker_idf)

		golog.Debug("Client Request: " + str)

	}
	ws.Write([]byte{'E', 'O', 'F'})
	ws.Close()

}

func WsReqWorkerAgentWithBufferio(wsconn *websocket.Conn, req_conn *net.TCPConn, sid string, worker_idf string) {


	area.WsConnRegister(wsconn, sid)
	//req_ready := fmt.Sprintf( `{"cmd":"req.connect", "client_idf":"%s" }`, sid    )
	req_ready := fmt.Sprintf(`%s||%s||%s||`, global.DATA_REQ_CONNECT, sid, worker_idf)

	req_conn.Write([]byte(req_ready + "\n"))

}

func WsHandleConnWithBufferio(conn *websocket.Conn, req_conn *net.TCPConn, sid string, worker_idf string) {

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

		msg, err := reader.ReadString('\n')
		//fmt.Println(  "handleConn ReadString: ", string(msg) )
		if err != nil {
			FreeWsConn(conn, sid)
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}
		if msg == "" {
			continue
		}
		go func(sid string, msg string, req_conn *net.TCPConn) {

			// fmt.Println(conn.RemoteAddr().String(), "receive str:", string(msg) )
			worker_idf :=""
			worker_data := fmt.Sprintf(`%s||%s||%s||%s`, global.DATA_REQ_MSG, sid, worker_idf, msg)
			//fmt.Println("req push worker_data 2:", worker_data)
			req_conn.Write([]byte(worker_data + "\n"))

		}(sid, msg, req_conn)

	}

}
