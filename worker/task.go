package worker

import (
	"net"

	"fmt"
	"morego/area"
	"morego/web"
	"github.com/antonholmquist/jason"
	"morego/golog"
	"strings"
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

	//sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	sid:=area.CreateSid()
	if(   true  ){
		json_ret := fmt.Sprintf(`{"ret":"ok","type":"%s","id":"%s"  }`,"welcome",sid)
		return json_ret
	}else{
		json_ret := fmt.Sprintf(`{"ret":"failed","type":"%s","id":"%s" }`,"failed",sid)
		return json_ret
	}

}

func (this TaskType)Authorize(  ) string {

	// 从数据库中查询token是否有效
	db := new(web.Mysql)
	_, err := db.Connect()
	if err != nil {
		json_ret := fmt.Sprintf(`{"ret":"failed","type":"%s","id":"%s" ,"msg":"%s"}`,"failed",this.Sid,"数据库连接失败:" + err.Error())
		return json_ret
	}


	// 获取当前用户信息
	data_json ,err_json:= jason.NewObjectFromBytes( []byte(this.Data ) )
	if( err_json!=nil ) {
		json_ret := fmt.Sprintf(`{"ret":"failed","type":"%s","id":"%s" ,"msg":"%s"}`,"failed",this.Sid,"解析认证数据失败:" + err_json.Error())
		return json_ret
	}
	_token,_ := data_json.GetString("token")
	_sid,_ := data_json.GetString("sid")
	my_record := web.GetUserRow(db.Db, _sid )
	if( my_record["token"]==_token ){
		json_ret := fmt.Sprintf(`{"ret":"ok","type":"%s","id":"%s"  }`,"welcome",_sid)
		return json_ret

	}else{
		json_ret := fmt.Sprintf(`{"ret":"failed","type":"%s","id":"%s" }`,"failed",this.Sid)
		return json_ret
	}


}


func (this TaskType)Push(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	data_json ,err_json:= jason.NewObjectFromBytes( []byte(this.Data ) )
	if( err_json!=nil ) {
		golog.Error("todpole message json err:",err_json.Error())
		return ""
	}
	to_sid, _ := data_json.GetString("sid")
	to_data, _ := data_json.GetString("data")
	sdk.Push(to_sid, this.Sid, to_data)

	return "";

}


func (this TaskType)Broadcast(  ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	data_json ,err_json:= jason.NewObjectFromBytes( []byte(this.Data ) )

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



func (this TaskType)GetUserSession(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	return sdk.GetSession( this.Sid )

}

func (this TaskType)JoinChannel(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	fmt.Println( "JoinChannel",this.Data  )
	if(   sdk.ChannelAddSid( this.Sid ,this.Data ) ){
		return "ok"
	}else{
		return "failed"
	}

}


func (this TaskType)LeaveChannel(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )

	if(   sdk.ChannelKickSid( this.Sid ,this.Data ) ){
		return "ok"
	}else{
		return "failed"
	}

}


func (this TaskType)KickSelf(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )

	if(   sdk.Kick( this.Sid ) ){
		return "ok"
	}else{
		return "failed"
	}

}



func (this TaskType)GetBase( conn *net.TCPConn, cmd string, req_sid string ,req_id int,req_data string ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	return sdk.GetBase()

}



