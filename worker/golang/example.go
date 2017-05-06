package golang

import (
	"net"
	"fmt"
	"morego/area"
	"github.com/antonholmquist/jason"
	"morego/golog"
	"strings"

)




func (this TaskType)Auth(  ) ReturnType {

	//sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )

	sid:=area.CreateSid()
	if(   sid!=""  ){
		//json_ret := fmt.Sprintf(`{"ret":"ok","type":"%s","id":"%s"  }`,"welcome",sid)
		ret := ReturnType{ "ok","welcome" ,sid,"" }
		return ret
	}else{
		//json_ret := fmt.Sprintf(`{"ret":"failed","type":"%s","id":"%s" }`,"failed",sid)

		ret := ReturnType{ "failed","failed" ,sid,"" }
		return ret
	}



}



func (this TaskType)Push(   ) interface{} {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid, this.Data.(string))
	data_json ,err_json:= jason.NewObjectFromBytes( this.Data.([]byte) )
	if( err_json!=nil ) {
		golog.Error("todpole message json err:",err_json.Error())
		return ""
	}
	to_sid, _ := data_json.GetString("sid")
	to_data, _ := data_json.GetString("data")
	sdk.Push(to_sid, this.Sid, to_data)

	return "";

}


func (this TaskType)Broadcast(  ) interface{}{

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )
	data_json ,err_json:= jason.NewObjectFromBytes( this.Data.([]byte) )

	if ( err_json != nil ) {
		golog.Error("broatcast data json format error")
		return ""
	}
	from_sid := this.Sid
	area_id, _ := data_json.GetString("area_id")
	to_data, _ := data_json.GetString("data")
	to_data = strings.Replace(to_data, "\n", "", -1)
	if( area_id=="global" ) {
		golog.Error("broatcast global failed")
		return ""
	}else{
		sdk.Broatcast( from_sid, area_id,to_data )
	}
	return ""
}



func (this TaskType)GetUserSession(   ) interface{} {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )
	return sdk.GetSession( this.Sid )

}

func (this TaskType)JoinChannel(   ) interface{} {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )
	fmt.Println( "JoinChannel",this.Data  )
	if(   sdk.ChannelAddSid( this.Sid ,this.Data.(string) ) ){
		return "ok"
	}else{
		return "failed"
	}

}


func (this TaskType)LeaveChannel(   ) interface{} {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )

	if(   sdk.ChannelKickSid( this.Sid ,this.Data.(string) ) ){
		return "ok"
	}else{
		return "failed"
	}

}


func (this TaskType)KickSelf(   ) interface{} {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )

	if(   sdk.Kick( this.Sid ) ){
		return "ok"
	}else{
		return "failed"
	}

}



func (this TaskType)GetBase( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )
	return sdk.GetBase()

}



