// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.12
// source: omni/specs/system.proto

package specs

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// DBVersionSpec keeps the current version of the COSI DB.
type DBVersionSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version uint64 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
}

func (x *DBVersionSpec) Reset() {
	*x = DBVersionSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omni_specs_system_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DBVersionSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DBVersionSpec) ProtoMessage() {}

func (x *DBVersionSpec) ProtoReflect() protoreflect.Message {
	mi := &file_omni_specs_system_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DBVersionSpec.ProtoReflect.Descriptor instead.
func (*DBVersionSpec) Descriptor() ([]byte, []int) {
	return file_omni_specs_system_proto_rawDescGZIP(), []int{0}
}

func (x *DBVersionSpec) GetVersion() uint64 {
	if x != nil {
		return x.Version
	}
	return 0
}

// SysVersionSpec keeps the current version of Omni.
type SysVersionSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BackendVersion string `protobuf:"bytes,1,opt,name=backend_version,json=backendVersion,proto3" json:"backend_version,omitempty"`
	InstanceName   string `protobuf:"bytes,2,opt,name=instance_name,json=instanceName,proto3" json:"instance_name,omitempty"`
}

func (x *SysVersionSpec) Reset() {
	*x = SysVersionSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_omni_specs_system_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SysVersionSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SysVersionSpec) ProtoMessage() {}

func (x *SysVersionSpec) ProtoReflect() protoreflect.Message {
	mi := &file_omni_specs_system_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SysVersionSpec.ProtoReflect.Descriptor instead.
func (*SysVersionSpec) Descriptor() ([]byte, []int) {
	return file_omni_specs_system_proto_rawDescGZIP(), []int{1}
}

func (x *SysVersionSpec) GetBackendVersion() string {
	if x != nil {
		return x.BackendVersion
	}
	return ""
}

func (x *SysVersionSpec) GetInstanceName() string {
	if x != nil {
		return x.InstanceName
	}
	return ""
}

var File_omni_specs_system_proto protoreflect.FileDescriptor

var file_omni_specs_system_proto_rawDesc = []byte{
	0x0a, 0x17, 0x6f, 0x6d, 0x6e, 0x69, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2f, 0x73, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x73, 0x70, 0x65, 0x63, 0x73,
	0x22, 0x29, 0x0a, 0x0d, 0x44, 0x42, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65,
	0x63, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0x5e, 0x0a, 0x0e, 0x53,
	0x79, 0x73, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x12, 0x27, 0x0a,
	0x0f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x42, 0x32, 0x5a, 0x30, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x64, 0x65, 0x72, 0x6f,
	0x6c, 0x61, 0x62, 0x73, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_omni_specs_system_proto_rawDescOnce sync.Once
	file_omni_specs_system_proto_rawDescData = file_omni_specs_system_proto_rawDesc
)

func file_omni_specs_system_proto_rawDescGZIP() []byte {
	file_omni_specs_system_proto_rawDescOnce.Do(func() {
		file_omni_specs_system_proto_rawDescData = protoimpl.X.CompressGZIP(file_omni_specs_system_proto_rawDescData)
	})
	return file_omni_specs_system_proto_rawDescData
}

var file_omni_specs_system_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_omni_specs_system_proto_goTypes = []interface{}{
	(*DBVersionSpec)(nil),  // 0: specs.DBVersionSpec
	(*SysVersionSpec)(nil), // 1: specs.SysVersionSpec
}
var file_omni_specs_system_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_omni_specs_system_proto_init() }
func file_omni_specs_system_proto_init() {
	if File_omni_specs_system_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_omni_specs_system_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DBVersionSpec); i {
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
		file_omni_specs_system_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SysVersionSpec); i {
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
			RawDescriptor: file_omni_specs_system_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_omni_specs_system_proto_goTypes,
		DependencyIndexes: file_omni_specs_system_proto_depIdxs,
		MessageInfos:      file_omni_specs_system_proto_msgTypes,
	}.Build()
	File_omni_specs_system_proto = out.File
	file_omni_specs_system_proto_rawDesc = nil
	file_omni_specs_system_proto_goTypes = nil
	file_omni_specs_system_proto_depIdxs = nil
}
