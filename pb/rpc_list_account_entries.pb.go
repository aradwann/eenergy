// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.19.6
// source: rpc_list_account_entries.proto

package pb

import (
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

type ListAccountEntriesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountId int64 `protobuf:"varint,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Limit     int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	Offset    int32 `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
}

func (x *ListAccountEntriesRequest) Reset() {
	*x = ListAccountEntriesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_list_account_entries_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAccountEntriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAccountEntriesRequest) ProtoMessage() {}

func (x *ListAccountEntriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_list_account_entries_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAccountEntriesRequest.ProtoReflect.Descriptor instead.
func (*ListAccountEntriesRequest) Descriptor() ([]byte, []int) {
	return file_rpc_list_account_entries_proto_rawDescGZIP(), []int{0}
}

func (x *ListAccountEntriesRequest) GetAccountId() int64 {
	if x != nil {
		return x.AccountId
	}
	return 0
}

func (x *ListAccountEntriesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *ListAccountEntriesRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

type ListAccountEntriesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Entries []*Entry `protobuf:"bytes,1,rep,name=entries,proto3" json:"entries,omitempty"`
}

func (x *ListAccountEntriesResponse) Reset() {
	*x = ListAccountEntriesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_list_account_entries_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListAccountEntriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListAccountEntriesResponse) ProtoMessage() {}

func (x *ListAccountEntriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_list_account_entries_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListAccountEntriesResponse.ProtoReflect.Descriptor instead.
func (*ListAccountEntriesResponse) Descriptor() ([]byte, []int) {
	return file_rpc_list_account_entries_proto_rawDescGZIP(), []int{1}
}

func (x *ListAccountEntriesResponse) GetEntries() []*Entry {
	if x != nil {
		return x.Entries
	}
	return nil
}

var File_rpc_list_account_entries_proto protoreflect.FileDescriptor

var file_rpc_list_account_entries_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x72, 0x70, 0x63, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x61, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x5f, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x02, 0x70, 0x62, 0x1a, 0x0d, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x68, 0x0a, 0x19, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x41, 0x0a,
	0x1a, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x45, 0x6e, 0x74, 0x72,
	0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x23, 0x0a, 0x07, 0x65,
	0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x70,
	0x62, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x65, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x42, 0x20, 0x5a, 0x1e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61,
	0x72, 0x61, 0x64, 0x77, 0x61, 0x6e, 0x6e, 0x2f, 0x65, 0x65, 0x6e, 0x65, 0x72, 0x67, 0x79, 0x2f,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_list_account_entries_proto_rawDescOnce sync.Once
	file_rpc_list_account_entries_proto_rawDescData = file_rpc_list_account_entries_proto_rawDesc
)

func file_rpc_list_account_entries_proto_rawDescGZIP() []byte {
	file_rpc_list_account_entries_proto_rawDescOnce.Do(func() {
		file_rpc_list_account_entries_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_list_account_entries_proto_rawDescData)
	})
	return file_rpc_list_account_entries_proto_rawDescData
}

var file_rpc_list_account_entries_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_rpc_list_account_entries_proto_goTypes = []interface{}{
	(*ListAccountEntriesRequest)(nil),  // 0: pb.ListAccountEntriesRequest
	(*ListAccountEntriesResponse)(nil), // 1: pb.ListAccountEntriesResponse
	(*Entry)(nil),                      // 2: pb.Entry
}
var file_rpc_list_account_entries_proto_depIdxs = []int32{
	2, // 0: pb.ListAccountEntriesResponse.entries:type_name -> pb.Entry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_list_account_entries_proto_init() }
func file_rpc_list_account_entries_proto_init() {
	if File_rpc_list_account_entries_proto != nil {
		return
	}
	file_account_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_rpc_list_account_entries_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAccountEntriesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_rpc_list_account_entries_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListAccountEntriesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_rpc_list_account_entries_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rpc_list_account_entries_proto_goTypes,
		DependencyIndexes: file_rpc_list_account_entries_proto_depIdxs,
		MessageInfos:      file_rpc_list_account_entries_proto_msgTypes,
	}.Build()
	File_rpc_list_account_entries_proto = out.File
	file_rpc_list_account_entries_proto_rawDesc = nil
	file_rpc_list_account_entries_proto_goTypes = nil
	file_rpc_list_account_entries_proto_depIdxs = nil
}