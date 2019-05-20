// Code generated by protoc-gen-go. DO NOT EDIT.
// source: log_propagation.proto

package antidote // import "github.com/dvasilas/proteus/protos/antidote"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type SubRequest struct {
	Timestamp            int64    `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubRequest) Reset()         { *m = SubRequest{} }
func (m *SubRequest) String() string { return proto.CompactTextString(m) }
func (*SubRequest) ProtoMessage()    {}
func (*SubRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{0}
}
func (m *SubRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubRequest.Unmarshal(m, b)
}
func (m *SubRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubRequest.Marshal(b, m, deterministic)
}
func (dst *SubRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubRequest.Merge(dst, src)
}
func (m *SubRequest) XXX_Size() int {
	return xxx_messageInfo_SubRequest.Size(m)
}
func (m *SubRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubRequest proto.InternalMessageInfo

func (m *SubRequest) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type LogOperation struct {
	Dc_ID                string                `protobuf:"bytes,1,opt,name=dc_ID,json=dcID" json:"dc_ID,omitempty"`
	Partition_ID         string                `protobuf:"bytes,2,opt,name=partition_ID,json=partitionID" json:"partition_ID,omitempty"`
	Key                  string                `protobuf:"bytes,3,opt,name=key" json:"key,omitempty"`
	Bucket               string                `protobuf:"bytes,4,opt,name=bucket" json:"bucket,omitempty"`
	CrdtType             string                `protobuf:"bytes,5,opt,name=crdt_type,json=crdtType" json:"crdt_type,omitempty"`
	CommitTime           int64                 `protobuf:"varint,6,opt,name=commit_time,json=commitTime" json:"commit_time,omitempty"`
	Payload              *LogOperation_Payload `protobuf:"bytes,7,opt,name=payload" json:"payload,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *LogOperation) Reset()         { *m = LogOperation{} }
func (m *LogOperation) String() string { return proto.CompactTextString(m) }
func (*LogOperation) ProtoMessage()    {}
func (*LogOperation) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{1}
}
func (m *LogOperation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogOperation.Unmarshal(m, b)
}
func (m *LogOperation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogOperation.Marshal(b, m, deterministic)
}
func (dst *LogOperation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogOperation.Merge(dst, src)
}
func (m *LogOperation) XXX_Size() int {
	return xxx_messageInfo_LogOperation.Size(m)
}
func (m *LogOperation) XXX_DiscardUnknown() {
	xxx_messageInfo_LogOperation.DiscardUnknown(m)
}

var xxx_messageInfo_LogOperation proto.InternalMessageInfo

func (m *LogOperation) GetDc_ID() string {
	if m != nil {
		return m.Dc_ID
	}
	return ""
}

func (m *LogOperation) GetPartition_ID() string {
	if m != nil {
		return m.Partition_ID
	}
	return ""
}

func (m *LogOperation) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *LogOperation) GetBucket() string {
	if m != nil {
		return m.Bucket
	}
	return ""
}

func (m *LogOperation) GetCrdtType() string {
	if m != nil {
		return m.CrdtType
	}
	return ""
}

func (m *LogOperation) GetCommitTime() int64 {
	if m != nil {
		return m.CommitTime
	}
	return 0
}

func (m *LogOperation) GetPayload() *LogOperation_Payload {
	if m != nil {
		return m.Payload
	}
	return nil
}

type LogOperation_StateDelta struct {
	Old                  *CrdtMapState `protobuf:"bytes,1,opt,name=old" json:"old,omitempty"`
	New                  *CrdtMapState `protobuf:"bytes,2,opt,name=new" json:"new,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *LogOperation_StateDelta) Reset()         { *m = LogOperation_StateDelta{} }
func (m *LogOperation_StateDelta) String() string { return proto.CompactTextString(m) }
func (*LogOperation_StateDelta) ProtoMessage()    {}
func (*LogOperation_StateDelta) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{1, 0}
}
func (m *LogOperation_StateDelta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogOperation_StateDelta.Unmarshal(m, b)
}
func (m *LogOperation_StateDelta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogOperation_StateDelta.Marshal(b, m, deterministic)
}
func (dst *LogOperation_StateDelta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogOperation_StateDelta.Merge(dst, src)
}
func (m *LogOperation_StateDelta) XXX_Size() int {
	return xxx_messageInfo_LogOperation_StateDelta.Size(m)
}
func (m *LogOperation_StateDelta) XXX_DiscardUnknown() {
	xxx_messageInfo_LogOperation_StateDelta.DiscardUnknown(m)
}

var xxx_messageInfo_LogOperation_StateDelta proto.InternalMessageInfo

func (m *LogOperation_StateDelta) GetOld() *CrdtMapState {
	if m != nil {
		return m.Old
	}
	return nil
}

func (m *LogOperation_StateDelta) GetNew() *CrdtMapState {
	if m != nil {
		return m.New
	}
	return nil
}

type LogOperation_Payload struct {
	// Types that are valid to be assigned to Val:
	//	*LogOperation_Payload_Delta
	//	*LogOperation_Payload_Op
	Val                  isLogOperation_Payload_Val `protobuf_oneof:"val"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *LogOperation_Payload) Reset()         { *m = LogOperation_Payload{} }
func (m *LogOperation_Payload) String() string { return proto.CompactTextString(m) }
func (*LogOperation_Payload) ProtoMessage()    {}
func (*LogOperation_Payload) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{1, 1}
}
func (m *LogOperation_Payload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogOperation_Payload.Unmarshal(m, b)
}
func (m *LogOperation_Payload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogOperation_Payload.Marshal(b, m, deterministic)
}
func (dst *LogOperation_Payload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogOperation_Payload.Merge(dst, src)
}
func (m *LogOperation_Payload) XXX_Size() int {
	return xxx_messageInfo_LogOperation_Payload.Size(m)
}
func (m *LogOperation_Payload) XXX_DiscardUnknown() {
	xxx_messageInfo_LogOperation_Payload.DiscardUnknown(m)
}

var xxx_messageInfo_LogOperation_Payload proto.InternalMessageInfo

type isLogOperation_Payload_Val interface {
	isLogOperation_Payload_Val()
}

type LogOperation_Payload_Delta struct {
	Delta *LogOperation_StateDelta `protobuf:"bytes,1,opt,name=delta,oneof"`
}
type LogOperation_Payload_Op struct {
	Op *Operation `protobuf:"bytes,2,opt,name=op,oneof"`
}

func (*LogOperation_Payload_Delta) isLogOperation_Payload_Val() {}
func (*LogOperation_Payload_Op) isLogOperation_Payload_Val()    {}

func (m *LogOperation_Payload) GetVal() isLogOperation_Payload_Val {
	if m != nil {
		return m.Val
	}
	return nil
}

func (m *LogOperation_Payload) GetDelta() *LogOperation_StateDelta {
	if x, ok := m.GetVal().(*LogOperation_Payload_Delta); ok {
		return x.Delta
	}
	return nil
}

func (m *LogOperation_Payload) GetOp() *Operation {
	if x, ok := m.GetVal().(*LogOperation_Payload_Op); ok {
		return x.Op
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*LogOperation_Payload) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _LogOperation_Payload_OneofMarshaler, _LogOperation_Payload_OneofUnmarshaler, _LogOperation_Payload_OneofSizer, []interface{}{
		(*LogOperation_Payload_Delta)(nil),
		(*LogOperation_Payload_Op)(nil),
	}
}

func _LogOperation_Payload_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*LogOperation_Payload)
	// val
	switch x := m.Val.(type) {
	case *LogOperation_Payload_Delta:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Delta); err != nil {
			return err
		}
	case *LogOperation_Payload_Op:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Op); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("LogOperation_Payload.Val has unexpected type %T", x)
	}
	return nil
}

func _LogOperation_Payload_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*LogOperation_Payload)
	switch tag {
	case 1: // val.delta
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(LogOperation_StateDelta)
		err := b.DecodeMessage(msg)
		m.Val = &LogOperation_Payload_Delta{msg}
		return true, err
	case 2: // val.op
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Operation)
		err := b.DecodeMessage(msg)
		m.Val = &LogOperation_Payload_Op{msg}
		return true, err
	default:
		return false, nil
	}
}

func _LogOperation_Payload_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*LogOperation_Payload)
	// val
	switch x := m.Val.(type) {
	case *LogOperation_Payload_Delta:
		s := proto.Size(x.Delta)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *LogOperation_Payload_Op:
		s := proto.Size(x.Op)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type CrdtMapState struct {
	State                []*CrdtMapState_MapState `protobuf:"bytes,1,rep,name=state" json:"state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *CrdtMapState) Reset()         { *m = CrdtMapState{} }
func (m *CrdtMapState) String() string { return proto.CompactTextString(m) }
func (*CrdtMapState) ProtoMessage()    {}
func (*CrdtMapState) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{2}
}
func (m *CrdtMapState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrdtMapState.Unmarshal(m, b)
}
func (m *CrdtMapState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrdtMapState.Marshal(b, m, deterministic)
}
func (dst *CrdtMapState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrdtMapState.Merge(dst, src)
}
func (m *CrdtMapState) XXX_Size() int {
	return xxx_messageInfo_CrdtMapState.Size(m)
}
func (m *CrdtMapState) XXX_DiscardUnknown() {
	xxx_messageInfo_CrdtMapState.DiscardUnknown(m)
}

var xxx_messageInfo_CrdtMapState proto.InternalMessageInfo

func (m *CrdtMapState) GetState() []*CrdtMapState_MapState {
	if m != nil {
		return m.State
	}
	return nil
}

type CrdtMapState_MapState struct {
	Object               *CrdtKeyType `protobuf:"bytes,1,opt,name=object" json:"object,omitempty"`
	Value                *CrdtValue   `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *CrdtMapState_MapState) Reset()         { *m = CrdtMapState_MapState{} }
func (m *CrdtMapState_MapState) String() string { return proto.CompactTextString(m) }
func (*CrdtMapState_MapState) ProtoMessage()    {}
func (*CrdtMapState_MapState) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{2, 0}
}
func (m *CrdtMapState_MapState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrdtMapState_MapState.Unmarshal(m, b)
}
func (m *CrdtMapState_MapState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrdtMapState_MapState.Marshal(b, m, deterministic)
}
func (dst *CrdtMapState_MapState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrdtMapState_MapState.Merge(dst, src)
}
func (m *CrdtMapState_MapState) XXX_Size() int {
	return xxx_messageInfo_CrdtMapState_MapState.Size(m)
}
func (m *CrdtMapState_MapState) XXX_DiscardUnknown() {
	xxx_messageInfo_CrdtMapState_MapState.DiscardUnknown(m)
}

var xxx_messageInfo_CrdtMapState_MapState proto.InternalMessageInfo

func (m *CrdtMapState_MapState) GetObject() *CrdtKeyType {
	if m != nil {
		return m.Object
	}
	return nil
}

func (m *CrdtMapState_MapState) GetValue() *CrdtValue {
	if m != nil {
		return m.Value
	}
	return nil
}

type Operation struct {
	Op                   []*Operation_Op `protobuf:"bytes,1,rep,name=op" json:"op,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *Operation) Reset()         { *m = Operation{} }
func (m *Operation) String() string { return proto.CompactTextString(m) }
func (*Operation) ProtoMessage()    {}
func (*Operation) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{3}
}
func (m *Operation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Operation.Unmarshal(m, b)
}
func (m *Operation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Operation.Marshal(b, m, deterministic)
}
func (dst *Operation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Operation.Merge(dst, src)
}
func (m *Operation) XXX_Size() int {
	return xxx_messageInfo_Operation.Size(m)
}
func (m *Operation) XXX_DiscardUnknown() {
	xxx_messageInfo_Operation.DiscardUnknown(m)
}

var xxx_messageInfo_Operation proto.InternalMessageInfo

func (m *Operation) GetOp() []*Operation_Op {
	if m != nil {
		return m.Op
	}
	return nil
}

type Operation_Op struct {
	Object               *CrdtKeyType `protobuf:"bytes,1,opt,name=object" json:"object,omitempty"`
	Update               *Update      `protobuf:"bytes,2,opt,name=update" json:"update,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Operation_Op) Reset()         { *m = Operation_Op{} }
func (m *Operation_Op) String() string { return proto.CompactTextString(m) }
func (*Operation_Op) ProtoMessage()    {}
func (*Operation_Op) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{3, 0}
}
func (m *Operation_Op) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Operation_Op.Unmarshal(m, b)
}
func (m *Operation_Op) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Operation_Op.Marshal(b, m, deterministic)
}
func (dst *Operation_Op) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Operation_Op.Merge(dst, src)
}
func (m *Operation_Op) XXX_Size() int {
	return xxx_messageInfo_Operation_Op.Size(m)
}
func (m *Operation_Op) XXX_DiscardUnknown() {
	xxx_messageInfo_Operation_Op.DiscardUnknown(m)
}

var xxx_messageInfo_Operation_Op proto.InternalMessageInfo

func (m *Operation_Op) GetObject() *CrdtKeyType {
	if m != nil {
		return m.Object
	}
	return nil
}

func (m *Operation_Op) GetUpdate() *Update {
	if m != nil {
		return m.Update
	}
	return nil
}

type CrdtKeyType struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=type" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CrdtKeyType) Reset()         { *m = CrdtKeyType{} }
func (m *CrdtKeyType) String() string { return proto.CompactTextString(m) }
func (*CrdtKeyType) ProtoMessage()    {}
func (*CrdtKeyType) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{4}
}
func (m *CrdtKeyType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrdtKeyType.Unmarshal(m, b)
}
func (m *CrdtKeyType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrdtKeyType.Marshal(b, m, deterministic)
}
func (dst *CrdtKeyType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrdtKeyType.Merge(dst, src)
}
func (m *CrdtKeyType) XXX_Size() int {
	return xxx_messageInfo_CrdtKeyType.Size(m)
}
func (m *CrdtKeyType) XXX_DiscardUnknown() {
	xxx_messageInfo_CrdtKeyType.DiscardUnknown(m)
}

var xxx_messageInfo_CrdtKeyType proto.InternalMessageInfo

func (m *CrdtKeyType) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *CrdtKeyType) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type CrdtValue struct {
	// Types that are valid to be assigned to Val:
	//	*CrdtValue_Str
	//	*CrdtValue_Int
	Val                  isCrdtValue_Val `protobuf_oneof:"val"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *CrdtValue) Reset()         { *m = CrdtValue{} }
func (m *CrdtValue) String() string { return proto.CompactTextString(m) }
func (*CrdtValue) ProtoMessage()    {}
func (*CrdtValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{5}
}
func (m *CrdtValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CrdtValue.Unmarshal(m, b)
}
func (m *CrdtValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CrdtValue.Marshal(b, m, deterministic)
}
func (dst *CrdtValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CrdtValue.Merge(dst, src)
}
func (m *CrdtValue) XXX_Size() int {
	return xxx_messageInfo_CrdtValue.Size(m)
}
func (m *CrdtValue) XXX_DiscardUnknown() {
	xxx_messageInfo_CrdtValue.DiscardUnknown(m)
}

var xxx_messageInfo_CrdtValue proto.InternalMessageInfo

type isCrdtValue_Val interface {
	isCrdtValue_Val()
}

type CrdtValue_Str struct {
	Str string `protobuf:"bytes,1,opt,name=str,oneof"`
}
type CrdtValue_Int struct {
	Int int64 `protobuf:"varint,2,opt,name=int,oneof"`
}

func (*CrdtValue_Str) isCrdtValue_Val() {}
func (*CrdtValue_Int) isCrdtValue_Val() {}

func (m *CrdtValue) GetVal() isCrdtValue_Val {
	if m != nil {
		return m.Val
	}
	return nil
}

func (m *CrdtValue) GetStr() string {
	if x, ok := m.GetVal().(*CrdtValue_Str); ok {
		return x.Str
	}
	return ""
}

func (m *CrdtValue) GetInt() int64 {
	if x, ok := m.GetVal().(*CrdtValue_Int); ok {
		return x.Int
	}
	return 0
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*CrdtValue) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CrdtValue_OneofMarshaler, _CrdtValue_OneofUnmarshaler, _CrdtValue_OneofSizer, []interface{}{
		(*CrdtValue_Str)(nil),
		(*CrdtValue_Int)(nil),
	}
}

func _CrdtValue_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*CrdtValue)
	// val
	switch x := m.Val.(type) {
	case *CrdtValue_Str:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		b.EncodeStringBytes(x.Str)
	case *CrdtValue_Int:
		b.EncodeVarint(2<<3 | proto.WireVarint)
		b.EncodeVarint(uint64(x.Int))
	case nil:
	default:
		return fmt.Errorf("CrdtValue.Val has unexpected type %T", x)
	}
	return nil
}

func _CrdtValue_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*CrdtValue)
	switch tag {
	case 1: // val.str
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeStringBytes()
		m.Val = &CrdtValue_Str{x}
		return true, err
	case 2: // val.int
		if wire != proto.WireVarint {
			return true, proto.ErrInternalBadWireType
		}
		x, err := b.DecodeVarint()
		m.Val = &CrdtValue_Int{int64(x)}
		return true, err
	default:
		return false, nil
	}
}

func _CrdtValue_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*CrdtValue)
	// val
	switch x := m.Val.(type) {
	case *CrdtValue_Str:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(len(x.Str)))
		n += len(x.Str)
	case *CrdtValue_Int:
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(x.Int))
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
}

type Update struct {
	OpType               string     `protobuf:"bytes,1,opt,name=op_type,json=opType" json:"op_type,omitempty"`
	Value                *CrdtValue `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Update) Reset()         { *m = Update{} }
func (m *Update) String() string { return proto.CompactTextString(m) }
func (*Update) ProtoMessage()    {}
func (*Update) Descriptor() ([]byte, []int) {
	return fileDescriptor_log_propagation_6d7d2d2a93f2854a, []int{6}
}
func (m *Update) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Update.Unmarshal(m, b)
}
func (m *Update) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Update.Marshal(b, m, deterministic)
}
func (dst *Update) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Update.Merge(dst, src)
}
func (m *Update) XXX_Size() int {
	return xxx_messageInfo_Update.Size(m)
}
func (m *Update) XXX_DiscardUnknown() {
	xxx_messageInfo_Update.DiscardUnknown(m)
}

var xxx_messageInfo_Update proto.InternalMessageInfo

func (m *Update) GetOpType() string {
	if m != nil {
		return m.OpType
	}
	return ""
}

func (m *Update) GetValue() *CrdtValue {
	if m != nil {
		return m.Value
	}
	return nil
}

func init() {
	proto.RegisterType((*SubRequest)(nil), "logpropagation.SubRequest")
	proto.RegisterType((*LogOperation)(nil), "logpropagation.LogOperation")
	proto.RegisterType((*LogOperation_StateDelta)(nil), "logpropagation.LogOperation.StateDelta")
	proto.RegisterType((*LogOperation_Payload)(nil), "logpropagation.LogOperation.Payload")
	proto.RegisterType((*CrdtMapState)(nil), "logpropagation.CrdtMapState")
	proto.RegisterType((*CrdtMapState_MapState)(nil), "logpropagation.CrdtMapState.MapState")
	proto.RegisterType((*Operation)(nil), "logpropagation.Operation")
	proto.RegisterType((*Operation_Op)(nil), "logpropagation.Operation.Op")
	proto.RegisterType((*CrdtKeyType)(nil), "logpropagation.CrdtKeyType")
	proto.RegisterType((*CrdtValue)(nil), "logpropagation.CrdtValue")
	proto.RegisterType((*Update)(nil), "logpropagation.Update")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Service service

type ServiceClient interface {
	WatchAsync(ctx context.Context, in *SubRequest, opts ...grpc.CallOption) (Service_WatchAsyncClient, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) WatchAsync(ctx context.Context, in *SubRequest, opts ...grpc.CallOption) (Service_WatchAsyncClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Service_serviceDesc.Streams[0], c.cc, "/logpropagation.Service/WatchAsync", opts...)
	if err != nil {
		return nil, err
	}
	x := &serviceWatchAsyncClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Service_WatchAsyncClient interface {
	Recv() (*LogOperation, error)
	grpc.ClientStream
}

type serviceWatchAsyncClient struct {
	grpc.ClientStream
}

func (x *serviceWatchAsyncClient) Recv() (*LogOperation, error) {
	m := new(LogOperation)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for Service service

type ServiceServer interface {
	WatchAsync(*SubRequest, Service_WatchAsyncServer) error
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_WatchAsync_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ServiceServer).WatchAsync(m, &serviceWatchAsyncServer{stream})
}

type Service_WatchAsyncServer interface {
	Send(*LogOperation) error
	grpc.ServerStream
}

type serviceWatchAsyncServer struct {
	grpc.ServerStream
}

func (x *serviceWatchAsyncServer) Send(m *LogOperation) error {
	return x.ServerStream.SendMsg(m)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "logpropagation.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchAsync",
			Handler:       _Service_WatchAsync_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "log_propagation.proto",
}

func init() {
	proto.RegisterFile("log_propagation.proto", fileDescriptor_log_propagation_6d7d2d2a93f2854a)
}

var fileDescriptor_log_propagation_6d7d2d2a93f2854a = []byte{
	// 596 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0xdf, 0x6e, 0xd3, 0x3e,
	0x18, 0x6d, 0x9a, 0x25, 0x5d, 0xbf, 0x4c, 0x3f, 0xfd, 0x64, 0xc4, 0x08, 0x59, 0x25, 0x46, 0x04,
	0x62, 0x62, 0x90, 0xa2, 0xee, 0x0e, 0x24, 0x10, 0xa3, 0x17, 0x1b, 0x7f, 0x34, 0xe4, 0x6d, 0x20,
	0x71, 0x53, 0xb9, 0x89, 0xd5, 0x99, 0x25, 0xb1, 0x49, 0x9c, 0x4e, 0xb9, 0xe1, 0x69, 0x78, 0x0d,
	0xde, 0x84, 0x87, 0x41, 0xb6, 0xd3, 0x3f, 0x54, 0xac, 0x08, 0xae, 0xe2, 0x9c, 0xef, 0x1c, 0xfb,
	0x9c, 0x2f, 0xfe, 0x02, 0x37, 0x53, 0x3e, 0x19, 0x89, 0x82, 0x0b, 0x32, 0x21, 0x92, 0xf1, 0x3c,
	0x12, 0x05, 0x97, 0x1c, 0xfd, 0x97, 0xf2, 0xc9, 0x12, 0x1a, 0x3e, 0x04, 0x38, 0xad, 0xc6, 0x98,
	0x7e, 0xa9, 0x68, 0x29, 0x51, 0x0f, 0xba, 0x92, 0x65, 0xb4, 0x94, 0x24, 0x13, 0xbe, 0xb5, 0x6b,
	0xed, 0xd9, 0x78, 0x01, 0x84, 0x3f, 0x6c, 0xd8, 0x7a, 0xcb, 0x27, 0x27, 0x82, 0x16, 0x5a, 0x8c,
	0x6e, 0x80, 0x93, 0xc4, 0xa3, 0xe3, 0xa1, 0xa6, 0x76, 0xf1, 0x46, 0x12, 0x1f, 0x0f, 0xd1, 0x5d,
	0xd8, 0x12, 0xa4, 0x90, 0x4c, 0x31, 0x54, 0xad, 0xad, 0x6b, 0xde, 0x1c, 0x3b, 0x1e, 0xa2, 0xff,
	0xc1, 0xbe, 0xa4, 0xb5, 0x6f, 0xeb, 0x8a, 0x5a, 0xa2, 0x6d, 0x70, 0xc7, 0x55, 0x7c, 0x49, 0xa5,
	0xbf, 0xa1, 0xc1, 0xe6, 0x0d, 0xed, 0x40, 0x37, 0x2e, 0x12, 0x39, 0x92, 0xb5, 0xa0, 0xbe, 0xa3,
	0x4b, 0x9b, 0x0a, 0x38, 0xab, 0x05, 0x45, 0x77, 0xc0, 0x8b, 0x79, 0x96, 0x31, 0x39, 0x52, 0x1e,
	0x7d, 0x57, 0xfb, 0x05, 0x03, 0x9d, 0xb1, 0x8c, 0xa2, 0xe7, 0xd0, 0x11, 0xa4, 0x4e, 0x39, 0x49,
	0xfc, 0xce, 0xae, 0xb5, 0xe7, 0x0d, 0xee, 0x45, 0xbf, 0xc6, 0x8f, 0x96, 0xe3, 0x44, 0xef, 0x0d,
	0x17, 0xcf, 0x44, 0x41, 0x0a, 0x70, 0x2a, 0x89, 0xa4, 0x43, 0x9a, 0x4a, 0x82, 0x22, 0xb0, 0x79,
	0x9a, 0xe8, 0xac, 0xde, 0xa0, 0xb7, 0xba, 0xd3, 0xab, 0x22, 0x91, 0xef, 0x88, 0xd0, 0x7c, 0xac,
	0x88, 0x8a, 0x9f, 0xd3, 0x2b, 0x9d, 0xff, 0x8f, 0xfc, 0x9c, 0x5e, 0x05, 0x5f, 0xa1, 0xd3, 0x38,
	0x40, 0x2f, 0xc0, 0x49, 0xd4, 0x99, 0xcd, 0x61, 0x0f, 0xd6, 0xda, 0x5e, 0x58, 0x3c, 0x6a, 0x61,
	0xa3, 0x43, 0xfb, 0xd0, 0xe6, 0xa2, 0x39, 0xfa, 0xf6, 0xaa, 0x7a, 0x2e, 0x3d, 0x6a, 0xe1, 0x36,
	0x17, 0x87, 0x0e, 0xd8, 0x53, 0x92, 0x86, 0xdf, 0x2d, 0xd8, 0x5a, 0x76, 0x85, 0x9e, 0x81, 0x53,
	0xaa, 0x85, 0x6f, 0xed, 0xda, 0x7b, 0xde, 0xe0, 0xfe, 0xba, 0x08, 0xd1, 0x3c, 0x8b, 0xd1, 0x04,
	0x02, 0x36, 0xe7, 0x1b, 0x1d, 0x80, 0xcb, 0xc7, 0x9f, 0x69, 0x2c, 0x9b, 0x3c, 0x3b, 0xbf, 0xdb,
	0xe9, 0x0d, 0xad, 0xd5, 0x57, 0xc5, 0x0d, 0x15, 0xf5, 0xc1, 0x99, 0x92, 0xb4, 0xa2, 0xd7, 0xa5,
	0x50, 0x9a, 0x0f, 0x8a, 0x80, 0x0d, 0x2f, 0xfc, 0x66, 0x41, 0x77, 0x71, 0x37, 0x1f, 0xe9, 0x0e,
	0x18, 0xe7, 0xbd, 0x6b, 0x3b, 0x10, 0x9d, 0x08, 0xd5, 0x82, 0x80, 0x41, 0xfb, 0x44, 0xfc, 0x9b,
	0xcf, 0x08, 0xdc, 0x4a, 0x24, 0xaa, 0x4d, 0xc6, 0xe8, 0xf6, 0xaa, 0xe8, 0x5c, 0x57, 0x71, 0xc3,
	0x0a, 0x0f, 0xc0, 0x5b, 0xda, 0x66, 0x36, 0x0b, 0xd6, 0x62, 0x16, 0x10, 0x6c, 0xe8, 0xeb, 0x6e,
	0x06, 0x47, 0xaf, 0xc3, 0xa7, 0xd0, 0x9d, 0xe7, 0x45, 0x08, 0xec, 0x52, 0x16, 0x46, 0x72, 0xd4,
	0xc2, 0xea, 0x45, 0x61, 0x2c, 0x97, 0x5a, 0x63, 0x2b, 0x8c, 0xe5, 0x72, 0xf6, 0x5d, 0x31, 0xb8,
	0xc6, 0x02, 0xba, 0x05, 0x1d, 0x2e, 0xcc, 0x2c, 0x99, 0xf3, 0x5c, 0x2e, 0xb4, 0x89, 0xbf, 0xed,
	0xf5, 0xe0, 0x1c, 0x3a, 0xa7, 0xb4, 0x98, 0xb2, 0x98, 0xa2, 0xd7, 0x00, 0x1f, 0x89, 0x8c, 0x2f,
	0x5e, 0x96, 0x75, 0x1e, 0xa3, 0x60, 0x55, 0xba, 0xf8, 0xbb, 0x04, 0xbd, 0x75, 0xd7, 0x38, 0x6c,
	0x3d, 0xb1, 0x0e, 0x1f, 0x7f, 0xda, 0x9f, 0x30, 0x79, 0x51, 0x8d, 0xa3, 0x98, 0x67, 0xfd, 0x64,
	0x4a, 0x4a, 0x96, 0x92, 0xb2, 0xaf, 0x7e, 0x5d, 0xb4, 0x32, 0x4f, 0x5e, 0xf6, 0x49, 0x2e, 0x59,
	0xc2, 0x25, 0x1d, 0xbb, 0x1a, 0x38, 0xf8, 0x19, 0x00, 0x00, 0xff, 0xff, 0xcd, 0x69, 0xa5, 0x25,
	0xec, 0x04, 0x00, 0x00,
}
