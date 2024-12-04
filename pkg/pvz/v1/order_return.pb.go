// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: pvz/v1/order_return.proto

package pvz

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OrderReturn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  int64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	OrderId int64 `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *OrderReturn) Reset() {
	*x = OrderReturn{}
	mi := &file_pvz_v1_order_return_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *OrderReturn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OrderReturn) ProtoMessage() {}

func (x *OrderReturn) ProtoReflect() protoreflect.Message {
	mi := &file_pvz_v1_order_return_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OrderReturn.ProtoReflect.Descriptor instead.
func (*OrderReturn) Descriptor() ([]byte, []int) {
	return file_pvz_v1_order_return_proto_rawDescGZIP(), []int{0}
}

func (x *OrderReturn) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *OrderReturn) GetOrderId() int64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type CreateOrderReturnRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId  uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	OrderId uint64 `protobuf:"varint,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *CreateOrderReturnRequest) Reset() {
	*x = CreateOrderReturnRequest{}
	mi := &file_pvz_v1_order_return_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateOrderReturnRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateOrderReturnRequest) ProtoMessage() {}

func (x *CreateOrderReturnRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pvz_v1_order_return_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateOrderReturnRequest.ProtoReflect.Descriptor instead.
func (*CreateOrderReturnRequest) Descriptor() ([]byte, []int) {
	return file_pvz_v1_order_return_proto_rawDescGZIP(), []int{1}
}

func (x *CreateOrderReturnRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateOrderReturnRequest) GetOrderId() uint64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type GiveOrderToCourierRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderId uint64 `protobuf:"varint,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
}

func (x *GiveOrderToCourierRequest) Reset() {
	*x = GiveOrderToCourierRequest{}
	mi := &file_pvz_v1_order_return_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GiveOrderToCourierRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GiveOrderToCourierRequest) ProtoMessage() {}

func (x *GiveOrderToCourierRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pvz_v1_order_return_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GiveOrderToCourierRequest.ProtoReflect.Descriptor instead.
func (*GiveOrderToCourierRequest) Descriptor() ([]byte, []int) {
	return file_pvz_v1_order_return_proto_rawDescGZIP(), []int{2}
}

func (x *GiveOrderToCourierRequest) GetOrderId() uint64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type GetOrderReturnsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page     int64 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PageSize int64 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *GetOrderReturnsRequest) Reset() {
	*x = GetOrderReturnsRequest{}
	mi := &file_pvz_v1_order_return_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrderReturnsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderReturnsRequest) ProtoMessage() {}

func (x *GetOrderReturnsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pvz_v1_order_return_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderReturnsRequest.ProtoReflect.Descriptor instead.
func (*GetOrderReturnsRequest) Descriptor() ([]byte, []int) {
	return file_pvz_v1_order_return_proto_rawDescGZIP(), []int{3}
}

func (x *GetOrderReturnsRequest) GetPage() int64 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *GetOrderReturnsRequest) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type GetOrderReturnsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderReturns []*OrderReturn `protobuf:"bytes,1,rep,name=order_returns,json=orderReturns,proto3" json:"order_returns,omitempty"`
}

func (x *GetOrderReturnsResponse) Reset() {
	*x = GetOrderReturnsResponse{}
	mi := &file_pvz_v1_order_return_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetOrderReturnsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetOrderReturnsResponse) ProtoMessage() {}

func (x *GetOrderReturnsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pvz_v1_order_return_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetOrderReturnsResponse.ProtoReflect.Descriptor instead.
func (*GetOrderReturnsResponse) Descriptor() ([]byte, []int) {
	return file_pvz_v1_order_return_proto_rawDescGZIP(), []int{4}
}

func (x *GetOrderReturnsResponse) GetOrderReturns() []*OrderReturn {
	if x != nil {
		return x.OrderReturns
	}
	return nil
}

var File_pvz_v1_order_return_proto protoreflect.FileDescriptor

var file_pvz_v1_order_return_proto_rawDesc = []byte{
	0x0a, 0x19, 0x70, 0x76, 0x7a, 0x2f, 0x76, 0x31, 0x2f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x72,
	0x65, 0x74, 0x75, 0x72, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x70, 0x76, 0x7a,
	0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64,
	0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x62, 0x65, 0x68, 0x61,
	0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x41, 0x0a, 0x0b, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x58, 0x0a,
	0x18, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72,
	0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x07,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3b, 0x0a, 0x19, 0x47, 0x69, 0x76, 0x65, 0x4f,
	0x72, 0x64, 0x65, 0x72, 0x54, 0x6f, 0x43, 0x6f, 0x75, 0x72, 0x69, 0x65, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x07, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x49, 0x64, 0x22, 0x61, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x4f, 0x72, 0x64, 0x65, 0x72,
	0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e,
	0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x42, 0x0a, 0xe0, 0x41,
	0x02, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x28, 0x01, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x27,
	0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x42, 0x0a, 0xe0, 0x41, 0x02, 0xfa, 0x42, 0x04, 0x22, 0x02, 0x28, 0x01, 0x52, 0x08, 0x70,
	0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x50, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4f, 0x72,
	0x64, 0x65, 0x72, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x35, 0x0a, 0x0d, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x74, 0x75,
	0x72, 0x6e, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x76, 0x7a, 0x2e,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x52, 0x0c, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x52, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x73, 0x42, 0x09, 0x5a, 0x07, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x76, 0x7a, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pvz_v1_order_return_proto_rawDescOnce sync.Once
	file_pvz_v1_order_return_proto_rawDescData = file_pvz_v1_order_return_proto_rawDesc
)

func file_pvz_v1_order_return_proto_rawDescGZIP() []byte {
	file_pvz_v1_order_return_proto_rawDescOnce.Do(func() {
		file_pvz_v1_order_return_proto_rawDescData = protoimpl.X.CompressGZIP(file_pvz_v1_order_return_proto_rawDescData)
	})
	return file_pvz_v1_order_return_proto_rawDescData
}

var file_pvz_v1_order_return_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pvz_v1_order_return_proto_goTypes = []any{
	(*OrderReturn)(nil),               // 0: pvz.OrderReturn
	(*CreateOrderReturnRequest)(nil),  // 1: pvz.CreateOrderReturnRequest
	(*GiveOrderToCourierRequest)(nil), // 2: pvz.GiveOrderToCourierRequest
	(*GetOrderReturnsRequest)(nil),    // 3: pvz.GetOrderReturnsRequest
	(*GetOrderReturnsResponse)(nil),   // 4: pvz.GetOrderReturnsResponse
}
var file_pvz_v1_order_return_proto_depIdxs = []int32{
	0, // 0: pvz.GetOrderReturnsResponse.order_returns:type_name -> pvz.OrderReturn
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_pvz_v1_order_return_proto_init() }
func file_pvz_v1_order_return_proto_init() {
	if File_pvz_v1_order_return_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pvz_v1_order_return_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pvz_v1_order_return_proto_goTypes,
		DependencyIndexes: file_pvz_v1_order_return_proto_depIdxs,
		MessageInfos:      file_pvz_v1_order_return_proto_msgTypes,
	}.Build()
	File_pvz_v1_order_return_proto = out.File
	file_pvz_v1_order_return_proto_rawDesc = nil
	file_pvz_v1_order_return_proto_goTypes = nil
	file_pvz_v1_order_return_proto_depIdxs = nil
}
