package main

import (
	"bufio"
	//"crypto/md5"
	"fmt"
	"morego/protocol"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
	"strings"
)

var Conns   []*net.TCPConn
var Sids []string


func createReqConns(num int64)  {


	Conns = make([]*net.TCPConn, 0)
	Sids = make([]string, 0)
	for i := 0; i < int(num); i++ {
		service := os.Args[1]
		tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
		if err != nil {
			fmt.Println(err.Error())
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		//defer conn.Close()
		Conns = append(Conns, conn)
		time.Sleep(10 * time.Millisecond)
		//srcData :=  []byte( strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(i)  )
		//str:=md5.Sum([]byte(srcData))
		data :=  protocol.WrapReqStr("Auth", strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(i) , i, strconv.FormatInt(int64(time.Now().Unix()), 10) ) // fmt.Sprintf("%d||%s||%x||%d||%d\n", protocol.TypeReq, "Auth",  md5.Sum([]byte(srcData)) , i, time.Now().Unix())
		conn.Write([]byte(data))
		r := bufio.NewReader(conn)
		for {
			str, _ := r.ReadString('\n')
			fmt.Println( "Auth:", str )
			msg_err,_type,cmd,req_sid,req_id,msg_data := protocol.ParseRplyData(str)
			if msg_err != nil {
				fmt.Println("msg auth error: ", msg_err.Error(),_type,cmd,req_sid,req_id,msg_data)
				continue
			}
			Sids  = append(  Sids, req_sid )
			break
		}
	}


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

	createReqConns(conn_num)
	//ch_success := make(chan int64, 0)

	var i int64

	for i = 0; i < conn_num; i++ {

		conn := Conns[i]

		go func(conn *net.TCPConn ,times int64, conn_num int64 ,i int) {
			//fmt.Println( conn )
			reader := bufio.NewReader(conn)
			var success int64
			success = 0
			req_sid := Sids[i]
			data :=  protocol.WrapReqStr("GetUserSession",req_sid,0,req_sid )
			conn.Write([]byte( data ))

			for {
				str, err := reader.ReadString('\n')
				if err != nil {
					fmt.Println("HandleConn connection error: ", err.Error())
					break
				}
				success++
				// fmt.Println("recv msg: ", str)
				msg_arr := strings.Split(str, "||")

				if len(msg_arr) <3 {
					fmt.Println(" recv msg length error: ", msg_arr )
					continue
				}
				_type ,_ :=strconv.Atoi(msg_arr[protocol.MSG_TYPE_INDEX])
				if( _type==protocol.TypeReply ){
					msg_err,_,cmd,req_sid,req_id,msg_data := protocol.ParseRplyData(str)
					if err != nil {
						fmt.Println("msg error: ", msg_err.Error(),msg_arr, msg_data)
						continue
					}
					// 登录认证,然后获取用户信息

					// 获取当前信息后 发送点对点信息
					if cmd=="GetUserSession"  {

						to_sid_index := i-1
						if to_sid_index<0 {
							to_sid_index = 0
						}
						to_sid:= Sids[to_sid_index]
						push_data := fmt.Sprintf(`{"sid":"%s","data":"%s"}`,to_sid,"md55555555555")
						data =  protocol.WrapPushStr("Push",req_sid,0,push_data )
						conn.Write([]byte( data ))

						// 发送点对点发送消息后 加入场景
						time.Sleep(100 * time.Millisecond)
						data =  protocol.WrapReqStr("JoinChannel",req_sid,req_id+1,"area-global" )
						conn.Write([]byte( data ))
						time.Sleep(5 * time.Second)


					}

					if cmd=="JoinChannel"  {

						push_data := fmt.Sprintf(`{"area_id":"area-global","data":"%s"}`,"md56666666666")
						data =  protocol.WrapBroatcastStr("Broadcast",req_sid,req_id+1, push_data )
						conn.Write([]byte( data ))
					}

					if cmd=="LeaveChannel"  {

						data =  protocol.WrapReqStr("KickSelf",req_sid,req_id+1,req_sid )
						conn.Write([]byte( data ))
					}
					if cmd=="KickSelf"  {

						conn.Close()
						return
					}
				}

				// 发送广播
				if _type==protocol.TypeBroadcast  {

					//fmt.Println("Broadcast:",msg_data)
					msg_err,form_sid,area_id,data := protocol.ParseRplyBrodcastData( str )
					if msg_err != nil {
						fmt.Println("broadcast reply error: ", msg_err.Error(),msg_arr)
						continue
					}
					fmt.Println( "broadcast recvice:", form_sid,area_id,data  )

					data =  protocol.WrapReqStr("LeaveChannel",req_sid,0,"area-global" )
					conn.Write([]byte( data ))

				}


			}

			return

		}( conn,  times, conn_num , int(i))


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
