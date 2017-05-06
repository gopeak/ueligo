package golang

import (
	"fmt"
	"strconv"
	"morego/golog"
	"morego/lib"
	"github.com/antonholmquist/jason"

)



func (this TaskType)Authorize(  ) ReturnType {

	// 从数据库中查询token是否有效
	db := new(lib.Mysql)
	_, err := db.Connect()
	if err != nil {
		//json_ret := fmt.Sprintf(`{"ret":"failed","type":"%s","id":"%s" ,"msg":"%s"}`,"failed",this.Sid,"数据库连接失败:" + err.Error())
		ret := ReturnType{ "failed","failed" ,this.Sid, "数据库连接失败:" + err.Error() }
		return ret
	}

	// 获取当前用户信息
	data_json ,err_json:= jason.NewObjectFromBytes( this.Data.([]byte) )
	if( err_json!=nil ) {
		ret := ReturnType{ "failed","failed" ,this.Sid, "解析认证数据失败:" + err_json.Error() }
		return ret
	}
	_token,_ := data_json.GetString("token")
	_sid,_ := data_json.GetString("sid")
	my_record := GetUserRow(db.Db, _sid )
	if( my_record["token"]==_token ){
		ret := ReturnType{ "ok","welcome" ,this.Sid,"认证成功"}
		return ret
	}else{
		ret := ReturnType{ "failed","failed" ,this.Sid, "认证失败"   }
		return ret
	}

}

func (this TaskType)SubscripeGroup(  ) ReturnType {

	// 从数据库中查询token是否有效
	db := new(lib.Mysql)
	_, err := db.Connect()
	if err != nil {
		ret := ReturnType{ "failed","subscripeGroupfailed" ,this.Sid, "数据库连接失败" + err.Error()  }
		return ret
	}

	// 获取当前用户信息
	data_json ,err_json:= jason.NewObjectFromBytes( this.Data.([]byte) )
	if( err_json!=nil ) {
		ret := ReturnType{ "failed","subscripeGroupfailed" ,this.Sid, "解析认证数据失败" + err_json.Error()  }
		return ret
	}
	_token,_ := data_json.GetString("token")
	_sid,_ := data_json.GetString("sid")
	my_record := GetUserRow(db.Db, _sid )

	if( my_record["token"]!=_token ){
		fmt.Println( "token:", my_record["token"], _token )
		ret := ReturnType{ "failed","subscripeGroupfailed" ,this.Sid, "数据认证错误" + err_json.Error()  }
		return ret
	}
	uid, _ := strconv.Atoi(my_record["id"])
	JoinChannel( db.Db, uid, my_record["sid"] )
	//json_ret := fmt.Sprintf(`{"ret":"ok","type":"%s","id":"%s" }`,"subscripeGroup",this.Sid)
	ret := ReturnType{ "ok","welcome" ,this.Sid,"认证成功"}
	return ret
}



func (this TaskType)PushMessage(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )
	fmt.Println( "PushMessage:",this.Sid, this.Data )
	data_json ,err_json:= jason.NewObjectFromBytes( this.Data.([]byte) )
	if( err_json!=nil ) {
		golog.Error(" push message json err:",err_json.Error())
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

func (this TaskType)PushGroupMessage(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid,this.Data.(string) )
	fmt.Println( "PushGroupMessage:",this.Sid, this.Data )
	data_json ,err_json:= jason.NewObjectFromBytes( this.Data.([]byte) )
	if( err_json!=nil ) {
		golog.Error(" push group  message json err:",err_json.Error())
		return ""
	}
	area_id,err1 := data_json.GetString("area_id")
	//message,err2 := data_json.GetString("msg")
	if( err1!=nil   ){
		golog.Error(" push group   message json err:",err1.Error() )
		return ""
	}
	sdk.Broatcast( this.Sid,area_id, this.Data.(string) )
	return "";
}


