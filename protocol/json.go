package protocol

import (
	"encoding/json"
)




type Json struct {

	ProtocolObj ProtocolType
	Data []byte

}

func (this *Json) Init(  ) *Json{

	this.ProtocolObj =   ProtocolType{}
	this.ProtocolObj.ReqObj = ReqRoot{}
	this.ProtocolObj.RespObj = ResponseRoot{}
	this.ProtocolObj.BroatcastObj = BroatcastRoot{}
	this.ProtocolObj.PushObj = PushRoot{}
	return this
}

func (this *Json)GetReqObj(  data []byte  ) (ReqRoot , error) {
	this.Data = data
	err:=json.Unmarshal( data, this.ProtocolObj.ReqObj)
	return this.ProtocolObj.ReqObj ,err
}

func (this *Json)GetRespObj(  data []byte  ) (ResponseRoot , error) {
	this.Data = data
	err:=json.Unmarshal( data, this.ProtocolObj.RespObj)
	return this.ProtocolObj.RespObj ,err
}


func (this *Json)GetBroatcastObj(  data []byte  )  (BroatcastRoot , error)  {
	this.Data = data
	err:=json.Unmarshal( data, this.ProtocolObj.BroatcastObj)
	return this.ProtocolObj.BroatcastObj,err
}

func (this *Json)GetPushObj(  data []byte  )  (PushRoot , error)  {
	this.Data = data
	err:=json.Unmarshal( data, this.ProtocolObj.PushObj)
	return this.ProtocolObj.PushObj,err
}


func (this *Json) WrapRespObj(  req_obj ReqRoot,invoker_ret interface{} ,status int ,msg string )  ResponseRoot   {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = req_obj.Header.Cmd
	resp_header_obj.SeqId = req_obj.Header.SeqId
	resp_header_obj.Gzip = req_obj.Header.Gzip
	resp_header_obj.Sid = req_obj.Header.Sid
	this.ProtocolObj.RespObj.Data = invoker_ret
	this.ProtocolObj.RespObj.Status = status
	this.ProtocolObj.RespObj.Msg = msg

	return this.ProtocolObj.RespObj
}

