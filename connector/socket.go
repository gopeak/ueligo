package connector

import (
	"bufio"
	"fmt"
	"net"
	"sync/atomic"
	"strings"
	"errors"
	"strconv"
	"morego/area"
	"morego/global"
	"morego/golog"
	"morego/protocol"
	"morego/worker"
	"github.com/antonholmquist/jason"

	"os"
	"time"
	"sync"
)

var ConnMlock *sync.RWMutex
var ChannelMlock *sync.RWMutex
var SessionMlock *sync.RWMutex
var UserChannelsMlock *sync.RWMutex



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
		fmt.Println("tcpAddr: ", tcpAddr)

		req_conn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer req_conn.Close()
		if err != nil {
			fmt.Println("req_conn net.DialTCP :", err.Error())
			return
		}


		go handleClientMsg(conn, req_conn, area.CreateSid())
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
		_,_,cmd,req_sid,_,msg_data := protocol.ParseRplyData(string(buf))
		fmt.Println( "handleWorkerResponse",string(buf) )
		if cmd==global.AuthCcmd && msg_data=="ok" {

			area.ConnRegister( conn,req_sid)
		}
		conn.Write(buf)

	}
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
			conn.Write([]byte(protocol.WrapRespErrStr("request data length error-->"+str)))
			continue
		}
		cmd := "user.getSession" //msg_arr[1];
		req_sid := msg_arr[2]
		req_id, _ := strconv.Atoi(msg_arr[3])
		data := msg_arr[4]
		resp_str := protocol.WrapRespStr(cmd, req_sid, req_id, data)
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
		fmt.Println( "HandleConn str: ",str)
		if err != nil {
			area.FreeConn(conn, sid)
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

	msg_err,_type,cmd,req_sid,req_id,req_data := protocol.ParseRplyData(str)
	if msg_err!=nil {
		return -1, msg_err
	}
	buf := []byte(str)
	//buf = append( buf, '\n')

	//  认证检查
	if ( cmd!=global.AuthCcmd && !area.CheckSid(req_sid) ) {
		area.FreeConn(conn, req_sid)
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
		req_data =msg_arr[protocol.MSG_DATA_INDEX]
		req_data = strings.Replace(req_data, "\n", "", -1)
		data_json, json_err := jason.NewObjectFromBytes([]byte(req_data))
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
		req_data =msg_arr[protocol.MSG_DATA_INDEX]
		req_data = strings.Replace(req_data, "\n", "", -1)
		data_json, json_err := jason.NewObjectFromBytes([]byte(req_data))
		if ( json_err != nil ) {
			err = errors.New("broatcast data json format error")
			return -3, err
		}
		area_id, _ := data_json.GetString("area_id")
		to_data, _ := data_json.GetString("data")
		if( area_id=="global" ) {
			err = errors.New("broatcast global failed")
			return -4, err
		}else{
			area.Broatcast( from_sid, area_id,to_data )
		}

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




