// automatically generated, do not modify

package protocol

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Data struct {
	_tab flatbuffers.Table
}

func GetRootAsData(buf []byte, offset flatbuffers.UOffsetT) *Data {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Data{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Data) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Data) _type() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 1
}

func (rcv *Data) Cmd() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Data) Sid() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *Data) ReqId() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Data) Data() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func DataStart(builder *flatbuffers.Builder) { builder.StartObject(5) }
func DataAdd_type(builder *flatbuffers.Builder, Type int16) { builder.PrependInt16Slot(0, Type, 1) }
func DataAddCmd(builder *flatbuffers.Builder, cmd flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(cmd), 0) }
func DataAddSid(builder *flatbuffers.Builder, sid flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(2, flatbuffers.UOffsetT(sid), 0) }
func DataAddReqId(builder *flatbuffers.Builder, reqId int32) { builder.PrependInt32Slot(3, reqId, 0) }
func DataAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(4, flatbuffers.UOffsetT(data), 0) }
func DataEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
