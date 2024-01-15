// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: shared/proto/trade.proto

package entities

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

type Trade struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID          int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Symbol      string   `protobuf:"bytes,2,opt,name=Symbol,proto3" json:"Symbol,omitempty"`
	Exchange    string   `protobuf:"bytes,3,opt,name=Exchange,proto3" json:"Exchange,omitempty"`
	Price       float64  `protobuf:"fixed64,4,opt,name=Price,proto3" json:"Price,omitempty"`
	Size        float64  `protobuf:"fixed64,5,opt,name=Size,proto3" json:"Size,omitempty"`
	Timestamp   int64    `protobuf:"varint,6,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	TakerSide   string   `protobuf:"bytes,7,opt,name=TakerSide,proto3" json:"TakerSide,omitempty"`
	Conditions  []string `protobuf:"bytes,8,rep,name=Conditions,proto3" json:"Conditions,omitempty"`
	Tape        string   `protobuf:"bytes,9,opt,name=Tape,proto3" json:"Tape,omitempty"`
	Fingerprint string   `protobuf:"bytes,10,opt,name=Fingerprint,proto3" json:"Fingerprint,omitempty"`
	Update      string   `protobuf:"bytes,11,opt,name=Update,proto3" json:"Update,omitempty"`
	Source      string   `protobuf:"bytes,13,opt,name=Source,proto3" json:"Source,omitempty"`
	AssetClass  string   `protobuf:"bytes,14,opt,name=AssetClass,proto3" json:"AssetClass,omitempty"`
}

func (x *Trade) Reset() {
	*x = Trade{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_proto_trade_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Trade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Trade) ProtoMessage() {}

func (x *Trade) ProtoReflect() protoreflect.Message {
	mi := &file_shared_proto_trade_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Trade.ProtoReflect.Descriptor instead.
func (*Trade) Descriptor() ([]byte, []int) {
	return file_shared_proto_trade_proto_rawDescGZIP(), []int{0}
}

func (x *Trade) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Trade) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Trade) GetExchange() string {
	if x != nil {
		return x.Exchange
	}
	return ""
}

func (x *Trade) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Trade) GetSize() float64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Trade) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *Trade) GetTakerSide() string {
	if x != nil {
		return x.TakerSide
	}
	return ""
}

func (x *Trade) GetConditions() []string {
	if x != nil {
		return x.Conditions
	}
	return nil
}

func (x *Trade) GetTape() string {
	if x != nil {
		return x.Tape
	}
	return ""
}

func (x *Trade) GetFingerprint() string {
	if x != nil {
		return x.Fingerprint
	}
	return ""
}

func (x *Trade) GetUpdate() string {
	if x != nil {
		return x.Update
	}
	return ""
}

func (x *Trade) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Trade) GetAssetClass() string {
	if x != nil {
		return x.AssetClass
	}
	return ""
}

var File_shared_proto_trade_proto protoreflect.FileDescriptor

var file_shared_proto_trade_proto_rawDesc = []byte{
	0x0a, 0x18, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74,
	0x72, 0x61, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x69, 0x65, 0x73, 0x22, 0xd7, 0x02, 0x0a, 0x05, 0x54, 0x72, 0x61, 0x64, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x16,
	0x0a, 0x06, 0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x53, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e,
	0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x69, 0x7a, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1c, 0x0a, 0x09,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x61,
	0x6b, 0x65, 0x72, 0x53, 0x69, 0x64, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54,
	0x61, 0x6b, 0x65, 0x72, 0x53, 0x69, 0x64, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x64,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x43, 0x6f,
	0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x54, 0x61, 0x70, 0x65,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x54, 0x61, 0x70, 0x65, 0x12, 0x20, 0x0a, 0x0b,
	0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0b, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x41, 0x73, 0x73, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x18, 0x0e, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0a, 0x41, 0x73, 0x73, 0x65, 0x74, 0x43, 0x6c, 0x61, 0x73, 0x73, 0x42, 0x0b,
	0x5a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_shared_proto_trade_proto_rawDescOnce sync.Once
	file_shared_proto_trade_proto_rawDescData = file_shared_proto_trade_proto_rawDesc
)

func file_shared_proto_trade_proto_rawDescGZIP() []byte {
	file_shared_proto_trade_proto_rawDescOnce.Do(func() {
		file_shared_proto_trade_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_proto_trade_proto_rawDescData)
	})
	return file_shared_proto_trade_proto_rawDescData
}

var file_shared_proto_trade_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_shared_proto_trade_proto_goTypes = []interface{}{
	(*Trade)(nil), // 0: entities.Trade
}
var file_shared_proto_trade_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_shared_proto_trade_proto_init() }
func file_shared_proto_trade_proto_init() {
	if File_shared_proto_trade_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_proto_trade_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Trade); i {
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
			RawDescriptor: file_shared_proto_trade_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_proto_trade_proto_goTypes,
		DependencyIndexes: file_shared_proto_trade_proto_depIdxs,
		MessageInfos:      file_shared_proto_trade_proto_msgTypes,
	}.Build()
	File_shared_proto_trade_proto = out.File
	file_shared_proto_trade_proto_rawDesc = nil
	file_shared_proto_trade_proto_goTypes = nil
	file_shared_proto_trade_proto_depIdxs = nil
}
