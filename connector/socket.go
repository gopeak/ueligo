package connector

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"../area"
	"../global"
	"../golog"
	"../protocol"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"
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
	//go stat_tick()
	listenAcceptTCP(listen)
}

/**
 *  处理客户端连接
 */
func listenAcceptTCP(listen *net.TCPListener) {

	for {
		conn, err := listen.AcceptTCP()
		defer conn.Close()
		if err != nil {
			golog.Error("AcceptTCP Exception::", err.Error())
			continue
		}
		atomic.AddInt32(&global.SumConnections, 1)
		conn.SetNoDelay(false)

		// 校验ip地址
		conn.SetKeepAlive(true)

		// 获取随机worker服务地址

		configAddr := global.GetRandWorkerAddr()
		tcpAddr, err := net.ResolveTCPAddr("tcp4", configAddr)
		if err != nil {
			fmt.Println("req_conn tcpAddr :", err.Error())
			return
		}
		//fmt.Println("tcpAddr: ", tcpAddr)

		req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer req_conn.Close()
		if err != nil {
			fmt.Println("req_conn net.DialTCP :", err.Error())
			return
		}

		go handleClient(conn, req_conn, area.CreateSid())
		go handleWorkerResponse(conn, req_conn)
		//go handleClientMsgSingle( conn ,CreateSid() )

	} //end for {

}

func handleWorkerResponse(conn *net.TCPConn, req_conn *net.TCPConn) {

	reader := bufio.NewReader(req_conn)
	for {
		buf, err := reader.ReadBytes('\n')
		//fmt.Println("worker_task response:", msg)
		if err != nil {
			fmt.Println("handleWorkerResponse ", "error: ", err.Error())
			req_conn.Close()
			break
		}
		if strings.Replace(string(buf), "\n", "", -1) == "" {
			continue
		}
		WorkerResponseProcess(nil, conn, buf)
		conn.Write(buf)

	}
}

func handleClientMsgSingle(conn *net.TCPConn, sid string) {

	//声明一个管道用于接收解包的数据
	qps := 0 // make(chan int64, 0)
	reader := bufio.NewReader(conn)

	for {
		if !global.Config.Enable {
			//conn.Write([]byte(fmt.Sprintf("%s\r\n", global.DISBALE_RESPONSE)))
			conn.Close()
			break
		}
		_,header, data, err := protocol.DecodePacket( reader )
		if err != nil {
			conn.Close()
			break
		}
		qps++
		if qps%100 == 0 {
			fmt.Println("qps: ", qps)
		}
		atomic.AddInt64(&global.Qps, 1)

		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		req_obj,err := protocolPacket.GetReqHeaderObj( header )
		buf,_ := protocolPacket.WrapResp( "GetUserSession", req_obj.Sid, req_obj.SeqId , 200, data )
		conn.Write( buf )

	}
}

func handleClient(conn *net.TCPConn, req_conn *net.TCPConn, sid string) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	last_sid := ""
	defer area.FreeConn(conn, last_sid)
	for {
		if !global.Config.Enable {
			area.FreeConn(conn, last_sid)
			break
		}
		protocolPacket := new(protocol.Pack)
		protocolPacket.Init()
		req_obj ,err := protocolPacket.GetReqObjByReader( reader )
		if err != nil {
			golog.Error("SocketHandle protocolPacket.GetReqObjByReader err : "  + err.Error())
			area.FreeConn(conn, last_sid)
			break
		}
		last_sid = req_obj.Header.Sid
		ret, ret_err := dispatchMsg( req_obj, conn, req_conn )
		if ret_err != nil {
			if ret < 0 {
				fmt.Println(ret_err.Error())
				continue
			}
			if ret == 0 {
				fmt.Println(ret_err.Error())
				break
			}
		}

	}
}

/**
 * 根据消息类型分发处理
 */
func dispatchMsg(req_obj *protocol.ReqRoot, conn *net.TCPConn, req_conn *net.TCPConn) (int, error) {

	var err error

	buf, _ := json.Marshal(req_obj)
	buf = append(buf, '\n')

	//  认证检查, @todo 通过sid和worker判断非认证接口不能提交到worker中
	if !global.IsAuthCmd(req_obj.Header.Cmd) && !area.CheckSid(req_obj.Header.Sid) {
		area.FreeConn(conn, req_obj.Header.Sid)
		err = errors.New("认证失败")
		return 0, err
	}
	// 提交给worker  @todo判断单机模式下不需要请求worker
	if req_conn != nil {
		go req_conn.Write(buf)
	}

	return 1, nil
}

func reqWorker(buf []byte, req_conn *net.TCPConn) {

	req_conn.Write(buf)
	return
	//fmt.Println("worker agent from ", worker_idf, " receive 3:", msg)
	/*
	msg := protocol.GetRootAsData(buf, 0)
	//  do some thing
	cmd := string(msg.Cmd())
	data := string(msg.Data())
	req_sid := string(msg.Sid())
	req_id := int64(msg.ReqId())
	golog.Info("HandleConn data:", cmd, data, req_sid, req_id)
	*/

}

func checkError(err error) {
	if err != nil {
		golog.Error(os.Stderr, "Connector error: %s", err.Error())
	}
}

func stat_tick() {

	timer := time.Tick(1000 * time.Millisecond)
	for _ = range timer {
		//ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }` , time.Now().Unix() );
		fmt.Println(time.Now().Unix(), " Connections: ", global.SumConnections, "  Qps: ", global.Qps)
	}
}

func user_tick(conn *net.TCPConn) {

	timer := time.Tick(5000 * time.Millisecond)
	for _ = range timer {
		ping := fmt.Sprintf(`{"cmd":"ping","ret":200,"time":%d }`, time.Now().Unix())
		go conn.Write([]byte(fmt.Sprintf("%s\r\n", ping)))
	}
}
