package golang

import (
	"fmt"
	"strconv"
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
	fmt.Println( "Authorize this.Data:",string(this.Data) )
	json_obj,err := jason.NewObjectFromBytes( this.Data )
	if err != nil {
		ret := ReturnType{ "failed","failed" ,this.Sid, "json err:" + err.Error() }
		return ret
	}
	_token ,err:= json_obj.GetString("token")
	if err != nil {
		ret := ReturnType{ "failed","failed" ,this.Sid, "json err:" + err.Error() }
		return ret
	}
	_sid ,err:= json_obj.GetString("sid")
	if err != nil {
		ret := ReturnType{ "failed","failed" ,this.Sid, "json err:" + err.Error() }
		return ret
	}
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
	json_obj,_ := jason.NewObjectFromBytes( this.Data )
	_token ,_:= json_obj.GetString("token")
	_sid,_ := json_obj.GetString("sid")
	my_record := GetUserRow(db.Db, _sid )

	if( my_record["token"]!=_token ){
		fmt.Println( "token:", my_record["token"], _token )
		ret := ReturnType{ "failed","subscripeGroupfailed" ,this.Sid, "数据认证错误"    }
		return ret
	}
	uid, _ := strconv.Atoi(my_record["id"])
	JoinChannel( db.Db, uid, my_record["sid"] )
	ret := ReturnType{ "ok","subscripeGroup" ,this.Sid,"订阅群组消息成功"}
	return ret
}



func (this TaskType)PushMessage(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid, this.Data )
	fmt.Println( "PushMessage:", this.Sid, this.Data )

	json_obj,_ := jason.NewObjectFromBytes( this.Data )
	to_sid ,_:= json_obj.GetString("sid")

	GetBaseCallback := func( resp string ) string {

		fmt.Println( "GetBaseCallback:", resp )
		return ""
	}
	sdk.ReqHubAsync( "GetBase", "", GetBaseCallback )

	sdk.Push(  this.Sid, to_sid,  this.Data )
	return "";
}

func (this TaskType)PushGroupMessage(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid, this.Data  )
	//fmt.Println( "PushGroupMessage:",this.Sid, this.Data )
	json_obj,_ := jason.NewObjectFromBytes( this.Data )
	area_id ,_:= json_obj.GetString("area_id")

	sdk.Broatcast( this.Sid,area_id, this.Data )
	return "";
}


