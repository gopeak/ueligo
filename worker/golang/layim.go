package golang

import (
	"fmt"
	"strconv"
	"morego/lib"

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
	data := this.Data.(map[string]interface{})
	_token := data["token"].(string)
	_sid := data["sid"].(string)
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
	data := this.Data.(map[string]interface{})
	_token := data["token"].(string)
	_sid := data["sid"].(string)
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

	data := this.Data.(map[string]interface{})
	to_sid := data["sid"].(string)

	GetBaseCallback := func( resp string ) string {

		fmt.Println( "GetBaseCallback:", resp )
		return ""
	}
	sdk.ReqHubAsync( "GetBase", "", GetBaseCallback )

	sdk.Push(  this.Sid, to_sid,  data )
	return "";
}

func (this TaskType)PushGroupMessage(   ) string {

	sdk:=new(Sdk).Init(this.Cmd,this.Sid,this.Reqid, this.Data  )
	//fmt.Println( "PushGroupMessage:",this.Sid, this.Data )

	data := this.Data.(map[string]interface{})
	area_id := data["area_id"].(string)

	sdk.Broatcast( this.Sid,area_id, data )
	return "";
}


