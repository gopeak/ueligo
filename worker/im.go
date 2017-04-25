package worker

import (
	"github.com/antonholmquist/jason"
	"morego/golog"
	"fmt"
)








func (this TaskType)PushMessage(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data )
	fmt.Println( "PushMessage:",this.Sid, this.Data )
	data_json ,err_json:= jason.NewObjectFromBytes( []byte(this.Data ) )
	if( err_json!=nil ) {
		golog.Error("todpole message json err:",err_json.Error())
		return ""
	}

	to_sid,err1 := data_json.GetString("sid")
	message,err2 := data_json.GetString("msg")
	if( err1!=nil || err2!=nil   ){
		//golog.Error("todpole message json err:",err1.Error()+err2.Error()+err3.Error())
		return ""
	}

	sdk.Push( this.Sid,to_sid,message  )
	return "";


}

