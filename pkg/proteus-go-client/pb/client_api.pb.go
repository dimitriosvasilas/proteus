// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/protobuf-spec/client_api.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type QueryReq struct {
	QueryStr             string   `protobuf:"bytes,1,opt,name=queryStr,proto3" json:"queryStr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryReq) Reset()         { *m = QueryReq{} }
func (m *QueryReq) String() string { return proto.CompactTextString(m) }
func (*QueryReq) ProtoMessage()    {}
func (*QueryReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_a58988a9bf4976bc, []int{0}
}

func (m *QueryReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryReq.Unmarshal(m, b)
}
func (m *QueryReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryReq.Marshal(b, m, deterministic)
}
func (m *QueryReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryReq.Merge(m, src)
}
func (m *QueryReq) XXX_Size() int {
	return xxx_messageInfo_QueryReq.Size(m)
}
func (m *QueryReq) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryReq.DiscardUnknown(m)
}

var xxx_messageInfo_QueryReq proto.InternalMessageInfo

func (m *QueryReq) GetQueryStr() string {
	if m != nil {
		return m.QueryStr
	}
	return ""
}

type QueryResp struct {
	RespRecord           []*QueryRespRecord `protobuf:"bytes,1,rep,name=respRecord,proto3" json:"respRecord,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *QueryResp) Reset()         { *m = QueryResp{} }
func (m *QueryResp) String() string { return proto.CompactTextString(m) }
func (*QueryResp) ProtoMessage()    {}
func (*QueryResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_a58988a9bf4976bc, []int{1}
}

func (m *QueryResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryResp.Unmarshal(m, b)
}
func (m *QueryResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryResp.Marshal(b, m, deterministic)
}
func (m *QueryResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryResp.Merge(m, src)
}
func (m *QueryResp) XXX_Size() int {
	return xxx_messageInfo_QueryResp.Size(m)
}
func (m *QueryResp) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryResp.DiscardUnknown(m)
}

var xxx_messageInfo_QueryResp proto.InternalMessageInfo

func (m *QueryResp) GetRespRecord() []*QueryRespRecord {
	if m != nil {
		return m.RespRecord
	}
	return nil
}

type QueryRespRecord struct {
	RecordId             string                          `protobuf:"bytes,1,opt,name=record_id,json=recordId,proto3" json:"record_id,omitempty"`
	Attributes           map[string][]byte               `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Timestamp            map[string]*timestamp.Timestamp `protobuf:"bytes,3,rep,name=timestamp,proto3" json:"timestamp,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *QueryRespRecord) Reset()         { *m = QueryRespRecord{} }
func (m *QueryRespRecord) String() string { return proto.CompactTextString(m) }
func (*QueryRespRecord) ProtoMessage()    {}
func (*QueryRespRecord) Descriptor() ([]byte, []int) {
	return fileDescriptor_a58988a9bf4976bc, []int{2}
}

func (m *QueryRespRecord) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryRespRecord.Unmarshal(m, b)
}
func (m *QueryRespRecord) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryRespRecord.Marshal(b, m, deterministic)
}
func (m *QueryRespRecord) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRespRecord.Merge(m, src)
}
func (m *QueryRespRecord) XXX_Size() int {
	return xxx_messageInfo_QueryRespRecord.Size(m)
}
func (m *QueryRespRecord) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRespRecord.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRespRecord proto.InternalMessageInfo

func (m *QueryRespRecord) GetRecordId() string {
	if m != nil {
		return m.RecordId
	}
	return ""
}

func (m *QueryRespRecord) GetAttributes() map[string][]byte {
	if m != nil {
		return m.Attributes
	}
	return nil
}

func (m *QueryRespRecord) GetTimestamp() map[string]*timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func init() {
	proto.RegisterType((*QueryReq)(nil), "pb.QueryReq")
	proto.RegisterType((*QueryResp)(nil), "pb.QueryResp")
	proto.RegisterType((*QueryRespRecord)(nil), "pb.QueryRespRecord")
	proto.RegisterMapType((map[string][]byte)(nil), "pb.QueryRespRecord.AttributesEntry")
	proto.RegisterMapType((map[string]*timestamp.Timestamp)(nil), "pb.QueryRespRecord.TimestampEntry")
}

func init() { proto.RegisterFile("api/protobuf-spec/client_api.proto", fileDescriptor_a58988a9bf4976bc) }

var fileDescriptor_a58988a9bf4976bc = []byte{
	// 360 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xdf, 0x4b, 0xc2, 0x40,
	0x1c, 0x6f, 0x93, 0xc4, 0x7d, 0xb5, 0x8c, 0xab, 0x87, 0xb1, 0x1e, 0x92, 0x05, 0x21, 0x84, 0x5b,
	0x68, 0x45, 0x04, 0x81, 0x16, 0x3d, 0xf8, 0xa6, 0x2b, 0x21, 0x7a, 0x91, 0xdb, 0x76, 0xad, 0xc3,
	0xe9, 0xce, 0xbb, 0x9b, 0xb0, 0x7f, 0xb2, 0xbf, 0x29, 0xdc, 0xdc, 0x9c, 0xe2, 0xdb, 0xe7, 0xee,
	0xfb, 0xf9, 0xf1, 0xfd, 0x70, 0x07, 0x26, 0x66, 0xd4, 0x66, 0x3c, 0x92, 0x91, 0x1b, 0xff, 0x74,
	0x04, 0x23, 0x9e, 0xed, 0x85, 0x94, 0x2c, 0xe4, 0x14, 0x33, 0x6a, 0xa5, 0x03, 0xa4, 0x32, 0xd7,
	0xb8, 0x0a, 0xa2, 0x28, 0x08, 0x49, 0x41, 0xb5, 0x25, 0x9d, 0x13, 0x21, 0xf1, 0x9c, 0x65, 0x24,
	0xf3, 0x06, 0x6a, 0xe3, 0x98, 0xf0, 0xc4, 0x21, 0x4b, 0x64, 0x40, 0x6d, 0xb9, 0xc6, 0x1f, 0x92,
	0xeb, 0x4a, 0x4b, 0x69, 0x6b, 0x4e, 0x71, 0x36, 0xfb, 0xa0, 0x6d, 0x78, 0x82, 0xa1, 0x1e, 0x00,
	0x27, 0x82, 0x39, 0xc4, 0x8b, 0xb8, 0xaf, 0x2b, 0xad, 0x4a, 0xbb, 0xde, 0x3d, 0xb7, 0x98, 0x6b,
	0x15, 0x94, 0x6c, 0xe4, 0x94, 0x68, 0xe6, 0x9f, 0x0a, 0xcd, 0xbd, 0x39, 0xba, 0x04, 0x8d, 0xa7,
	0x68, 0x4a, 0xfd, 0x3c, 0x32, 0xbb, 0x18, 0xfa, 0xe8, 0x0d, 0x00, 0x4b, 0xc9, 0xa9, 0x1b, 0x4b,
	0x22, 0x74, 0x35, 0x4d, 0xb9, 0x3e, 0x90, 0x62, 0x0d, 0x0a, 0xd6, 0xfb, 0x42, 0xf2, 0xc4, 0x29,
	0xc9, 0x50, 0x1f, 0xb4, 0xa2, 0xb2, 0x5e, 0x49, 0x3d, 0xcc, 0x43, 0x1e, 0x9f, 0x39, 0x29, 0xb3,
	0xd8, 0x8a, 0x8c, 0x17, 0x68, 0xee, 0x05, 0xa0, 0x33, 0xa8, 0xcc, 0x48, 0xb2, 0x59, 0x78, 0x0d,
	0xd1, 0x05, 0x1c, 0xaf, 0x70, 0x18, 0x13, 0x5d, 0x6d, 0x29, 0xed, 0x86, 0x93, 0x1d, 0x9e, 0xd5,
	0x27, 0xc5, 0xf8, 0x82, 0xd3, 0x5d, 0xef, 0x03, 0xea, 0xbb, 0xb2, 0xba, 0xde, 0x35, 0xac, 0xec,
	0xd5, 0xac, 0xfc, 0xd5, 0xb6, 0xdb, 0x95, 0x9c, 0xbb, 0x0f, 0x50, 0x1d, 0x8f, 0x26, 0x83, 0xd1,
	0x10, 0xdd, 0x02, 0xa4, 0x7d, 0x26, 0x0b, 0xcc, 0x13, 0xd4, 0x28, 0xf5, 0x5b, 0x1a, 0x27, 0x3b,
	0x6d, 0xcd, 0xa3, 0xd7, 0xc7, 0xef, 0xfb, 0x80, 0xca, 0xdf, 0xd8, 0xb5, 0xbc, 0x68, 0x6e, 0xfb,
	0x2b, 0x2c, 0x68, 0x88, 0x45, 0xfa, 0x43, 0x48, 0x2c, 0x6c, 0x36, 0x0b, 0x72, 0xdc, 0x09, 0xa2,
	0x4e, 0xf6, 0xab, 0x6c, 0xe6, 0xba, 0xd5, 0x74, 0x9b, 0xde, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x2e, 0x58, 0xc8, 0xe4, 0x7b, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QPUAPIClient is the client API for QPUAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QPUAPIClient interface {
	QueryUnary(ctx context.Context, in *QueryReq, opts ...grpc.CallOption) (*QueryResp, error)
}

type qPUAPIClient struct {
	cc *grpc.ClientConn
}

func NewQPUAPIClient(cc *grpc.ClientConn) QPUAPIClient {
	return &qPUAPIClient{cc}
}

func (c *qPUAPIClient) QueryUnary(ctx context.Context, in *QueryReq, opts ...grpc.CallOption) (*QueryResp, error) {
	out := new(QueryResp)
	err := c.cc.Invoke(ctx, "/pb.QPUAPI/QueryUnary", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QPUAPIServer is the server API for QPUAPI service.
type QPUAPIServer interface {
	QueryUnary(context.Context, *QueryReq) (*QueryResp, error)
}

func RegisterQPUAPIServer(s *grpc.Server, srv QPUAPIServer) {
	s.RegisterService(&_QPUAPI_serviceDesc, srv)
}

func _QPUAPI_QueryUnary_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QPUAPIServer).QueryUnary(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.QPUAPI/QueryUnary",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QPUAPIServer).QueryUnary(ctx, req.(*QueryReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _QPUAPI_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.QPUAPI",
	HandlerType: (*QPUAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryUnary",
			Handler:    _QPUAPI_QueryUnary_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/protobuf-spec/client_api.proto",
}