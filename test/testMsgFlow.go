package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"morego/protocol"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

//var Conns = make([]*net.TCPConn, 1000)

func createReqConns(num int64) []*net.TCPConn {

	var conns []*net.TCPConn
	conns = make([]*net.TCPConn, 0)

	for i := 0; i < int(num); i++ {
		service := os.Args[1]
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
			fmt.Println(err.Error())
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		defer conn.Close()
		conns = append(conns, conn)
		time.Sleep(100 * time.Millisecond)
	}
	return conns

} //

func getReqStrData(num int64) string {

	data := ""

	type_ := protocol.TypeReq
	cmd := "Auth"
	req_data := time.Now().Unix()
	srcData :=  []byte( strconv.FormatInt(time.Now().Unix(), 10) )
	sid := md5.Sum([]byte(srcData))
	req_id := num
	data = fmt.Sprintf("%d||%s||%x||%d||%d", type_, cmd, sid, req_id, req_data)

	return data

} //

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	//start := time.Now().Unix()
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage:%s  host:port connections send_times packet_type ", os.Args[0])
		os.Exit(1)
	}
	//go end_hook( start )
	packet_type := "str"
	if len(os.Args) > 4 {
		packet_type = string(os.Args[4])
	}
	fmt.Println(" packet_type : ", packet_type)

	times, _ := strconv.ParseInt(os.Args[3], 10, 32)
	conn_num, _ := strconv.ParseInt(os.Args[2], 10, 32)
	fmt.Println("Connections and  times:", conn_num, times)

	//conns := createReqConns(conn_num)
	//ch_success := make(chan int64, 0)

	var i int64

	for i = 0; i < conn_num; i++ {

		service := os.Args[1]
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
			fmt.Println(err.Error())
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		srcData :=  []byte( strconv.FormatInt(time.Now().Unix(), 10) )
		data := fmt.Sprintf("%d||%s||%x||%d||%d\n", protocol.TypeReq, "Auth",  md5.Sum([]byte(srcData)) , i, time.Now().Unix())
		//getReqStrData(conn_num)
		_, errw := conn.Write([]byte(data))
		if err != nil {
			fmt.Println(errw.Error())
		}
		go func(conn *net.TCPConn ,times int64, conn_num int64) {

			reader := bufio.NewReader(conn)
			var success int64
			success = 0
			for {
				str, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("HandleConn connection error: ", err.Error())
					break
				}
				success++
				fmt.Println("recv msg: ", str)
				msg_err,_type,cmd,req_sid,req_id,msg_data := protocol.ParseData(str)
				if err != nil {
					fmt.Println("msg error: ", msg_err.Error(),_type,cmd,req_sid,req_id,msg_data)
					continue
				}
				fmt.Printf( " cmd: %s\n", cmd )
				// 登录认证,然后获取用户信息
				if cmd=="Auth" {
					//fmt.Printf( " sid: %s\n", sid )
					//str=fmt.Sprintf( `{ "cmd":"socket.getSession","params":{"sid":"%s"} }` ,req_sid )
					data = fmt.Sprintf("%d||%s||%s||%d||%s\n", protocol.TypeReq, "GetUserSession", req_sid, req_id+1, req_sid)
				 	conn.Write([]byte( data ))
				}
				// 获取当前信息后 发送点对点信息
				if cmd=="GetUserSession"  {

					to_sid:=""
					data = fmt.Sprintf("%d||%s||%s||%d||%s\n", protocol.TypePush, "Push", req_sid, req_id+1, to_sid)
					conn.Write([]byte( data ))

					i++
					if( i>=times ){
						fmt.Println( " i : ", i )
						fmt.Println( " conn close! " ,i ,"\n" )
						break
					}
				}
				// 发送点对点发送消息后 发送广播





			}

			return

		}(conn,  times, conn_num)


	}
	select {

	}
	/*
	var qps int64
	var recv_times int64
	qps = 0
	recv_times = 0
	for {
		select {
		case r := <-ch_success:
			recv_times++
		//fmt.Println("recv_times:", recv_times)
			qps = qps + r
			if recv_times == conn_num-1 {
				fmt.Printf(".")
				end := time.Now().Unix()
				els_time := end - start
				//fmt.Println("time:", els_time, qps)
				fmt.Printf("\nels_time:%d recv_times:%d qps:%d", els_time, recv_times, qps)
				return
			}
		default:
			fmt.Printf(".")
			time.Sleep(100 * time.Millisecond)
		}
	}
	*/

}
