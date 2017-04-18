package worker

import (
	"github.com/antonholmquist/jason"
	"morego/golog"
	"fmt"
	"strconv"
)





func (this TaskType)Authorize(  ) string {

	return this.Auth()


}


func (this TaskType)Message(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )

	data_json ,err_json:= jason.NewObjectFromBytes( []byte(this.Data ) )
	if( err_json!=nil ) {
		golog.Error("todpole message json err:",err_json.Error())
		return ""
	}

	_type,err1 := data_json.GetString("type")
	_,err2 := data_json.GetString("message")
	sid,err3 := data_json.GetString("id")
	if( err1!=nil || err2!=nil || err3!=nil ){
		//golog.Error("todpole message json err:",err1.Error()+err2.Error()+err3.Error())
		return ""
	}
	//broatcast_msg := fmt.Sprintf(`{"type":"message","message":"%s","id":"%s" }`,message,sid)
	sdk.Broatcast( sid,"area-global",this.Data  )
	json_ret := fmt.Sprintf(`{"type":"%s","id":"%s" }`,_type,sid)
	return json_ret;


}

func (this TaskType)Update(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )

	data_json ,err_json:= jason.NewObjectFromBytes( []byte(this.Data ) )
	if( err_json!=nil ) {
		golog.Error("todpole message json err:",err_json.Error())
		return ""
	}

	_type,_ := data_json.GetString("type")
	angle,_ := data_json.GetInt64("angle")
	_id,_ := data_json.GetInt64("id")
	momentum,_ := data_json.GetInt64("momentum")
	x,_ := data_json.GetInt64("x")
	y,_ := data_json.GetInt64("y")
	name,err_name := data_json.GetString("name")
	if( err_name!=nil ) {
		name = "Guest."+strconv.FormatInt(_id,10);
	}
	angle = angle+0
	momentum = momentum+0
	x = x+0
	y = y+0
	broatcast_data := fmt.Sprintf(`{"type":"%s","id":"%s","angle":%d,"momentum":%d,"x":%d,"y":%d,"life":1,"name":"%s","authorized":%s}`,
		_type,_id,angle,momentum,x,y,name,"false" )
	sdk.Broatcast( this.Sid,"area-global",broatcast_data )

	json_ret := fmt.Sprintf(`{"type":"%s","id":"%d" }`,"none",_id)
	return json_ret;


}



