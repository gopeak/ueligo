package golang

import (
	"morego/protocol"
	"morego/lib/websocket"
	"net"
)

type TaskType struct {
	Conn *net.TCPConn

	WsConn *websocket.Conn

	Cmd string

	Sid string

	Reqid int

	Data interface{}

	ReqObj *protocol.ReqRoot
}

type ReturnType struct {
	Ret string `json:"ret"`

	Type string `json:"type"`

	Sid string `json:"sid"`

	Msg string `json:"msg"`
}

func (this *TaskType) Init(conn *net.TCPConn, req_obj *protocol.ReqRoot) *TaskType {

	//  cmd string,sid string,reqid int,data string
	this.Cmd = req_obj.Header.Cmd
	this.Sid = req_obj.Header.Sid
	this.Reqid = req_obj.Header.SeqId
	this.Data = req_obj.Data
	this.Conn = conn
	this.ReqObj = req_obj
	return this
}

func (this *TaskType) WsInit(wsconn *websocket.Conn, req_obj *protocol.ReqRoot) *TaskType {

	//  cmd string,sid string,reqid int,data string
	this.Cmd = req_obj.Header.Cmd
	this.Sid = req_obj.Header.Sid
	this.Reqid = req_obj.Header.SeqId
	this.Data = req_obj.Data
	this.WsConn = wsconn
	this.ReqObj = req_obj
	return this
}