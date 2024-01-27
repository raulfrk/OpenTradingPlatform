// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.25.1
// source: proto/news.proto

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

type News struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          int64            `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Author      string           `protobuf:"bytes,2,opt,name=Author,proto3" json:"Author,omitempty"`
	CreatedAt   int64            `protobuf:"varint,3,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt   int64            `protobuf:"varint,4,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
	Headline    string           `protobuf:"bytes,5,opt,name=Headline,proto3" json:"Headline,omitempty"`
	Summary     string           `protobuf:"bytes,6,opt,name=Summary,proto3" json:"Summary,omitempty"`
	Content     string           `protobuf:"bytes,7,opt,name=Content,proto3" json:"Content,omitempty"`
	URL         string           `protobuf:"bytes,8,opt,name=URL,proto3" json:"URL,omitempty"`
	Symbols     []string         `protobuf:"bytes,9,rep,name=Symbols,proto3" json:"Symbols,omitempty"`
	Fingerprint string           `protobuf:"bytes,10,opt,name=Fingerprint,proto3" json:"Fingerprint,omitempty"`
	Source      string           `protobuf:"bytes,11,opt,name=Source,proto3" json:"Source,omitempty"`
	Sentiments  []*NewsSentiment `protobuf:"bytes,12,rep,name=Sentiments,proto3" json:"Sentiments,omitempty"`
}

func (x *News) Reset() {
	*x = News{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_news_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *News) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*News) ProtoMessage() {}

func (x *News) ProtoReflect() protoreflect.Message {
	mi := &file_proto_news_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use News.ProtoReflect.Descriptor instead.
func (*News) Descriptor() ([]byte, []int) {
	return file_proto_news_proto_rawDescGZIP(), []int{0}
}

func (x *News) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *News) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *News) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

func (x *News) GetUpdatedAt() int64 {
	if x != nil {
		return x.UpdatedAt
	}
	return 0
}

func (x *News) GetHeadline() string {
	if x != nil {
		return x.Headline
	}
	return ""
}

func (x *News) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *News) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *News) GetURL() string {
	if x != nil {
		return x.URL
	}
	return ""
}

func (x *News) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

func (x *News) GetFingerprint() string {
	if x != nil {
		return x.Fingerprint
	}
	return ""
}

func (x *News) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *News) GetSentiments() []*NewsSentiment {
	if x != nil {
		return x.Sentiments
	}
	return nil
}

type NewsSentiment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp                int64  `protobuf:"varint,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	News                     *News  `protobuf:"bytes,2,opt,name=News,proto3" json:"News,omitempty"`
	Sentiment                string `protobuf:"bytes,3,opt,name=Sentiment,proto3" json:"Sentiment,omitempty"`
	SentimentAnalysisProcess string `protobuf:"bytes,4,opt,name=SentimentAnalysisProcess,proto3" json:"SentimentAnalysisProcess,omitempty"`
	Fingerprint              string `protobuf:"bytes,5,opt,name=Fingerprint,proto3" json:"Fingerprint,omitempty"`
	LLM                      string `protobuf:"bytes,6,opt,name=LLM,proto3" json:"LLM,omitempty"`
	Symbol                   string `protobuf:"bytes,7,opt,name=Symbol,proto3" json:"Symbol,omitempty"`
	SystemPrompt             string `protobuf:"bytes,8,opt,name=SystemPrompt,proto3" json:"SystemPrompt,omitempty"`
	Failed                   bool   `protobuf:"varint,9,opt,name=Failed,proto3" json:"Failed,omitempty"`
}

func (x *NewsSentiment) Reset() {
	*x = NewsSentiment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_news_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewsSentiment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewsSentiment) ProtoMessage() {}

func (x *NewsSentiment) ProtoReflect() protoreflect.Message {
	mi := &file_proto_news_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewsSentiment.ProtoReflect.Descriptor instead.
func (*NewsSentiment) Descriptor() ([]byte, []int) {
	return file_proto_news_proto_rawDescGZIP(), []int{1}
}

func (x *NewsSentiment) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *NewsSentiment) GetNews() *News {
	if x != nil {
		return x.News
	}
	return nil
}

func (x *NewsSentiment) GetSentiment() string {
	if x != nil {
		return x.Sentiment
	}
	return ""
}

func (x *NewsSentiment) GetSentimentAnalysisProcess() string {
	if x != nil {
		return x.SentimentAnalysisProcess
	}
	return ""
}

func (x *NewsSentiment) GetFingerprint() string {
	if x != nil {
		return x.Fingerprint
	}
	return ""
}

func (x *NewsSentiment) GetLLM() string {
	if x != nil {
		return x.LLM
	}
	return ""
}

func (x *NewsSentiment) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *NewsSentiment) GetSystemPrompt() string {
	if x != nil {
		return x.SystemPrompt
	}
	return ""
}

func (x *NewsSentiment) GetFailed() bool {
	if x != nil {
		return x.Failed
	}
	return false
}

var File_proto_news_proto protoreflect.FileDescriptor

var file_proto_news_proto_rawDesc = []byte{
	0x0a, 0x10, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6e, 0x65, 0x77, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x08, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x22, 0xd9, 0x02, 0x0a,
	0x04, 0x4e, 0x65, 0x77, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x41, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x1c, 0x0a,
	0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x48, 0x65, 0x61,
	0x64, 0x6c, 0x69, 0x6e, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x48, 0x65, 0x61,
	0x64, 0x6c, 0x69, 0x6e, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x53, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12,
	0x18, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x55, 0x52, 0x4c,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x55, 0x52, 0x4c, 0x12, 0x18, 0x0a, 0x07, 0x53,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x53, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72, 0x70,
	0x72, 0x69, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x46, 0x69, 0x6e, 0x67,
	0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x37, 0x0a, 0x0a, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0c, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2e, 0x4e,
	0x65, 0x77, 0x73, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0a, 0x53, 0x65,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x22, 0xb3, 0x02, 0x0a, 0x0d, 0x4e, 0x65, 0x77,
	0x73, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x22, 0x0a, 0x04, 0x4e, 0x65, 0x77, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65,
	0x73, 0x2e, 0x4e, 0x65, 0x77, 0x73, 0x52, 0x04, 0x4e, 0x65, 0x77, 0x73, 0x12, 0x1c, 0x0a, 0x09,
	0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x53, 0x65, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x3a, 0x0a, 0x18, 0x53, 0x65,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x18, 0x53, 0x65,
	0x6e, 0x74, 0x69, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x6e, 0x61, 0x6c, 0x79, 0x73, 0x69, 0x73, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x12, 0x20, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x67, 0x65, 0x72,
	0x70, 0x72, 0x69, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x46, 0x69, 0x6e,
	0x67, 0x65, 0x72, 0x70, 0x72, 0x69, 0x6e, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4c, 0x4c, 0x4d, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x4c, 0x4c, 0x4d, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x53, 0x79, 0x6d, 0x62,
	0x6f, 0x6c, 0x12, 0x22, 0x0a, 0x0c, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x50, 0x72, 0x6f, 0x6d,
	0x70, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d,
	0x50, 0x72, 0x6f, 0x6d, 0x70, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x42, 0x0b,
	0x5a, 0x09, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x69, 0x65, 0x73, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_proto_news_proto_rawDescOnce sync.Once
	file_proto_news_proto_rawDescData = file_proto_news_proto_rawDesc
)

func file_proto_news_proto_rawDescGZIP() []byte {
	file_proto_news_proto_rawDescOnce.Do(func() {
		file_proto_news_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_news_proto_rawDescData)
	})
	return file_proto_news_proto_rawDescData
}

var file_proto_news_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_news_proto_goTypes = []interface{}{
	(*News)(nil),          // 0: entities.News
	(*NewsSentiment)(nil), // 1: entities.NewsSentiment
}
var file_proto_news_proto_depIdxs = []int32{
	1, // 0: entities.News.Sentiments:type_name -> entities.NewsSentiment
	0, // 1: entities.NewsSentiment.News:type_name -> entities.News
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_news_proto_init() }
func file_proto_news_proto_init() {
	if File_proto_news_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_news_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*News); i {
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
		file_proto_news_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewsSentiment); i {
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
			RawDescriptor: file_proto_news_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_news_proto_goTypes,
		DependencyIndexes: file_proto_news_proto_depIdxs,
		MessageInfos:      file_proto_news_proto_msgTypes,
	}.Build()
	File_proto_news_proto = out.File
	file_proto_news_proto_rawDesc = nil
	file_proto_news_proto_goTypes = nil
	file_proto_news_proto_depIdxs = nil
}
