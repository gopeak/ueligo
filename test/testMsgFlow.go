package main

import (
	"bufio"
	//"crypto/md5"
	"fmt"
	"morego/protocol"
	"morego/worker/golang"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
	"strings"
	"encoding/json"
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
		//data :=  protocol.WrapReqStr("Auth", strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(i) , i, strconv.FormatInt(int64(time.Now().Unix()), 10) ) // fmt.Sprintf("%d||%s||%x||%d||%d\n", protocol.TypeReq, "Auth",  md5.Sum([]byte(srcData)) , i, time.Now().Unix())

		// str:=fmt.Sprintf("%d||%s||%s||%d||%s\n" ,TypeReq, cmd,from_sid ,req_id, data) ;
		protocolPack:= new(protocol.Pack)
		protocolPack.Init()
		req_obj_header := &protocol.ReqHeader{}
		req_obj_header.Cmd = "Auth"
		req_obj_header.Sid = strconv.FormatInt(time.Now().Unix(), 10)+strconv.Itoa(i)
		req_obj_header.NoResp = false
		req_obj_header.Token = ""
		req_obj_header.SeqId = i
		req_obj_header.Version = "1.0"
		header_buf ,_ := json.Marshal( req_obj_header )
		data :=  strconv.FormatInt(int64(time.Now().Unix()), 10)
		buf,_ := protocol.EncodePacket( protocol.TypeReq,header_buf,[]byte(data)  )

		conn.Write( buf )
		r := bufio.NewReader(conn)
		for {
			_, resp_header,resp_data,err :=  protocol.DecodePacket( r )
			fmt.Println( "Auth:",  resp_header, resp_data )

			ret_obj := new(golang.ReturnType)
			 json.Unmarshal( resp_data, ret_obj )
			if err != nil {
				fmt.Println(" protocol.DecodePacket err: ", err.Error() )
				continue
			}
			Sids  = append(  Sids, ret_obj.Sid )
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
			//data :=  protocol.WrapReqStr("GetUserSession",req_sid,0,req_sid )
			protocolPack:= new(protocol.Pack)
			protocolPack.Init()
			req_obj_header := &protocol.ReqHeader{}
			req_obj_header.Cmd = "GetUserSession"
			req_obj_header.Sid = req_sid
			req_obj_header.SeqId = 0
			header_buf ,_ := json.Marshal( req_obj_header )

			buf,_ := protocol.EncodePacket( protocol.TypeReq,header_buf,[]byte(req_sid)  )
			conn.Write([]byte( buf ))

			for {

				_type, resp_header,resp_data,err :=  protocol.DecodePacket( reader )
				fmt.Println( "Auth:",  resp_header, resp_data )
				if err != nil {
					fmt.Println("HandleConn connection error: ", err.Error())
					break
				}
				success++


				if( _type==protocol.TypeReply ){
					resp_header_obj,msg_err := protocolPack.GetRespHeaderObj( resp_header )
					if msg_err != nil {
						fmt.Println("msg error: ", msg_err.Error() )
						continue
					}
					// 登录认证,然后获取用户信息
					req_id := resp_header_obj.SeqId
					// 获取当前信息后 发送点对点信息
					if resp_header_obj.Cmd=="GetUserSession"  {

						to_sid_index := i-1
						if to_sid_index<0 {
							to_sid_index = 0
						}
						to_sid:= Sids[to_sid_index]
						push_data := fmt.Sprintf(`{"sid":"%s","data":"%s"}`,to_sid,"md55555555555")
						req_obj_header := &protocol.ReqHeader{}
						req_obj_header.Cmd = "Push"
						req_obj_header.Sid = req_sid
						req_obj_header.SeqId = 0
						header_buf ,_ := json.Marshal( req_obj_header )
						buf,_ := protocol.EncodePacket( protocol.TypeReq,header_buf,[]byte(push_data)  )
						conn.Write([]byte( buf ))

						// 发送点对点发送消息后 加入场景
						time.Sleep(100 * time.Millisecond)
						req_obj_header.Cmd = "JoinChannel"
						req_obj_header.Sid = req_sid
						req_obj_header.SeqId = req_id+1
						header_buf ,_ = json.Marshal( req_obj_header )
						buf,_ = protocol.EncodePacket( protocol.TypeReq,header_buf,[]byte("area-global")  )
						conn.Write([]byte( buf ))

						time.Sleep(5 * time.Second)

					}

					if resp_header_obj.Cmd=="JoinChannel"  {
						req_obj_header := &protocol.ReqHeader{}
						req_obj_header.Cmd = "Broadcast"
						req_obj_header.Sid = req_sid
						req_obj_header.SeqId = req_id+1
						header_buf ,_ := json.Marshal( req_obj_header )
						push_data := fmt.Sprintf(`{"area_id":"area-global","data":"%s"}`,"md56666666666")
						buf,_ := protocol.EncodePacket( protocol.TypeReq,header_buf,[]byte(push_data)  )
						conn.Write([]byte( buf ))
					}

					if resp_header_obj.Cmd=="LeaveChannel"  {

						req_obj_header := &protocol.ReqHeader{}
						req_obj_header.Cmd = "KickSelf"
						req_obj_header.Sid = req_sid
						req_obj_header.SeqId = req_id+1
						header_buf ,_ := json.Marshal( req_obj_header )
						buf,_ := protocol.EncodePacket( protocol.TypeReq,header_buf,[]byte(req_sid)  )
						conn.Write([]byte( buf ))
					}
					if resp_header_obj.Cmd=="KickSelf"  {

						conn.Close()
						return
					}
				}

				// 发送广播
				if _type==protocol.TypeBroadcast  {

					//fmt.Println("Broadcast:",msg_data)
					resp_broatcast_obj,msg_err  := protocolPack.GetBroatcastHeaderObj( resp_header )
					if msg_err != nil {
						fmt.Println("broadcast reply error: ", msg_err.Error(),resp_broatcast_obj)
						continue
					}
					fmt.Println( "broadcast recvice:",resp_broatcast_obj, string(resp_data)  )

					req_obj_header := &protocol.ReqHeader{}
					req_obj_header.Cmd = "LeaveChannel"
					req_obj_header.Sid = req_sid
					req_obj_header.SeqId = 0
					header_buf ,_ := json.Marshal( req_obj_header )
					buf,_ := protocol.EncodePacket( protocol.TypeReq,header_buf,[]byte("area-global" )  )
					conn.Write([]byte( buf ))
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
