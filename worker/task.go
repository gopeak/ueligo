package worker

import (
	"net"

	"fmt"
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


func (this TaskType)Auth(  ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	fmt.Println(sdk)
	return "ok";

}


func (this TaskType)GetUserSession(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	return sdk.GetSession( this.Sid )

}

func (this TaskType)JoinChannel(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )

	if(   sdk.ChannelAddSid( this.Sid ,this.Data ) ){
		return "ok"
	}else{
		return "failed"
	}

}


func (this TaskType)GetBase( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	return sdk.GetBase()

}



