//通讯协议处理，主要处理封包和解包的过程
package protocol


import (

	"encoding/binary"
	"errors"
	"fmt"
	"bufio"
	"bytes"
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

