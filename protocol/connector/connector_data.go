// automatically generated, do not modify

package connector

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type connector_data struct {
	_tab flatbuffers.Table
}

func GetRootAsconnector_data(buf []byte, offset flatbuffers.UOffsetT) *connector_data {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &connector_data{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *connector_data) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *connector_data) Type() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 1
}

func (rcv *connector_data) Data() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *connector_data) Reqid() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *connector_data) Keeplive() int16 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		return rcv._tab.GetInt16(o + rcv._tab.Pos)
	}
	return 1
}

func connector_dataStart(builder *flatbuffers.Builder) { builder.StartObject(4) }
func connector_dataAddType(builder *flatbuffers.Builder, type int16) { builder.PrependInt16Slot(0, type, 1) }
func connector_dataAddData(builder *flatbuffers.Builder, data flatbuffers.UOffsetT) { builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(data), 0) }
func connector_dataAddReqid(builder *flatbuffers.Builder, reqid int32) { builder.PrependInt32Slot(2, reqid, 0) }
func connector_dataAddKeeplive(builder *flatbuffers.Builder, keeplive int16) { builder.PrependInt16Slot(3, keeplive, 1) }
func connector_dataEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT { return builder.EndObject() }
