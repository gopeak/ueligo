package connector

import (
	"bufio"
	//"encoding/json"
	"fmt"
	//"math/rand"
	"net"
	"morego/area"
	"morego/global"
	"morego/golog"
	"morego/lib/antonholmquist/jason"
	//flatbuffers "github.com/google/flatbuffers/go"
	"morego/protocol"
	"morego/worker"
	"sync/atomic"
	//"time"
	//"encoding/json"
	"strings"
	"errors"
	"strconv"

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
	go stat_tick()
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
		fmt.Println("tcpAddr: ", tcpAddr)

		req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer req_conn.Close()
		if err != nil {
			fmt.Println("req_conn net.DialTCP :", err.Error())
			return
		}


		go handleClientMsg(conn, req_conn, CreateSid())
		go handleWorkerResponse(conn, req_conn)
		//go handleClientMsgSingle( conn ,CreateSid() )


	} //end for {

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

func  ReadBytes(delim byte) (line []byte, err error) {
	// Use ReadSlice to look for array,
	// accumulating full buffers.
	var b *bufio.Reader

	var frag []byte
	var full [][]byte

	for {
		var e error
		frag, e = b.ReadSlice(delim)
		if e == nil { // got final fragment
			break
		}
		if e != bufio.ErrBufferFull { // unexpected error
			err = e
			break
		}

		// Make a copy of the buffer.
		buf := make([]byte, len(frag))
		copy(buf, frag)
		full = append(full, buf)
	}

	// Allocate new buffer to hold the full pieces and the fragment.
	n := 0
	for i := range full {
		n += len(full[i])
	}
	n += len(frag)

	// Copy full pieces and fragment in.
	buf := make([]byte, n)
	n = 0
	for i := range full {
		n += copy(buf[n:], full[i])
	}
	copy(buf[n:], frag)
	return buf, err
}

func handleClientMsgSingle(conn *net.TCPConn,   sid string) {

	//声明一个管道用于接收解包的数据
	qps := 0;// make(chan int64, 0)
	reader := bufio.NewReader(  conn  )

	for {
		if !global.Config.Enable {
			//conn.Write([]byte(fmt.Sprintf("%s\r\n", global.DISBALE_RESPONSE)))
			conn.Close()
			break
		}

		buf,err := reader.ReadBytes('\n')
		if err != nil {
			//fmt.Println("err ReadString:", err.Error())
			conn.Close()
			break
		}
		qps++
		if( qps%100==0){
			fmt.Println( "qps: ", qps )
		}
		atomic.AddInt64(&global.Qps, 1)
		str := string(buf)
		//fmt.Println( "HandleConn str: ",str)

		msg_arr := strings.Split(str, "||")
		if len(msg_arr) < 5 {
			conn.Write([]byte(worker.WrapRespErrStr("request data length error-->"+str)))
			continue
		}
		cmd := "user.getSession" //msg_arr[1];
		req_sid := msg_arr[2]
		req_id, _ := strconv.Atoi(msg_arr[3])
		data := msg_arr[4]
		resp_str := worker.WrapRespStr(cmd, req_sid, req_id, data)

		conn.Write([]byte(resp_str))

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
		//fmt.Println( "HandleConn str: ",str)
		if err != nil {
			FreeConn(conn, sid)
			//fmt.Println( "HandleConn connection error: ", err.Error())
			break
		}

		ret, ret_err := dispatchMsg(str, conn, req_conn)

		if ( ret_err != nil ) {
			if ( ret < 0 ) {
				fmt.Println(ret_err.Error(), str)
				continue
			}
			if ( ret == 0 ) {
				fmt.Println(ret_err.Error(), str)
				break
			}
		}

	}
}

/**
 * 根据消息类型分发处理
 */
func dispatchMsg(str string, conn *net.TCPConn, req_conn *net.TCPConn) (int, error) {

	var err error
	msg_arr := strings.Split(str, "||")
	if len(msg_arr) < 5 {
		conn.Close()
		err = errors.New("request data length error")
		return -1, err
	}
	_type,_ := strconv.Atoi(msg_arr[protocol.MSG_TYPE_INDEX])
	cmd := msg_arr[protocol.MSG_CMD_INDEX];
	req_sid := msg_arr[protocol.MSG_SID_INDEX]
	req_id :=msg_arr[protocol.MSG_REQID_INDEX]
	req_data := msg_arr[protocol.MSG_DATA_INDEX]
	buf := []byte(str)
	buf = append( buf, '\n')

	//  认证检查
	if ( cmd!=global.AuthCcmd && !CheckSid(req_sid) ) {
		FreeConn(conn, req_sid)
		err = errors.New("认证失败")
		return 0, err
	}
	// 请求
	if _type == protocol.TypeReq {
		// 如果是单机模式,则直接调用
		if( global.SingleMode ){
			go worker.Invoker( conn,cmd, req_sid,req_id,req_data )
		}else{
			go req_conn.Write(buf)
		}
	}
	if _type == protocol.TypePush {
		from_sid := msg_arr[protocol.MSG_SID_INDEX]
		data_json, json_err := jason.NewObjectFromBytes([]byte(msg_arr[protocol.MSG_DATA_INDEX]))
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
		from_sid := msg_arr[protocol.MSG_SID_INDEX]
		data_json, json_err := jason.NewObjectFromBytes([]byte(msg_arr[protocol.MSG_DATA_INDEX]))
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

/**
 * 认证
 */
func auth(token string, conn *net.TCPConn)   {

}



/**
 * 广播
 */
func broadcast(sid string, area_id string, data string, conn *net.TCPConn) {

}


