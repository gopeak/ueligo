package worker

import (
	"net"

)



type ReturnType struct {
	Data string
	Type  string
}





func (this ReturnType)Auth( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {


	return "ok";

}


func (this ReturnType)GetUserSession( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {

	return GetSessionStr( req_sid )

}