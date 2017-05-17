package main

import (
	"bufio"
	//"crypto/md5"
	"fmt"
	"morego/protocol"
	"net"
	"time"
)

func main() {

	go Server("0.0.0.0", 7003)

	time.Sleep(3 * time.Second)

	go client_side()

	select {

	}

}


/**
 * 监听客户端连接
 */
func Server(ip string, port int) {


	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), port, ""})
	if err != nil {
		fmt.Println("ListenTCP Exception:", err.Error())
		return
	}
	// 初始化
	fmt.Println("  Server :", ip, port)
	for {
		conn, err := listen.AcceptTCP()
		defer conn.Close()
		if err != nil {
			fmt.Println("AcceptTCP Exception::", err.Error())
			continue
		}
		// 校验ip地址
		conn.SetKeepAlive(true)

		go handleClientMsg(conn)
	}
}
func handleClientMsg(conn *net.TCPConn) {

	//声明一个管道用于接收解包的数据
	reader := bufio.NewReader(conn)
	fmt.Println("HandleConn client: ", conn.RemoteAddr() )
	for {

		header, data, err := protocol.DecodePacket( reader )
		fmt.Println("  server recvice header: ", string(header), " data:", string(data))
		if err != nil {
			fmt.Println("HandleConn connection error: ", err.Error())
			break
		}


	}
}

func client_side() {

	// 客户端请求
	fmt.Println("  client_side " )
	service := "127.0.0.1:7003"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if ( err != nil ) {
		fmt.Println("  net.DialTCP error: ", err.Error())
		return
	}

	//fmt.Println( conn )
	reader := bufio.NewReader(conn)

	header := []byte( `{"cmd":"Auth","sid":"123456"}`)
	data := []byte(`{"user":"simarui","pass":123"}`)
	buf,err := protocol.EncodePacket( header, data)
	if ( err != nil ) {
		fmt.Println(" protocol.EncodePacket error: ", err.Error())
		return
	}
	fmt.Println(" conn.Write: ", string(buf) )
	conn.Write( buf )
	for {
		resp_header, resp_data, err := protocol.DecodePacket(reader)
		if err != nil {
			fmt.Println(" connection error: ", err.Error())
			break
		}
		fmt.Println("  resp header: ", string(resp_header), " data:", string(resp_data))
		break
	}

}