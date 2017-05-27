package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/antonholmquist/jason"
	"morego/util"
)

type Json struct {
	ProtocolObj ProtocolType
	Data        []byte
}

func (this *Json) Init() *Json {

	this.ProtocolObj = ProtocolType{}
	this.ProtocolObj.ReqObj = ReqRoot{}
	this.ProtocolObj.RespObj = ResponseRoot{}
	this.ProtocolObj.BroatcastObj = BroatcastRoot{}
	this.ProtocolObj.PushObj = PushRoot{}
	return this
}



func (this *Json) GetReqObj(data []byte) (*ReqRoot, error) {
	this.Data = data
	stb := &ReqRoot{}

	json_obj,err := jason.NewObjectFromBytes( data )
	req_header_obj := ReqHeader{}
	req_header_obj.Cmd , _ = json_obj.GetString("header","cmd")
	seq_id  , _ := json_obj.GetInt64("header","seq_id")
	req_header_obj.SeqId = int( seq_id )
	req_header_obj.Gzip ,  _ = json_obj.GetBoolean("header","gzip")
	req_header_obj.Sid ,   _ = json_obj.GetString("header","sid")
	req_header_obj.Token , _ = json_obj.GetString("header","token")
	stb.Type,_ = json_obj.GetString("type")
	data_mix ,_:= json_obj.GetInterface("data")
	data_buf := util.Convert2Byte(data_mix)
	fmt.Println( "ws GetReqObj data_str:",string(data_buf) )
	stb.Header = req_header_obj
	stb.Data = data_buf
	return stb, err
}

func (this *Json) GetRespObj(data []byte) (*ResponseRoot, error) {
	this.Data = data
	stb := &ResponseRoot{}
	err := json.Unmarshal(data,  stb)

	//this.ProtocolObj.RespObj = stb
	return stb, err
}

func (this *Json) GetBroatcastObj(data []byte) (*BroatcastRoot, error) {
	this.Data = data
	stb := &BroatcastRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.BroatcastObj = stb
	return stb, err
}

func (this *Json) GetPushObj(data []byte) (*PushRoot, error) {
	this.Data = data
	stb := &PushRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.PushObj = stb
	return stb, err
}

func (this *Json) WrapRespObj( req_obj *ReqRoot, invoker_ret []byte, status int ) ResponseRoot {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = req_obj.Header.Cmd
	resp_header_obj.SeqId = req_obj.Header.SeqId
	resp_header_obj.Gzip = req_obj.Header.Gzip
	resp_header_obj.Sid = req_obj.Header.Sid
	resp_header_obj.Status = status
	this.ProtocolObj.RespObj.Header =resp_header_obj
	this.ProtocolObj.RespObj.Data = invoker_ret
	this.ProtocolObj.RespObj.Type = TypeResp

	return this.ProtocolObj.RespObj
}

func (this *Json) WrapResp(   header []byte, data []byte, status int, msg string )  []byte  {

	header = util.TrimX001( header )
	data = util.TrimX001( data )
	data_str := string(data)
	if( util.TrimStr(data_str)==""){
		data_str = "{}"
	}
	header_str := string(header)
	if( util.TrimStr(header_str)==""){
		header_str = "{}"
	}
	return []byte(fmt.Sprintf(`{"type":"%s","status":%d,"msg":"%s","header":%s,"data":%s}`,
			TypeResp, status, msg, header_str ,data_str ))


}

func (this *Json) WrapPushRespObj(to_sid string, from_sid string , data string ) PushRoot {

	push_header_obj := PushHeader{}
	push_header_obj.Sid = from_sid
	var map_data map[string]interface{}
	err:=json.Unmarshal( []byte(data), map_data )
	var to_data_buf []byte
	if err==nil {
		tmp ,err:=json.Marshal( map_data )
		if err==nil {
			to_data_buf = tmp
		}
	}else{
		to_data_buf =  []byte(data)
	}

	push_obj := PushRoot{}
	push_obj.Header =push_header_obj
	push_obj.Data  = to_data_buf
	push_obj.Type  = TypeResp

	return push_obj
}


func (this *Json) WrapBroatcastRespObj(channel_id string, from_sid string , data []byte) BroatcastRoot {

	broatcast_header_obj := BroatcastHeader{}
	broatcast_header_obj.Sid = from_sid
	broatcast_header_obj.ChannelId = channel_id

	broatcast_obj := BroatcastRoot{}
	broatcast_obj.Header =broatcast_header_obj
	broatcast_obj.Data  = data
	broatcast_obj.Type  = "broatcast"

	return broatcast_obj
}

/**
 * 封包返回客户端错误的消息
 */
func (this *Json) WrapRespErr(err string) []byte {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = "WrapRespErr"
	resp_header_obj.SeqId = 0
	resp_header_obj.Sid = ""
	resp_header_obj.Status = 500
	this.ProtocolObj.RespObj.Header =resp_header_obj
	this.ProtocolObj.RespObj.Data = []byte(err)
	this.ProtocolObj.RespObj.Type = TypeError

	buf,_ := json.Marshal( this.ProtocolObj.RespObj )

	return buf
}


