package protocol

import (
	"encoding/json"
)




type Json struct {

	ProtocolObj ProtocolType
	Data []byte
	Init func()

}

func (this *Json) Init(  ) *Json{

	this.ProtocolObj = new( ProtocolType)
	this.ProtocolObj.ReqObj = new( ReqRoot )
	this.ProtocolObj.RespObj = new( ResponseRoot )
	this.ProtocolObj.BroatcastObj = new( BroatcastRoot )
	this.ProtocolObj.PushObj = new( PushRoot )
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


func (this *Json)SetReqObj(  ReqObj ReqRoot   )  {

	 this.ProtocolObj.ReqObj = ReqObj
}

func (this *Json)SetBroatcastObj(  BroatcastObj BroatcastRoot )  {

	this.ProtocolObj.BroatcastObj = BroatcastObj
}

func (this *Json)SetPushObj(  PushObj PushRoot )  {

	this.ProtocolObj.PushObj = PushObj
}

