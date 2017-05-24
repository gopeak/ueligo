// automatically generated, do not modify

package protocol

import (
	flatbuffers "github.com/google/flatbuffers/go"
	"morego/protocol/hub"
	"encoding/binary"
	"bytes"
	"bufio"
)


// main
func MakeHubReq(  cmd  , sid  , req_id   , data string) []byte {
	// re-use the already-allocated Builder:
	b := flatbuffers.NewBuilder(0)
	b.Reset()

	// create the name object and get its offset:
	cmd_position := b.CreateByteString( []byte(cmd) )
	sid_position := b.CreateByteString( []byte(sid) )
	data_position := b.CreateByteString( []byte(data) )
	req_id_position := b.CreateByteString( []byte(req_id) )

	// write the User object:
	hub.HubReqStart(b)
	hub.HubReqAddCmd( b, cmd_position )
	hub.HubReqAddSid( b,sid_position )
	hub.HubReqAddReqId(b,req_id_position)
	hub.HubReqAddData( b, data_position )
	end_position := hub.HubReqEnd( b )

	// finish the write operations by our User the root object:
	b.Finish(end_position)

	// return the byte slice containing encoded data:
	return b.Bytes[b.Head():]
}

func ReadHubReq(buf []byte) ( cmd  , sid   , req_id string , data []byte) {
	// initialize a hub_req reader from the given buffer:
	hub_req := hub.GetRootAsHubReq(buf, 0)

	cmd = string(hub_req.Cmd())
	sid = string(hub_req.Sid())
	req_id = string(hub_req.ReqId())
	data = hub_req.Data()

	return
}



// main
func MakeHubResp(  cmd  ,   req_id   ,err  ,data string) []byte {
	// re-use the already-allocated Builder:
	b := flatbuffers.NewBuilder(0)
	b.Reset()

	// create the name object and get its offset:
	cmd_position := b.CreateByteString( []byte(cmd) )
	err_position := b.CreateByteString( []byte(err) )
	data_position := b.CreateByteString( []byte(data) )
	reqid_position := b.CreateByteString( []byte(req_id) )

	// write the User object:
	hub.HubRespStart(b)
	hub.HubRespAddCmd( b, cmd_position )
	hub.HubRespAddReqId(b,reqid_position)
	hub.HubRespAddErr( b,err_position )
	hub.HubRespAddData( b, data_position )
	end_position := hub.HubRespEnd( b )
	b.Finish(end_position)

	// return the byte slice containing encoded data:
	data_buf,_ := Packet(b.Bytes[b.Head():])
	return data_buf
}

func ReadHubResp(buf []byte) ( cmd  , req_id   ,err string ,  data []byte) {
	// initialize a hub_req reader from the given buffer:

	hub_resp := hub.GetRootAsHubResp(buf, 0)

	cmd = string(hub_resp.Cmd())
	err = string(hub_resp.Err())
	req_id = string(hub_resp.ReqId())
	data = hub_resp.Data()

	return
}


func Packet(buf []byte) ([]byte, error) {

	var length int32 = int32(len( string(buf) ))
	// fmt.Println( "Set length :", length )
	var pkg *bytes.Buffer = new(bytes.Buffer)
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}

	err = binary.Write(pkg, binary.LittleEndian, buf )
	if err != nil {
		return nil, err
	}

	return pkg.Bytes(), nil

}



func Unpack(reader *bufio.Reader) ([]byte, error) {

	lengthByte, _ := reader.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	var length int32
	err := binary.Read(lengthBuff, binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	// fmt.Println( "Get length :", length )
	if int32(reader.Buffered()) < length+4 {
		return nil, err
	}

	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return nil, err
	}
	return  pack[4:] , nil
}



