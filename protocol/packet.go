//通讯协议处理，主要处理封包和解包的过程
package protocol
 
 
import (
	"bufio"
	"bytes"
	"encoding/binary"

)


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


