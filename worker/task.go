package worker

import (
	"net"
)



type ReturnType struct {
	Data string
	Type  string
}





func (this ReturnType)auth( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {


	return "ok";

}