//通讯协议处理，主要处理封包和解包的过程
package protocol


import (

	"encoding/binary"
	"errors"
	"fmt"
	"bufio"
	"bytes"
)

var RPC_MAGIC = [4]byte{'u', 'e', 'l', 'i'}

//|总长度4|头长度4|header|data|checksum
type ClientPacket struct {
	TotalSize uint32    //4
	HeaderSize uint32   //4
	Header   []byte
	Payload   []byte
	Checksum  uint32    //4
}



func EncodePacket( header []byte,  payload []byte) ( []byte ,error ){
	// len(totaol)+ len(header) + len(Checksum) == 12
	var pkg *bytes.Buffer = new(bytes.Buffer)
	totalsize := uint32(len(string(header)) +  len(string(payload)) )
	err:=binary.Write( pkg , binary.LittleEndian, totalsize)
	if err != nil {
		return nil, err
	}

	headersize := uint32(len(string(header)))
	err = binary.Write( pkg, binary.LittleEndian, headersize)
	if err != nil {
		return nil, err
	}


	err = binary.Write( pkg, binary.LittleEndian, header)
	if err != nil {
		return nil, err
	}
	err = binary.Write(pkg, binary.LittleEndian, payload )
	if err != nil {
		return nil, err
	}

	// write checksum
	return  pkg.Bytes(),nil
}

func Packet2(buf []byte) ([]byte, error) {

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

func Unpack2(reader *bufio.Reader) ([]byte, error) {

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




func DecodePacket(r *bufio.Reader) (  []byte,   []byte, error) {
	var totalsize , headersize uint32

	lengthByte, _ := r.Peek(4)
	lengthBuff := bytes.NewBuffer(lengthByte)
	err := binary.Read(lengthBuff, binary.LittleEndian, &totalsize)
	if err != nil {
		return nil,nil,errors.New("read total size error") // errors.Annotate(err, "read total size")
	}
	fmt.Println( "totalsize" , totalsize)

	hlengthByte, _ := r.Peek(8)
	hlengthBuff := bytes.NewBuffer(hlengthByte)
	binary.Read(hlengthBuff, binary.LittleEndian, &headersize)

	fmt.Println( "headersize" , headersize)

	// at least len(magic) + len(checksum)
	if totalsize < 12 {
		return nil, nil,errors.New( fmt.Sprintf("bad packet. header:%d", totalsize))
	}


	pack := make([]byte, int(totalsize))
	_, err = r.Read(pack)
	if err != nil {
		return nil,nil, err
	}
	header := pack[:headersize]

	payload := pack[headersize:totalsize]

	fmt.Println( "header:" , string(header))
	fmt.Println( "payload:" ,  string(payload) )
	return header,payload, nil
}