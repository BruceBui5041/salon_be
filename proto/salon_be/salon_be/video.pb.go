// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v3.21.12
// source: salon_be/video.proto

package proto

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

type ProcessState int32

const (
	ProcessState_PENDING    ProcessState = 0
	ProcessState_PROCESSING ProcessState = 1
	ProcessState_DONE       ProcessState = 2
	ProcessState_ERROR      ProcessState = 3
	ProcessState_INQUEUE    ProcessState = 4
)

// Enum value maps for ProcessState.
var (
	ProcessState_name = map[int32]string{
		0: "PENDING",
		1: "PROCESSING",
		2: "DONE",
		3: "ERROR",
		4: "INQUEUE",
	}
	ProcessState_value = map[string]int32{
		"PENDING":    0,
		"PROCESSING": 1,
		"DONE":       2,
		"ERROR":      3,
		"INQUEUE":    4,
	}
)

func (x ProcessState) Enum() *ProcessState {
	p := new(ProcessState)
	*p = x
	return p
}

func (x ProcessState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProcessState) Descriptor() protoreflect.EnumDescriptor {
	return file_salon_be_video_proto_enumTypes[0].Descriptor()
}

func (ProcessState) Type() protoreflect.EnumType {
	return &file_salon_be_video_proto_enumTypes[0]
}

func (x ProcessState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProcessState.Descriptor instead.
func (ProcessState) EnumDescriptor() ([]byte, []int) {
	return file_salon_be_video_proto_rawDescGZIP(), []int{0}
}

type ProcessResolution int32

const (
	ProcessResolution_ProcessResolution_NONE ProcessResolution = 0
	ProcessResolution_RESOLUTION_360P        ProcessResolution = 1
	ProcessResolution_RESOLUTION_480P        ProcessResolution = 2
	ProcessResolution_RESOLUTION_720P        ProcessResolution = 3
	ProcessResolution_RESOLUTION_1080P       ProcessResolution = 4
)

// Enum value maps for ProcessResolution.
var (
	ProcessResolution_name = map[int32]string{
		0: "ProcessResolution_NONE",
		1: "RESOLUTION_360P",
		2: "RESOLUTION_480P",
		3: "RESOLUTION_720P",
		4: "RESOLUTION_1080P",
	}
	ProcessResolution_value = map[string]int32{
		"ProcessResolution_NONE": 0,
		"RESOLUTION_360P":        1,
		"RESOLUTION_480P":        2,
		"RESOLUTION_720P":        3,
		"RESOLUTION_1080P":       4,
	}
)

func (x ProcessResolution) Enum() *ProcessResolution {
	p := new(ProcessResolution)
	*p = x
	return p
}

func (x ProcessResolution) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ProcessResolution) Descriptor() protoreflect.EnumDescriptor {
	return file_salon_be_video_proto_enumTypes[1].Descriptor()
}

func (ProcessResolution) Type() protoreflect.EnumType {
	return &file_salon_be_video_proto_enumTypes[1]
}

func (x ProcessResolution) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ProcessResolution.Descriptor instead.
func (ProcessResolution) EnumDescriptor() ([]byte, []int) {
	return file_salon_be_video_proto_rawDescGZIP(), []int{1}
}

type VideoInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId           string            `protobuf:"bytes,1,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	ServiceId         string            `protobuf:"bytes,2,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	Description       string            `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	UploadedBy        string            `protobuf:"bytes,4,opt,name=uploaded_by,json=uploadedBy,proto3" json:"uploaded_by,omitempty"`
	Timestamp         string            `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	S3Key             string            `protobuf:"bytes,6,opt,name=s3_key,json=s3Key,proto3" json:"s3_key,omitempty"`
	Retry             int32             `protobuf:"varint,7,opt,name=retry,proto3" json:"retry,omitempty"`
	RequestResolution ProcessResolution `protobuf:"varint,8,opt,name=request_resolution,json=requestResolution,proto3,enum=salonbe.ProcessResolution" json:"request_resolution,omitempty"`
}

func (x *VideoInfo) Reset() {
	*x = VideoInfo{}
	mi := &file_salon_be_video_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoInfo) ProtoMessage() {}

func (x *VideoInfo) ProtoReflect() protoreflect.Message {
	mi := &file_salon_be_video_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoInfo.ProtoReflect.Descriptor instead.
func (*VideoInfo) Descriptor() ([]byte, []int) {
	return file_salon_be_video_proto_rawDescGZIP(), []int{0}
}

func (x *VideoInfo) GetVideoId() string {
	if x != nil {
		return x.VideoId
	}
	return ""
}

func (x *VideoInfo) GetServiceId() string {
	if x != nil {
		return x.ServiceId
	}
	return ""
}

func (x *VideoInfo) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *VideoInfo) GetUploadedBy() string {
	if x != nil {
		return x.UploadedBy
	}
	return ""
}

func (x *VideoInfo) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *VideoInfo) GetS3Key() string {
	if x != nil {
		return x.S3Key
	}
	return ""
}

func (x *VideoInfo) GetRetry() int32 {
	if x != nil {
		return x.Retry
	}
	return 0
}

func (x *VideoInfo) GetRequestResolution() ProcessResolution {
	if x != nil {
		return x.RequestResolution
	}
	return ProcessResolution_ProcessResolution_NONE
}

type VideoProcessProgress struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VideoId           string            `protobuf:"bytes,1,opt,name=video_id,json=videoId,proto3" json:"video_id,omitempty"`
	ServiceId         string            `protobuf:"bytes,2,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	Timestamp         int64             `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Progress          int32             `protobuf:"varint,4,opt,name=progress,proto3" json:"progress,omitempty"`
	State             ProcessState      `protobuf:"varint,5,opt,name=state,proto3,enum=salonbe.ProcessState" json:"state,omitempty"`
	ErrorMsg          string            `protobuf:"bytes,6,opt,name=error_msg,json=errorMsg,proto3" json:"error_msg,omitempty"`
	RequestResolution ProcessResolution `protobuf:"varint,7,opt,name=request_resolution,json=requestResolution,proto3,enum=salonbe.ProcessResolution" json:"request_resolution,omitempty"`
	Duration          int32             `protobuf:"varint,8,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *VideoProcessProgress) Reset() {
	*x = VideoProcessProgress{}
	mi := &file_salon_be_video_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VideoProcessProgress) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VideoProcessProgress) ProtoMessage() {}

func (x *VideoProcessProgress) ProtoReflect() protoreflect.Message {
	mi := &file_salon_be_video_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VideoProcessProgress.ProtoReflect.Descriptor instead.
func (*VideoProcessProgress) Descriptor() ([]byte, []int) {
	return file_salon_be_video_proto_rawDescGZIP(), []int{1}
}

func (x *VideoProcessProgress) GetVideoId() string {
	if x != nil {
		return x.VideoId
	}
	return ""
}

func (x *VideoProcessProgress) GetServiceId() string {
	if x != nil {
		return x.ServiceId
	}
	return ""
}

func (x *VideoProcessProgress) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *VideoProcessProgress) GetProgress() int32 {
	if x != nil {
		return x.Progress
	}
	return 0
}

func (x *VideoProcessProgress) GetState() ProcessState {
	if x != nil {
		return x.State
	}
	return ProcessState_PENDING
}

func (x *VideoProcessProgress) GetErrorMsg() string {
	if x != nil {
		return x.ErrorMsg
	}
	return ""
}

func (x *VideoProcessProgress) GetRequestResolution() ProcessResolution {
	if x != nil {
		return x.RequestResolution
	}
	return ProcessResolution_ProcessResolution_NONE
}

func (x *VideoProcessProgress) GetDuration() int32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

type ProcessNewVideoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status       string `protobuf:"bytes,1,opt,name=status,proto3" json:"status,omitempty"`
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	Retry        int32  `protobuf:"varint,3,opt,name=retry,proto3" json:"retry,omitempty"`
}

func (x *ProcessNewVideoResponse) Reset() {
	*x = ProcessNewVideoResponse{}
	mi := &file_salon_be_video_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProcessNewVideoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessNewVideoResponse) ProtoMessage() {}

func (x *ProcessNewVideoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_salon_be_video_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessNewVideoResponse.ProtoReflect.Descriptor instead.
func (*ProcessNewVideoResponse) Descriptor() ([]byte, []int) {
	return file_salon_be_video_proto_rawDescGZIP(), []int{2}
}

func (x *ProcessNewVideoResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *ProcessNewVideoResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *ProcessNewVideoResponse) GetRetry() int32 {
	if x != nil {
		return x.Retry
	}
	return 0
}

var File_salon_be_video_proto protoreflect.FileDescriptor

var file_salon_be_video_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x61, 0x6c, 0x6f, 0x6e, 0x5f, 0x62, 0x65, 0x2f, 0x76, 0x69, 0x64, 0x65, 0x6f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x61, 0x6c, 0x6f, 0x6e, 0x62, 0x65, 0x22,
	0x9e, 0x02, 0x0a, 0x09, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x0a,
	0x08, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x76, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x49, 0x64, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65,
	0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x75, 0x70, 0x6c,
	0x6f, 0x61, 0x64, 0x65, 0x64, 0x5f, 0x62, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x75, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x42, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x15, 0x0a, 0x06, 0x73, 0x33, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x33, 0x4b, 0x65, 0x79, 0x12,
	0x14, 0x0a, 0x05, 0x72, 0x65, 0x74, 0x72, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x72, 0x65, 0x74, 0x72, 0x79, 0x12, 0x49, 0x0a, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x5f, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x1a, 0x2e, 0x73, 0x61, 0x6c, 0x6f, 0x6e, 0x62, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63,
	0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x11, 0x72,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e,
	0x22, 0xbb, 0x02, 0x0a, 0x14, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x50, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x76, 0x69, 0x64,
	0x65, 0x6f, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x69, 0x64,
	0x65, 0x6f, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x67, 0x72, 0x65, 0x73, 0x73, 0x12, 0x2b, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x73,
	0x61, 0x6c, 0x6f, 0x6e, 0x62, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x73, 0x67, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x4d, 0x73, 0x67, 0x12, 0x49, 0x0a, 0x12, 0x72, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x73, 0x61, 0x6c, 0x6f, 0x6e, 0x62, 0x65, 0x2e, 0x50, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x11, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x6c,
	0x0a, 0x17, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x64, 0x65,
	0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x65, 0x74, 0x72, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x72, 0x65, 0x74, 0x72, 0x79, 0x2a, 0x4d, 0x0a, 0x0c,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x0b, 0x0a, 0x07,
	0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a, 0x50, 0x52, 0x4f,
	0x43, 0x45, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x44, 0x4f, 0x4e,
	0x45, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x03, 0x12, 0x0b,
	0x0a, 0x07, 0x49, 0x4e, 0x51, 0x55, 0x45, 0x55, 0x45, 0x10, 0x04, 0x2a, 0x84, 0x01, 0x0a, 0x11,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1a, 0x0a, 0x16, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x52, 0x65, 0x73, 0x6f,
	0x6c, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x13, 0x0a,
	0x0f, 0x52, 0x45, 0x53, 0x4f, 0x4c, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x33, 0x36, 0x30, 0x50,
	0x10, 0x01, 0x12, 0x13, 0x0a, 0x0f, 0x52, 0x45, 0x53, 0x4f, 0x4c, 0x55, 0x54, 0x49, 0x4f, 0x4e,
	0x5f, 0x34, 0x38, 0x30, 0x50, 0x10, 0x02, 0x12, 0x13, 0x0a, 0x0f, 0x52, 0x45, 0x53, 0x4f, 0x4c,
	0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x37, 0x32, 0x30, 0x50, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10,
	0x52, 0x45, 0x53, 0x4f, 0x4c, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x31, 0x30, 0x38, 0x30, 0x50,
	0x10, 0x04, 0x32, 0x6a, 0x0a, 0x16, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x50, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x50, 0x0a, 0x16,
	0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4e, 0x65, 0x77, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x2e, 0x73, 0x61, 0x6c, 0x6f, 0x6e, 0x62, 0x65,
	0x2e, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x20, 0x2e, 0x73, 0x61, 0x6c,
	0x6f, 0x6e, 0x62, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4e, 0x65, 0x77, 0x56,
	0x69, 0x64, 0x65, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x10,
	0x5a, 0x0e, 0x73, 0x61, 0x6c, 0x6f, 0x6e, 0x5f, 0x62, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_salon_be_video_proto_rawDescOnce sync.Once
	file_salon_be_video_proto_rawDescData = file_salon_be_video_proto_rawDesc
)

func file_salon_be_video_proto_rawDescGZIP() []byte {
	file_salon_be_video_proto_rawDescOnce.Do(func() {
		file_salon_be_video_proto_rawDescData = protoimpl.X.CompressGZIP(file_salon_be_video_proto_rawDescData)
	})
	return file_salon_be_video_proto_rawDescData
}

var file_salon_be_video_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_salon_be_video_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_salon_be_video_proto_goTypes = []any{
	(ProcessState)(0),               // 0: salonbe.ProcessState
	(ProcessResolution)(0),          // 1: salonbe.ProcessResolution
	(*VideoInfo)(nil),               // 2: salonbe.VideoInfo
	(*VideoProcessProgress)(nil),    // 3: salonbe.VideoProcessProgress
	(*ProcessNewVideoResponse)(nil), // 4: salonbe.ProcessNewVideoResponse
}
var file_salon_be_video_proto_depIdxs = []int32{
	1, // 0: salonbe.VideoInfo.request_resolution:type_name -> salonbe.ProcessResolution
	0, // 1: salonbe.VideoProcessProgress.state:type_name -> salonbe.ProcessState
	1, // 2: salonbe.VideoProcessProgress.request_resolution:type_name -> salonbe.ProcessResolution
	2, // 3: salonbe.VideoProcessingService.ProcessNewVideoRequest:input_type -> salonbe.VideoInfo
	4, // 4: salonbe.VideoProcessingService.ProcessNewVideoRequest:output_type -> salonbe.ProcessNewVideoResponse
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_salon_be_video_proto_init() }
func file_salon_be_video_proto_init() {
	if File_salon_be_video_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_salon_be_video_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_salon_be_video_proto_goTypes,
		DependencyIndexes: file_salon_be_video_proto_depIdxs,
		EnumInfos:         file_salon_be_video_proto_enumTypes,
		MessageInfos:      file_salon_be_video_proto_msgTypes,
	}.Build()
	File_salon_be_video_proto = out.File
	file_salon_be_video_proto_rawDesc = nil
	file_salon_be_video_proto_goTypes = nil
	file_salon_be_video_proto_depIdxs = nil
}
