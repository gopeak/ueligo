package golang

import (

	"net"
)


type TaskType struct {

	Conn * net.TCPConn

	Cmd string

	Sid string

	Reqid int

	Data string


}

func (this *TaskType) Init( conn *net.TCPConn,cmd string,sid string,reqid int,data string ) *TaskType{

	this.Cmd = cmd
	this.Sid = sid
	this.Reqid = reqid
	this.Data = data
	this.Conn = conn
	return this
}

