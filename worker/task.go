package worker

import (
	"net"

	"fmt"
)



type ReturnType struct {
	Data string
	Type  string
}





func (this ReturnType)Auth( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {


	return "ok";

}


func (this ReturnType)GetUserSession( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {

	sdk:=Sdk{ Connected:false,Cmd:cmd,Sid:req_sid,Reqid:req_id,Data:req_data}
	this.Data=sdk.GetSessionStr( req_sid )
	fmt.Println( this.Data )
	return this.Data

}

func (this ReturnType)JoinChannel( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {

	sdk:=Sdk{ Connected:false,Cmd:cmd,Sid:req_sid,Reqid:req_id,Data:req_data}
	if(   sdk.ChannelAddSid( req_sid ,req_data) ){
		return "ok"
	}else{
		return "failed"
	}

}


func (this ReturnType)GetBase( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {

	sdk:=Sdk{ Connected:false,Cmd:cmd,Sid:req_sid,Reqid:req_id,Data:req_data}
	return sdk.GetBase()

}



