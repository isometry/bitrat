// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.12.3
// source: bitrat.proto

package bitratpb

import (
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

type RecordSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Algorithm   string               `protobuf:"bytes,1,opt,name=Algorithm,proto3" json:"Algorithm,omitempty"`
	PathHashMap map[string]*HashData `protobuf:"bytes,2,rep,name=PathHashMap,proto3" json:"PathHashMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// repeated Record Record = 2;
	Statistics *Statistics `protobuf:"bytes,3,opt,name=Statistics,proto3" json:"Statistics,omitempty"`
}

func (x *RecordSet) Reset() {
	*x = RecordSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bitrat_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordSet) ProtoMessage() {}

func (x *RecordSet) ProtoReflect() protoreflect.Message {
	mi := &file_bitrat_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordSet.ProtoReflect.Descriptor instead.
func (*RecordSet) Descriptor() ([]byte, []int) {
	return file_bitrat_proto_rawDescGZIP(), []int{0}
}

func (x *RecordSet) GetAlgorithm() string {
	if x != nil {
		return x.Algorithm
	}
	return ""
}

func (x *RecordSet) GetPathHashMap() map[string]*HashData {
	if x != nil {
		return x.PathHashMap
	}
	return nil
}

func (x *RecordSet) GetStatistics() *Statistics {
	if x != nil {
		return x.Statistics
	}
	return nil
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Path string               `protobuf:"bytes,1,opt,name=Path,proto3" json:"Path,omitempty"`
	Hash []byte               `protobuf:"bytes,2,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Size int64                `protobuf:"varint,3,opt,name=Size,proto3" json:"Size,omitempty"`
	Time *timestamp.Timestamp `protobuf:"bytes,4,opt,name=Time,proto3" json:"Time,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bitrat_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_bitrat_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_bitrat_proto_rawDescGZIP(), []int{1}
}

func (x *Record) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *Record) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *Record) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Record) GetTime() *timestamp.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type Statistics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NumFiles    int64              `protobuf:"varint,1,opt,name=NumFiles,proto3" json:"NumFiles,omitempty"`
	TotalBytes  int64              `protobuf:"varint,2,opt,name=TotalBytes,proto3" json:"TotalBytes,omitempty"`
	ElapsedTime *duration.Duration `protobuf:"bytes,3,opt,name=ElapsedTime,proto3" json:"ElapsedTime,omitempty"`
	TotalTime   *duration.Duration `protobuf:"bytes,4,opt,name=TotalTime,proto3" json:"TotalTime,omitempty"`
	Parallel    int32              `protobuf:"varint,5,opt,name=Parallel,proto3" json:"Parallel,omitempty"`
}

func (x *Statistics) Reset() {
	*x = Statistics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bitrat_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Statistics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Statistics) ProtoMessage() {}

func (x *Statistics) ProtoReflect() protoreflect.Message {
	mi := &file_bitrat_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Statistics.ProtoReflect.Descriptor instead.
func (*Statistics) Descriptor() ([]byte, []int) {
	return file_bitrat_proto_rawDescGZIP(), []int{2}
}

func (x *Statistics) GetNumFiles() int64 {
	if x != nil {
		return x.NumFiles
	}
	return 0
}

func (x *Statistics) GetTotalBytes() int64 {
	if x != nil {
		return x.TotalBytes
	}
	return 0
}

func (x *Statistics) GetElapsedTime() *duration.Duration {
	if x != nil {
		return x.ElapsedTime
	}
	return nil
}

func (x *Statistics) GetTotalTime() *duration.Duration {
	if x != nil {
		return x.TotalTime
	}
	return nil
}

func (x *Statistics) GetParallel() int32 {
	if x != nil {
		return x.Parallel
	}
	return 0
}

type HashData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Hash    []byte               `protobuf:"bytes,1,opt,name=Hash,proto3" json:"Hash,omitempty"`
	Size    int64                `protobuf:"varint,2,opt,name=Size,proto3" json:"Size,omitempty"`
	ModTime *timestamp.Timestamp `protobuf:"bytes,3,opt,name=ModTime,proto3" json:"ModTime,omitempty"`
}

func (x *HashData) Reset() {
	*x = HashData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_bitrat_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HashData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HashData) ProtoMessage() {}

func (x *HashData) ProtoReflect() protoreflect.Message {
	mi := &file_bitrat_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HashData.ProtoReflect.Descriptor instead.
func (*HashData) Descriptor() ([]byte, []int) {
	return file_bitrat_proto_rawDescGZIP(), []int{3}
}

func (x *HashData) GetHash() []byte {
	if x != nil {
		return x.Hash
	}
	return nil
}

func (x *HashData) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *HashData) GetModTime() *timestamp.Timestamp {
	if x != nil {
		return x.ModTime
	}
	return nil
}

var File_bitrat_proto protoreflect.FileDescriptor

var file_bitrat_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x62, 0x69, 0x74, 0x72, 0x61, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08,
	0x62, 0x69, 0x74, 0x72, 0x61, 0x74, 0x70, 0x62, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xfb, 0x01, 0x0a, 0x09, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x41, 0x6c, 0x67, 0x6f, 0x72,
	0x69, 0x74, 0x68, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x41, 0x6c, 0x67, 0x6f,
	0x72, 0x69, 0x74, 0x68, 0x6d, 0x12, 0x46, 0x0a, 0x0b, 0x50, 0x61, 0x74, 0x68, 0x48, 0x61, 0x73,
	0x68, 0x4d, 0x61, 0x70, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x62, 0x69, 0x74,
	0x72, 0x61, 0x74, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x53, 0x65, 0x74, 0x2e,
	0x50, 0x61, 0x74, 0x68, 0x48, 0x61, 0x73, 0x68, 0x4d, 0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x52, 0x0b, 0x50, 0x61, 0x74, 0x68, 0x48, 0x61, 0x73, 0x68, 0x4d, 0x61, 0x70, 0x12, 0x34, 0x0a,
	0x0a, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x62, 0x69, 0x74, 0x72, 0x61, 0x74, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x52, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74,
	0x69, 0x63, 0x73, 0x1a, 0x52, 0x0a, 0x10, 0x50, 0x61, 0x74, 0x68, 0x48, 0x61, 0x73, 0x68, 0x4d,
	0x61, 0x70, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x28, 0x0a, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x62, 0x69, 0x74, 0x72, 0x61,
	0x74, 0x70, 0x62, 0x2e, 0x48, 0x61, 0x73, 0x68, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x74, 0x0a, 0x06, 0x52, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x50, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x50, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x69, 0x7a,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x2e, 0x0a,
	0x04, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x04, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xda, 0x01,
	0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x1a, 0x0a, 0x08,
	0x4e, 0x75, 0x6d, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08,
	0x4e, 0x75, 0x6d, 0x46, 0x69, 0x6c, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x6f, 0x74, 0x61,
	0x6c, 0x42, 0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x54, 0x6f,
	0x74, 0x61, 0x6c, 0x42, 0x79, 0x74, 0x65, 0x73, 0x12, 0x3b, 0x0a, 0x0b, 0x45, 0x6c, 0x61, 0x70,
	0x73, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x44, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x45, 0x6c, 0x61, 0x70, 0x73, 0x65,
	0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x09, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x50, 0x61, 0x72, 0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x50, 0x61, 0x72, 0x61, 0x6c, 0x6c, 0x65, 0x6c, 0x22, 0x68, 0x0a, 0x08, 0x48, 0x61,
	0x73, 0x68, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x12, 0x0a, 0x04, 0x53, 0x69,
	0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x34,
	0x0a, 0x07, 0x4d, 0x6f, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x4d, 0x6f, 0x64,
	0x54, 0x69, 0x6d, 0x65, 0x42, 0x2e, 0x5a, 0x2c, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x69, 0x73, 0x6f, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2f, 0x62, 0x69, 0x74, 0x72,
	0x61, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x62, 0x69, 0x74, 0x72,
	0x61, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_bitrat_proto_rawDescOnce sync.Once
	file_bitrat_proto_rawDescData = file_bitrat_proto_rawDesc
)

func file_bitrat_proto_rawDescGZIP() []byte {
	file_bitrat_proto_rawDescOnce.Do(func() {
		file_bitrat_proto_rawDescData = protoimpl.X.CompressGZIP(file_bitrat_proto_rawDescData)
	})
	return file_bitrat_proto_rawDescData
}

var file_bitrat_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_bitrat_proto_goTypes = []interface{}{
	(*RecordSet)(nil),           // 0: bitratpb.RecordSet
	(*Record)(nil),              // 1: bitratpb.Record
	(*Statistics)(nil),          // 2: bitratpb.Statistics
	(*HashData)(nil),            // 3: bitratpb.HashData
	nil,                         // 4: bitratpb.RecordSet.PathHashMapEntry
	(*timestamp.Timestamp)(nil), // 5: google.protobuf.Timestamp
	(*duration.Duration)(nil),   // 6: google.protobuf.Duration
}
var file_bitrat_proto_depIdxs = []int32{
	4, // 0: bitratpb.RecordSet.PathHashMap:type_name -> bitratpb.RecordSet.PathHashMapEntry
	2, // 1: bitratpb.RecordSet.Statistics:type_name -> bitratpb.Statistics
	5, // 2: bitratpb.Record.Time:type_name -> google.protobuf.Timestamp
	6, // 3: bitratpb.Statistics.ElapsedTime:type_name -> google.protobuf.Duration
	6, // 4: bitratpb.Statistics.TotalTime:type_name -> google.protobuf.Duration
	5, // 5: bitratpb.HashData.ModTime:type_name -> google.protobuf.Timestamp
	3, // 6: bitratpb.RecordSet.PathHashMapEntry.value:type_name -> bitratpb.HashData
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_bitrat_proto_init() }
func file_bitrat_proto_init() {
	if File_bitrat_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_bitrat_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordSet); i {
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
		file_bitrat_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
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
		file_bitrat_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Statistics); i {
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
		file_bitrat_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HashData); i {
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
			RawDescriptor: file_bitrat_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_bitrat_proto_goTypes,
		DependencyIndexes: file_bitrat_proto_depIdxs,
		MessageInfos:      file_bitrat_proto_msgTypes,
	}.Build()
	File_bitrat_proto = out.File
	file_bitrat_proto_rawDesc = nil
	file_bitrat_proto_goTypes = nil
	file_bitrat_proto_depIdxs = nil
}
