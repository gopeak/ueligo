// automatically generated, do not modify

package protocol

import (
	"strings"
	"strconv"
	"errors"
)


// 使用字符串分割时数据块定义
const (
	MSG_TYPE_INDEX    = 0
	MSG_CMD_INDEX     = 1
	MSG_SID_INDEX =  2
	MSG_REQID_INDEX   = 3
	MSG_DATA_INDEX   = 4
)

// 数据类型
const (
	TypeReq = 1
	TypeReply = 2
	TypePush = 3
	TypeBroadcast = 4
	TypeError = 4
)

func ParseData(str string) ( error ,int ,string,string,int,string){

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

	return err,_type,cmd,req_sid,req_id,req_data
}