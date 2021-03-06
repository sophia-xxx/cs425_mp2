// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.21.0-devel
// 	protoc        v3.13.0
// source: tcpmessage.proto

package MessagePackage

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type MsgType int32

const (
	MsgType_PUT_MASTER       MsgType = 0  // Put request send to master
	MsgType_PUT_MASTER_REP   MsgType = 1  // Master reply to put request, send back target ip
	MsgType_PUT_P2P          MsgType = 2  // Put request send to target ip
	MsgType_PUT_P2P_ACK      MsgType = 3  // Target ip get the file header successfully and send back ack, once this is received, start sending file
	MsgType_WRITE_ACK        MsgType = 4  // when finish write, send ACK to client
	MsgType_GET_MASTER       MsgType = 5  // Get request send to master
	MsgType_GET_MASTER_REP   MsgType = 6  // Master reply to get request, send back target ip
	MsgType_GET_P2P          MsgType = 7  // Get request send to target ip
	MsgType_GET_P2P_ACK      MsgType = 8  // Target ip send back ack (may require to tell back file size)
	MsgType_GET_P2P_SIZE_ACK MsgType = 9  // The get request initiator tell file source ip that it get the file size info successfully, once this is received, start sending file
	MsgType_DELETE           MsgType = 10 // master send delete request to file node
	MsgType_DELETE_ACK       MsgType = 11 // server reply to DELETE message
	MsgType_DELETE_MASTER    MsgType = 12 // client send deletee request to master
	MsgType_LIST             MsgType = 13 // client send LIST request to master,
	MsgType_LIST_REP         MsgType = 14 // master send sdfsfilename list to server
)

// Enum value maps for MsgType.
var (
	MsgType_name = map[int32]string{
		0:  "PUT_MASTER",
		1:  "PUT_MASTER_REP",
		2:  "PUT_P2P",
		3:  "PUT_P2P_ACK",
		4:  "WRITE_ACK",
		5:  "GET_MASTER",
		6:  "GET_MASTER_REP",
		7:  "GET_P2P",
		8:  "GET_P2P_ACK",
		9:  "GET_P2P_SIZE_ACK",
		10: "DELETE",
		11: "DELETE_ACK",
		12: "DELETE_MASTER",
		13: "LIST",
		14: "LIST_REP",
	}
	MsgType_value = map[string]int32{
		"PUT_MASTER":       0,
		"PUT_MASTER_REP":   1,
		"PUT_P2P":          2,
		"PUT_P2P_ACK":      3,
		"WRITE_ACK":        4,
		"GET_MASTER":       5,
		"GET_MASTER_REP":   6,
		"GET_P2P":          7,
		"GET_P2P_ACK":      8,
		"GET_P2P_SIZE_ACK": 9,
		"DELETE":           10,
		"DELETE_ACK":       11,
		"DELETE_MASTER":    12,
		"LIST":             13,
		"LIST_REP":         14,
	}
)

func (x MsgType) Enum() *MsgType {
	p := new(MsgType)
	*p = x
	return p
}

func (x MsgType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MsgType) Descriptor() protoreflect.EnumDescriptor {
	return file_tcpmessage_proto_enumTypes[0].Descriptor()
}

func (MsgType) Type() protoreflect.EnumType {
	return &file_tcpmessage_proto_enumTypes[0]
}

func (x MsgType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MsgType.Descriptor instead.
func (MsgType) EnumDescriptor() ([]byte, []int) {
	return file_tcpmessage_proto_rawDescGZIP(), []int{0}
}

type TCPMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      MsgType  `protobuf:"varint,1,opt,name=Type,proto3,enum=tutorial.MsgType" json:"Type,omitempty"`
	FileName  string   `protobuf:"bytes,2,opt,name=fileName,proto3" json:"fileName,omitempty"`
	SenderIP  string   `protobuf:"bytes,3,opt,name=senderIP,proto3" json:"senderIP,omitempty"`
	PayLoad   []string `protobuf:"bytes,4,rep,name=PayLoad,proto3" json:"PayLoad,omitempty"`
	FileSize  int32    `protobuf:"varint,5,opt,name=fileSize,proto3" json:"fileSize,omitempty"`
	LocalPath string   `protobuf:"bytes,6,opt,name=localPath,proto3" json:"localPath,omitempty"`
}

func (x *TCPMessage) Reset() {
	*x = TCPMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_tcpmessage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TCPMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TCPMessage) ProtoMessage() {}

func (x *TCPMessage) ProtoReflect() protoreflect.Message {
	mi := &file_tcpmessage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TCPMessage.ProtoReflect.Descriptor instead.
func (*TCPMessage) Descriptor() ([]byte, []int) {
	return file_tcpmessage_proto_rawDescGZIP(), []int{0}
}

func (x *TCPMessage) GetType() MsgType {
	if x != nil {
		return x.Type
	}
	return MsgType_PUT_MASTER
}

func (x *TCPMessage) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

func (x *TCPMessage) GetSenderIP() string {
	if x != nil {
		return x.SenderIP
	}
	return ""
}

func (x *TCPMessage) GetPayLoad() []string {
	if x != nil {
		return x.PayLoad
	}
	return nil
}

func (x *TCPMessage) GetFileSize() int32 {
	if x != nil {
		return x.FileSize
	}
	return 0
}

func (x *TCPMessage) GetLocalPath() string {
	if x != nil {
		return x.LocalPath
	}
	return ""
}

var File_tcpmessage_proto protoreflect.FileDescriptor

var file_tcpmessage_proto_rawDesc = []byte{
	0x0a, 0x10, 0x74, 0x63, 0x70, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x74, 0x75, 0x74, 0x6f, 0x72, 0x69, 0x61, 0x6c, 0x22, 0xbf, 0x01, 0x0a,
	0x0a, 0x54, 0x43, 0x50, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x25, 0x0a, 0x04, 0x54,
	0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x74, 0x75, 0x74, 0x6f,
	0x72, 0x69, 0x61, 0x6c, 0x2e, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x50, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x50, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x61,
	0x79, 0x4c, 0x6f, 0x61, 0x64, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x50, 0x61, 0x79,
	0x4c, 0x6f, 0x61, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x50, 0x61, 0x74, 0x68, 0x2a, 0xf9,
	0x01, 0x0a, 0x07, 0x4d, 0x73, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x55,
	0x54, 0x5f, 0x4d, 0x41, 0x53, 0x54, 0x45, 0x52, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x50, 0x55,
	0x54, 0x5f, 0x4d, 0x41, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x50, 0x10, 0x01, 0x12, 0x0b,
	0x0a, 0x07, 0x50, 0x55, 0x54, 0x5f, 0x50, 0x32, 0x50, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b, 0x50,
	0x55, 0x54, 0x5f, 0x50, 0x32, 0x50, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09,
	0x57, 0x52, 0x49, 0x54, 0x45, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x04, 0x12, 0x0e, 0x0a, 0x0a, 0x47,
	0x45, 0x54, 0x5f, 0x4d, 0x41, 0x53, 0x54, 0x45, 0x52, 0x10, 0x05, 0x12, 0x12, 0x0a, 0x0e, 0x47,
	0x45, 0x54, 0x5f, 0x4d, 0x41, 0x53, 0x54, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x50, 0x10, 0x06, 0x12,
	0x0b, 0x0a, 0x07, 0x47, 0x45, 0x54, 0x5f, 0x50, 0x32, 0x50, 0x10, 0x07, 0x12, 0x0f, 0x0a, 0x0b,
	0x47, 0x45, 0x54, 0x5f, 0x50, 0x32, 0x50, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x08, 0x12, 0x14, 0x0a,
	0x10, 0x47, 0x45, 0x54, 0x5f, 0x50, 0x32, 0x50, 0x5f, 0x53, 0x49, 0x5a, 0x45, 0x5f, 0x41, 0x43,
	0x4b, 0x10, 0x09, 0x12, 0x0a, 0x0a, 0x06, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x10, 0x0a, 0x12,
	0x0e, 0x0a, 0x0a, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x41, 0x43, 0x4b, 0x10, 0x0b, 0x12,
	0x11, 0x0a, 0x0d, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x4d, 0x41, 0x53, 0x54, 0x45, 0x52,
	0x10, 0x0c, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x49, 0x53, 0x54, 0x10, 0x0d, 0x12, 0x0c, 0x0a, 0x08,
	0x4c, 0x49, 0x53, 0x54, 0x5f, 0x52, 0x45, 0x50, 0x10, 0x0e, 0x42, 0x12, 0x5a, 0x10, 0x2e, 0x2f,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x61, 0x67, 0x65, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_tcpmessage_proto_rawDescOnce sync.Once
	file_tcpmessage_proto_rawDescData = file_tcpmessage_proto_rawDesc
)

func file_tcpmessage_proto_rawDescGZIP() []byte {
	file_tcpmessage_proto_rawDescOnce.Do(func() {
		file_tcpmessage_proto_rawDescData = protoimpl.X.CompressGZIP(file_tcpmessage_proto_rawDescData)
	})
	return file_tcpmessage_proto_rawDescData
}

var file_tcpmessage_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_tcpmessage_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_tcpmessage_proto_goTypes = []interface{}{
	(MsgType)(0),       // 0: tutorial.MsgType
	(*TCPMessage)(nil), // 1: tutorial.TCPMessage
}
var file_tcpmessage_proto_depIdxs = []int32{
	0, // 0: tutorial.TCPMessage.Type:type_name -> tutorial.MsgType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_tcpmessage_proto_init() }
func file_tcpmessage_proto_init() {
	if File_tcpmessage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_tcpmessage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TCPMessage); i {
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
			RawDescriptor: file_tcpmessage_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_tcpmessage_proto_goTypes,
		DependencyIndexes: file_tcpmessage_proto_depIdxs,
		EnumInfos:         file_tcpmessage_proto_enumTypes,
		MessageInfos:      file_tcpmessage_proto_msgTypes,
	}.Build()
	File_tcpmessage_proto = out.File
	file_tcpmessage_proto_rawDesc = nil
	file_tcpmessage_proto_goTypes = nil
	file_tcpmessage_proto_depIdxs = nil
}
