//通讯协议处理，主要处理封包和解包的过程
package protocol


import (

	"encoding/binary"
	"errors"
	"fmt"
	"bufio"
	"bytes"
	"encoding/json"
	"./worker"
)


//|总长度4|type4|头长度4|header|data|checksum
type ClientPacket struct {
	TotalSize uint32    //4
	Type      uint8
	HeaderSize uint32   //4
	Header   []byte
	Payload   []byte
	Checksum  uint32    //4
}


func EncodePacket(  _type int32, header []byte,  payload []byte) ( []byte ,error ){
	// len(totaol)+ len(header) + len(Checksum) == 12
	var pkg *bytes.Buffer = new(bytes.Buffer)
	totalsize := uint32(len(string(header)) +  len(string(payload)) )

	// set totalsize
	err:=binary.Write( pkg , binary.LittleEndian, totalsize)
	if err != nil {
		return nil, err
	}
	fmt.Println( "set totalsize" , totalsize )

	// set type
	err = binary.Write( pkg, binary.LittleEndian, _type)
	if err != nil {
		return nil, err
	}
	// set headersize
	headersize := uint32(len(string(header)))
	err = binary.Write( pkg, binary.LittleEndian, headersize)
	if err != nil {
		return nil, err
	}
	fmt.Println( "set theadersize" , headersize )

	// set header
	err = binary.Write( pkg, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}

	// set payload
	err = binary.Write(pkg, binary.LittleEndian, payload )
	if err != nil {
		return nil, err
	}
	// write checksum
	return  pkg.Bytes(),nil
}


func DecodePacket(r *bufio.Reader) ( uint32, []byte,  []byte, error) {
	var totalsize , headersize uint32

	lengthByte, _ := r.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	err := binary.Read(lengthBuff, binary.LittleEndian, &totalsize)
	if err != nil {
		return 0,nil,nil,errors.New("read total size error") // errors.Annotate(err, "read total size")
	}
	fmt.Println( "totalsize" , totalsize)
	if totalsize < 12 {
		return 0,nil, nil,errors.New( fmt.Sprintf("bad packet. totalsize:%d", totalsize))
	}

	pack := make([]byte, 16+int(totalsize))
	_, err = r.Read(pack)
	if err != nil {
		return 0,nil,nil, err
	}
	if len(pack)<4 {
		return 0,nil,nil, errors.New("read headersize error")
	}

	_type :=   uint32(pack[4] )
	fmt.Println( "type:" ,_type )

	headersize =   uint32(pack[8] )
	fmt.Println( "headersize"  ,headersize )
	header :=  pack[12:headersize+12]

	payload := pack[(12+headersize):(totalsize+16)]
	fmt.Println( "header:" , string(header))
	fmt.Println( "payload:" ,  string(payload) )
	return _type,header,payload, nil
}


type Pack struct {
	ProtocolObj ProtocolType
	Data        []byte
}

func (this *Pack) Init() *Json {

	this.ProtocolObj = ProtocolType{}
	this.ProtocolObj.ReqObj = ReqRoot{}
	this.ProtocolObj.RespObj = ResponseRoot{}
	this.ProtocolObj.BroatcastObj = BroatcastRoot{}
	this.ProtocolObj.PushObj = PushRoot{}
	return this
}

func (this *Pack) GetReqObjByReader( reader *bufio.Reader ) (*ReqRoot, error) {

	stb := &ReqRoot{}
	_type,header,data,err := DecodePacket( reader )
	if err!=nil {
		return stb, err
	}
	return this.GetReqObj( _type,header,data )

}

func (this *Pack) GetReqObj( _type int32 ,header []byte, data []byte ) (*ReqRoot, error) {

	stb := &ReqRoot{}

	stb.Type = fmt.Sprintf( "%d", _type )
	err :=json.Unmarshal(header, stb.Header)
	if err!=nil {
		return stb, err
	}
	err = json.Unmarshal(data, stb.Data)
	//this.ProtocolObj.ReqObj = stb
	return stb, err
}


func (this *Pack) GetRespHeaderObj(  header []byte) (*RespHeader, error) {

	stb := &RespHeader{}
	err := json.Unmarshal( header, stb )
	return stb, err
}

func (this *Pack) GetRespObj(  data []byte) (*ResponseRoot, error) {
	this.Data = data
	stb := &ResponseRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.RespObj = stb
	return stb, err
}

func (this *Pack) GetBroatcastObj(data []byte) (*BroatcastRoot, error) {
	this.Data = data
	stb := &BroatcastRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.BroatcastObj = stb
	return stb, err
}

func (this *Pack) GetPushObj(data []byte) (*PushRoot, error) {
	this.Data = data
	stb := &PushRoot{}
	err := json.Unmarshal(data, stb)
	//this.ProtocolObj.PushObj = stb
	return stb, err
}


func (this *Pack) WrapRespObj( req_obj *ReqRoot, invoker_ret interface{}, status int ) []byte {

	resp_header_obj := RespHeader{}
	resp_header_obj.Cmd = req_obj.Header.Cmd
	resp_header_obj.SeqId = req_obj.Header.SeqId
	resp_header_obj.Gzip = req_obj.Header.Gzip
	resp_header_obj.Sid = req_obj.Header.Sid
	resp_header_obj.Status = status
	this.ProtocolObj.RespObj.Header =resp_header_obj
	this.ProtocolObj.RespObj.Data = invoker_ret
	this.ProtocolObj.RespObj.Type = "2"


	return this.ProtocolObj.RespObj
}

func (this *Pack) WrapPushRespObj(to_sid string, from_sid string , data interface{}) PushRoot {

	push_header_obj := PushHeader{}
	push_header_obj.Sid = from_sid

	push_obj := PushRoot{}
	push_obj.Header =push_header_obj
	push_obj.Data  = data
	push_obj.Type  = "push"

	return push_obj
}


func (this *Pack) WrapBroatcastRespObj(channel_id string, from_sid string , data interface{}) BroatcastRoot {

	broatcast_header_obj := BroatcastHeader{}
	broatcast_header_obj.Sid = from_sid
	broatcast_header_obj.ChannelId = channel_id

	broatcast_obj := BroatcastRoot{}
	broatcast_obj.Header =broatcast_header_obj
	broatcast_obj.Data  = data
	broatcast_obj.Type  = "broatcast"

	return broatcast_obj
}

func (this *Pack) GetBroatcastHeaderObj(  header []byte) (*RespHeader, error) {

	stb := &BroatcastHeader{}
	err := json.Unmarshal( header, stb )
	return stb, err
}
