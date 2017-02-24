// Code generated by protoc-gen-go.
// source: log.proto
// DO NOT EDIT!

/*
Package log_proto is a generated protocol buffer package.

It is generated from these files:
	log.proto

It has these top-level messages:
	LogPackage
*/
package log_proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Log package
type LogPackage struct {
	Project          *string `protobuf:"bytes,1,opt,name=project" json:"project,omitempty"`
	Service          *string `protobuf:"bytes,2,opt,name=service" json:"service,omitempty"`
	Level            *uint32 `protobuf:"varint,3,opt,name=level" json:"level,omitempty"`
	Log              []byte  `protobuf:"bytes,4,opt,name=log" json:"log,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *LogPackage) Reset()                    { *m = LogPackage{} }
func (m *LogPackage) String() string            { return proto.CompactTextString(m) }
func (*LogPackage) ProtoMessage()               {}
func (*LogPackage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LogPackage) GetProject() string {
	if m != nil && m.Project != nil {
		return *m.Project
	}
	return ""
}

func (m *LogPackage) GetService() string {
	if m != nil && m.Service != nil {
		return *m.Service
	}
	return ""
}

func (m *LogPackage) GetLevel() uint32 {
	if m != nil && m.Level != nil {
		return *m.Level
	}
	return 0
}

func (m *LogPackage) GetLog() []byte {
	if m != nil {
		return m.Log
	}
	return nil
}

func init() {
	proto.RegisterType((*LogPackage)(nil), "log_proto.LogPackage")
}

func init() { proto.RegisterFile("log.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 122 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcc, 0xc9, 0x4f, 0xd7,
	0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x02, 0x31, 0xe3, 0xc1, 0x4c, 0xa5, 0x0c, 0x2e, 0x2e, 0x9f,
	0xfc, 0xf4, 0x80, 0xc4, 0xe4, 0xec, 0xc4, 0xf4, 0x54, 0x21, 0x09, 0x2e, 0xf6, 0x82, 0xa2, 0xfc,
	0xac, 0xd4, 0xe4, 0x12, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0x17, 0x24, 0x53, 0x9c,
	0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a, 0xc1, 0x04, 0x91, 0x81, 0x72, 0x85, 0x44, 0xb8, 0x58, 0x73,
	0x52, 0xcb, 0x52, 0x73, 0x24, 0x98, 0x15, 0x18, 0x35, 0x78, 0x83, 0x20, 0x1c, 0x21, 0x01, 0x2e,
	0xe6, 0x9c, 0xfc, 0x74, 0x09, 0x16, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x10, 0x13, 0x10, 0x00, 0x00,
	0xff, 0xff, 0xd6, 0x25, 0xfd, 0x3c, 0x80, 0x00, 0x00, 0x00,
}
