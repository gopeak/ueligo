// automatically generated, do not modify

package protocol

import (
	"strings"
	"strconv"
	"errors"
	"fmt"
)


// 使用字符串分割时数据块定义
const (
	MSG_TYPE_INDEX    = 0
	MSG_CMD_INDEX     = 1
	MSG_SID_INDEX =  2
	MSG_REQID_INDEX   = 3
	MSG_CHANNEL_INDEX   = 3
	MSG_DATA_INDEX   = 4
)

// 数据类型
const (
	TypeReq = 1
	TypeReply = 2
	TypePush = 3
	TypeBroadcast = 4
	TypeError = 5
	TypeJoinChannel = 6
	TypePing = 7
)

/**
 * 封包返回客户端错误的消息
 */
func WrapRespErrStr(err string) string {
	str := fmt.Sprintf("%d||%s||%s||%d||%s\n", TypeError, "", "", 0, err)
	return str
}

/**
 * 客户端封包返回数据
 */
func WrapRespStr(cmd string, from_sid string, req_id int, data string) string {
	str := fmt.Sprintf("%d||%s||%s||%d||%s\n", TypeReply, cmd, from_sid, req_id, data)
	return str
}

/**
 * 客户端封包请求数据
 */
func WrapReqStr(  cmd string, from_sid string, req_id int,data string ) string {
	str:=fmt.Sprintf("%d||%s||%s||%d||%s\n" ,TypeReq, cmd,from_sid ,req_id, data) ;
	return str
}


// 解析服务器端请求返回的数据
func ParseReqData(str string) ( error ,int ,string,string,int,string){

	msg_arr := strings.Split(str, "||")
	var err error
	err = nil
	if len(msg_arr) < 4 {
		err = errors.New("ParseReqData data length error")
		return err,0,"","",0,""
	}
	_type,_ := strconv.Atoi(msg_arr[MSG_TYPE_INDEX])
	cmd := msg_arr[MSG_CMD_INDEX];
	req_sid := msg_arr[MSG_SID_INDEX]
	req_id ,_ :=strconv.Atoi(msg_arr[MSG_REQID_INDEX])
	req_data := msg_arr[MSG_DATA_INDEX]
	// 去除换行符
	req_data = strings.Replace(req_data, "\n", "", -1)
	return err,_type,cmd,req_sid,req_id,req_data
}

// 解析服务器端请求返回的数据
func ParseRplyData(str string) ( error ,int ,string,string,int,string){

	msg_arr := strings.Split(str, "||")
	var err error
	err = nil
	if len(msg_arr) < 5 {
		err = errors.New("request data length error")
		return err,0,"","",0,""
	}
	_type,_ := strconv.Atoi(msg_arr[MSG_TYPE_INDEX])
	cmd := msg_arr[MSG_CMD_INDEX];
	req_sid := msg_arr[MSG_SID_INDEX]
	req_id ,_ :=strconv.Atoi(msg_arr[MSG_REQID_INDEX])
	req_data := msg_arr[MSG_DATA_INDEX]
	// 去除换行符
	req_data = strings.Replace(req_data, "\n", "", -1)
	return err,_type,cmd,req_sid,req_id,req_data
}

/**
 * 客户端封包请求数据
 */
func WrapPushStr(  cmd string, from_sid string, req_id int,data string ) string {
	str:=fmt.Sprintf("%d||%s||%s||%s||%s\n" ,TypePush, cmd,from_sid ,req_id, data) ;
	return str
}

func WrapPushRespStr(  from_sid string, data string ) string {
	str:=fmt.Sprintf("%d||%s||%s||0||%s\n" ,TypePush, "",from_sid ,data) ;
	return str
}

// 解析服务器端PUSH过来的消息
func ParseRplyPushData(str string) ( error , string, string ){

	msg_arr := strings.Split(str, "||")
	var err error
	err = nil
	if len(msg_arr) < 3 {
		err = errors.New("request data length error")
		return err,"",""
	}

	from_sid := msg_arr[MSG_SID_INDEX]
	push_data := msg_arr[MSG_DATA_INDEX]
	// 去除换行符
	push_data = strings.Replace(push_data, "\n", "", -1)
	return err,from_sid,push_data
}

/**
 * 客户端封包广播数据
 */
func WrapBroatcastStr(  cmd string, from_sid string, req_id int,data string ) string {
	str:=fmt.Sprintf("%d||%s||%s||%d||%s\n" ,TypeBroadcast, cmd,from_sid ,req_id, data) ;
	return str
}

func WrapBroatcastRespStr(  from_sid string, area_id string, data string ) string {
	str:=fmt.Sprintf("%d||%s||%s||%s||%s\n" , TypeBroadcast,"",from_sid ,area_id,data) ;
	return str
}

func ParseRplyBrodcastData(str string) ( error ,string, string, string ){

	msg_arr := strings.Split(str, "||")
	var err error
	err = nil
	if len(msg_arr) < 4 {
		err = errors.New("request data length error")
		return err,"","",""
	}
	from_sid := msg_arr[MSG_SID_INDEX]
	area_id := msg_arr[MSG_CHANNEL_INDEX]
	broadcast_data := msg_arr[MSG_DATA_INDEX]
	// 去除换行符
	broadcast_data = strings.Replace(broadcast_data, "\n", "", -1)
	return err,from_sid,area_id,broadcast_data
}

// 封装请求Hub数据
func WrapReqHubStr(  cmd string, from_sid string, req_id int,data string ) string {
	str:=fmt.Sprintf("%s||%s||%d||%s\n" , cmd,from_sid ,req_id, data) ;
	return str
}

/**
 * 客户端封包返回数据
 */
func WrapHubRespStr(cmd string, from_sid string, req_id int, data string) string {
	str := fmt.Sprintf("%s||%s||%d||%s\n",  cmd, from_sid, req_id, data)
	return str
}

/**
 * 封包返回客户端错误的消息
 */
func WrapHubRespErrStr(err string,cmd string ) string {
	str := fmt.Sprintf("%d||%s||%s||%d||%s\n", TypeError, cmd, "", 0, err)
	return str
}

/**
 * 解析worker server返回的数据
 */
func ParseHubRplyData(str string) ( error ,string,string,int,string){

	msg_arr := strings.Split(str, "||")
	var err error
	err = nil
	if len(msg_arr) < 4 {
		err = errors.New("request data length error")
		return err,"","",0,""
	}
	cmd := msg_arr[0];
	req_sid := msg_arr[1]
	req_id ,_ :=strconv.Atoi(msg_arr[2])
	req_data := msg_arr[3]
	// 去除换行符
	req_data = strings.Replace(req_data, "\n", "", -1)
	return err,cmd,req_sid,req_id,req_data
}
/**
 * 解析Hub请求的数据
 */
func ParseHubReqData(str string) ( error ,string,string,int,string){


	msg_arr := strings.Split(str, "||")
	var err error
	err = nil

	if len(msg_arr) < 4 {
		err = errors.New("request data length error:"+str)
		return err,"","",0,""
	}
	cmd := msg_arr[0];
	req_sid := msg_arr[1]
	req_id ,_ :=strconv.Atoi(msg_arr[2])
	req_data := msg_arr[3]
	// 去除换行符
	req_data = strings.Replace(req_data, "\n", "", -1)
	return err,cmd,req_sid,req_id,req_data
}
