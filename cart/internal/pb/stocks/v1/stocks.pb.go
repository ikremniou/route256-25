// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: stocks/v1/stocks.proto

package stocks_v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StocksInfoRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Sku           int64                  `protobuf:"varint,1,opt,name=sku,proto3" json:"sku,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StocksInfoRequest) Reset() {
	*x = StocksInfoRequest{}
	mi := &file_stocks_v1_stocks_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StocksInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StocksInfoRequest) ProtoMessage() {}

func (x *StocksInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_stocks_v1_stocks_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StocksInfoRequest.ProtoReflect.Descriptor instead.
func (*StocksInfoRequest) Descriptor() ([]byte, []int) {
	return file_stocks_v1_stocks_proto_rawDescGZIP(), []int{0}
}

func (x *StocksInfoRequest) GetSku() int64 {
	if x != nil {
		return x.Sku
	}
	return 0
}

type StocksInfoResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Count         uint32                 `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StocksInfoResponse) Reset() {
	*x = StocksInfoResponse{}
	mi := &file_stocks_v1_stocks_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StocksInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StocksInfoResponse) ProtoMessage() {}

func (x *StocksInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stocks_v1_stocks_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StocksInfoResponse.ProtoReflect.Descriptor instead.
func (*StocksInfoResponse) Descriptor() ([]byte, []int) {
	return file_stocks_v1_stocks_proto_rawDescGZIP(), []int{1}
}

func (x *StocksInfoResponse) GetCount() uint32 {
	if x != nil {
		return x.Count
	}
	return 0
}

var File_stocks_v1_stocks_proto protoreflect.FileDescriptor

var file_stocks_v1_stocks_proto_rawDesc = string([]byte{
	0x0a, 0x16, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x63,
	0x6b, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x73,
	0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e,
	0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2e,
	0x0a, 0x11, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x19, 0x0a, 0x03, 0x73, 0x6b, 0x75, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x42, 0x07, 0xba, 0x48, 0x04, 0x22, 0x02, 0x20, 0x00, 0x52, 0x03, 0x73, 0x6b, 0x75, 0x22, 0x2a,
	0x0a, 0x12, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x32, 0x6f, 0x0a, 0x0d, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x73, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5e, 0x0a, 0x0a, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1c, 0x2e, 0x73, 0x74, 0x6f, 0x63,
	0x6b, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x73,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x13, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0d, 0x12, 0x0b,
	0x2f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x42, 0x27, 0x5a, 0x25, 0x72,
	0x6f, 0x75, 0x74, 0x65, 0x32, 0x35, 0x36, 0x2f, 0x6c, 0x6f, 0x6d, 0x73, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x73, 0x74, 0x6f, 0x63, 0x6b,
	0x73, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_stocks_v1_stocks_proto_rawDescOnce sync.Once
	file_stocks_v1_stocks_proto_rawDescData []byte
)

func file_stocks_v1_stocks_proto_rawDescGZIP() []byte {
	file_stocks_v1_stocks_proto_rawDescOnce.Do(func() {
		file_stocks_v1_stocks_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_stocks_v1_stocks_proto_rawDesc), len(file_stocks_v1_stocks_proto_rawDesc)))
	})
	return file_stocks_v1_stocks_proto_rawDescData
}

var file_stocks_v1_stocks_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_stocks_v1_stocks_proto_goTypes = []any{
	(*StocksInfoRequest)(nil),  // 0: stocks.v1.StocksInfoRequest
	(*StocksInfoResponse)(nil), // 1: stocks.v1.StocksInfoResponse
}
var file_stocks_v1_stocks_proto_depIdxs = []int32{
	0, // 0: stocks.v1.StocksService.StocksInfo:input_type -> stocks.v1.StocksInfoRequest
	1, // 1: stocks.v1.StocksService.StocksInfo:output_type -> stocks.v1.StocksInfoResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_stocks_v1_stocks_proto_init() }
func file_stocks_v1_stocks_proto_init() {
	if File_stocks_v1_stocks_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_stocks_v1_stocks_proto_rawDesc), len(file_stocks_v1_stocks_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stocks_v1_stocks_proto_goTypes,
		DependencyIndexes: file_stocks_v1_stocks_proto_depIdxs,
		MessageInfos:      file_stocks_v1_stocks_proto_msgTypes,
	}.Build()
	File_stocks_v1_stocks_proto = out.File
	file_stocks_v1_stocks_proto_goTypes = nil
	file_stocks_v1_stocks_proto_depIdxs = nil
}
