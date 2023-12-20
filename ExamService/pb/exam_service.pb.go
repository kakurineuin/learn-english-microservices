// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.25.1
// source: exam_service.proto

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

type CreateExamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic       string `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	IsPublic    bool   `protobuf:"varint,3,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	UserId      string `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *CreateExamRequest) Reset() {
	*x = CreateExamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateExamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateExamRequest) ProtoMessage() {}

func (x *CreateExamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateExamRequest.ProtoReflect.Descriptor instead.
func (*CreateExamRequest) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{0}
}

func (x *CreateExamRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *CreateExamRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreateExamRequest) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *CreateExamRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type CreateExamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExamId string `protobuf:"bytes,1,opt,name=exam_id,json=examId,proto3" json:"exam_id,omitempty"`
}

func (x *CreateExamResponse) Reset() {
	*x = CreateExamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateExamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateExamResponse) ProtoMessage() {}

func (x *CreateExamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateExamResponse.ProtoReflect.Descriptor instead.
func (*CreateExamResponse) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{1}
}

func (x *CreateExamResponse) GetExamId() string {
	if x != nil {
		return x.ExamId
	}
	return ""
}

type UpdateExamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExamId      string `protobuf:"bytes,1,opt,name=exam_id,json=examId,proto3" json:"exam_id,omitempty"`
	Topic       string `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	IsPublic    bool   `protobuf:"varint,4,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	UserId      string `protobuf:"bytes,5,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *UpdateExamRequest) Reset() {
	*x = UpdateExamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateExamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateExamRequest) ProtoMessage() {}

func (x *UpdateExamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateExamRequest.ProtoReflect.Descriptor instead.
func (*UpdateExamRequest) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateExamRequest) GetExamId() string {
	if x != nil {
		return x.ExamId
	}
	return ""
}

func (x *UpdateExamRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *UpdateExamRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *UpdateExamRequest) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *UpdateExamRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type UpdateExamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExamId string `protobuf:"bytes,1,opt,name=exam_id,json=examId,proto3" json:"exam_id,omitempty"`
}

func (x *UpdateExamResponse) Reset() {
	*x = UpdateExamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateExamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateExamResponse) ProtoMessage() {}

func (x *UpdateExamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateExamResponse.ProtoReflect.Descriptor instead.
func (*UpdateExamResponse) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateExamResponse) GetExamId() string {
	if x != nil {
		return x.ExamId
	}
	return ""
}

type FindExamsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageIndex int64  `protobuf:"varint,1,opt,name=page_index,json=pageIndex,proto3" json:"page_index,omitempty"`
	PageSize  int64  `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	UserId    string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *FindExamsRequest) Reset() {
	*x = FindExamsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindExamsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindExamsRequest) ProtoMessage() {}

func (x *FindExamsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindExamsRequest.ProtoReflect.Descriptor instead.
func (*FindExamsRequest) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{4}
}

func (x *FindExamsRequest) GetPageIndex() int64 {
	if x != nil {
		return x.PageIndex
	}
	return 0
}

func (x *FindExamsRequest) GetPageSize() int64 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *FindExamsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type FindExamsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Total     int64   `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	PageCount int64   `protobuf:"varint,2,opt,name=page_count,json=pageCount,proto3" json:"page_count,omitempty"`
	Exams     []*Exam `protobuf:"bytes,3,rep,name=exams,proto3" json:"exams,omitempty"`
}

func (x *FindExamsResponse) Reset() {
	*x = FindExamsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FindExamsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FindExamsResponse) ProtoMessage() {}

func (x *FindExamsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FindExamsResponse.ProtoReflect.Descriptor instead.
func (*FindExamsResponse) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{5}
}

func (x *FindExamsResponse) GetTotal() int64 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *FindExamsResponse) GetPageCount() int64 {
	if x != nil {
		return x.PageCount
	}
	return 0
}

func (x *FindExamsResponse) GetExams() []*Exam {
	if x != nil {
		return x.Exams
	}
	return nil
}

type Exam struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Topic       string   `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Tags        []string `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	IsPublic    bool     `protobuf:"varint,5,opt,name=is_public,json=isPublic,proto3" json:"is_public,omitempty"`
	UserId      string   `protobuf:"bytes,6,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *Exam) Reset() {
	*x = Exam{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Exam) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Exam) ProtoMessage() {}

func (x *Exam) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Exam.ProtoReflect.Descriptor instead.
func (*Exam) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{6}
}

func (x *Exam) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Exam) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *Exam) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *Exam) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Exam) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *Exam) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type DeleteExamRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExamId string `protobuf:"bytes,1,opt,name=exam_id,json=examId,proto3" json:"exam_id,omitempty"`
	UserId string `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *DeleteExamRequest) Reset() {
	*x = DeleteExamRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteExamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteExamRequest) ProtoMessage() {}

func (x *DeleteExamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteExamRequest.ProtoReflect.Descriptor instead.
func (*DeleteExamRequest) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteExamRequest) GetExamId() string {
	if x != nil {
		return x.ExamId
	}
	return ""
}

func (x *DeleteExamRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type DeleteExamResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DeleteExamResponse) Reset() {
	*x = DeleteExamResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_exam_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteExamResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteExamResponse) ProtoMessage() {}

func (x *DeleteExamResponse) ProtoReflect() protoreflect.Message {
	mi := &file_exam_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteExamResponse.ProtoReflect.Descriptor instead.
func (*DeleteExamResponse) Descriptor() ([]byte, []int) {
	return file_exam_service_proto_rawDescGZIP(), []int{8}
}

var File_exam_service_proto protoreflect.FileDescriptor

var file_exam_service_proto_rawDesc = []byte{
	0x0a, 0x12, 0x65, 0x78, 0x61, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x70, 0x62, 0x22, 0x81, 0x01, 0x0a, 0x11, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14,
	0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74,
	0x6f, 0x70, 0x69, 0x63, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x2d, 0x0a, 0x12,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x78, 0x61, 0x6d, 0x49, 0x64, 0x22, 0x9a, 0x01, 0x0a, 0x11,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x65, 0x78, 0x61, 0x6d, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f,
	0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63,
	0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x12,
	0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x2d, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x17,
	0x0a, 0x07, 0x65, 0x78, 0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x65, 0x78, 0x61, 0x6d, 0x49, 0x64, 0x22, 0x67, 0x0a, 0x10, 0x46, 0x69, 0x6e, 0x64, 0x45,
	0x78, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70,
	0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x70, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61,
	0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70,
	0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x22, 0x68, 0x0a, 0x11, 0x46, 0x69, 0x6e, 0x64, 0x45, 0x78, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x1d, 0x0a, 0x0a, 0x70,
	0x61, 0x67, 0x65, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x09, 0x70, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x1e, 0x0a, 0x05, 0x65, 0x78,
	0x61, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x70, 0x62, 0x2e, 0x45,
	0x78, 0x61, 0x6d, 0x52, 0x05, 0x65, 0x78, 0x61, 0x6d, 0x73, 0x22, 0x98, 0x01, 0x0a, 0x04, 0x45,
	0x78, 0x61, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x45, 0x0a, 0x11, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45,
	0x78, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x65, 0x78,
	0x61, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x65, 0x78, 0x61,
	0x6d, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x14, 0x0a, 0x12,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x32, 0xfe, 0x01, 0x0a, 0x0b, 0x45, 0x78, 0x61, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d,
	0x12, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3b, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x12, 0x15, 0x2e,
	0x70, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x09,
	0x46, 0x69, 0x6e, 0x64, 0x45, 0x78, 0x61, 0x6d, 0x73, 0x12, 0x14, 0x2e, 0x70, 0x62, 0x2e, 0x46,
	0x69, 0x6e, 0x64, 0x45, 0x78, 0x61, 0x6d, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x70, 0x62, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x45, 0x78, 0x61, 0x6d, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3b, 0x0a, 0x0a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x45, 0x78, 0x61, 0x6d, 0x12, 0x15, 0x2e, 0x70, 0x62, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x70, 0x62,
	0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x45, 0x78, 0x61, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_exam_service_proto_rawDescOnce sync.Once
	file_exam_service_proto_rawDescData = file_exam_service_proto_rawDesc
)

func file_exam_service_proto_rawDescGZIP() []byte {
	file_exam_service_proto_rawDescOnce.Do(func() {
		file_exam_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_exam_service_proto_rawDescData)
	})
	return file_exam_service_proto_rawDescData
}

var file_exam_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_exam_service_proto_goTypes = []interface{}{
	(*CreateExamRequest)(nil),  // 0: pb.CreateExamRequest
	(*CreateExamResponse)(nil), // 1: pb.CreateExamResponse
	(*UpdateExamRequest)(nil),  // 2: pb.UpdateExamRequest
	(*UpdateExamResponse)(nil), // 3: pb.UpdateExamResponse
	(*FindExamsRequest)(nil),   // 4: pb.FindExamsRequest
	(*FindExamsResponse)(nil),  // 5: pb.FindExamsResponse
	(*Exam)(nil),               // 6: pb.Exam
	(*DeleteExamRequest)(nil),  // 7: pb.DeleteExamRequest
	(*DeleteExamResponse)(nil), // 8: pb.DeleteExamResponse
}
var file_exam_service_proto_depIdxs = []int32{
	6, // 0: pb.FindExamsResponse.exams:type_name -> pb.Exam
	0, // 1: pb.ExamService.CreateExam:input_type -> pb.CreateExamRequest
	2, // 2: pb.ExamService.UpdateExam:input_type -> pb.UpdateExamRequest
	4, // 3: pb.ExamService.FindExams:input_type -> pb.FindExamsRequest
	7, // 4: pb.ExamService.DeleteExam:input_type -> pb.DeleteExamRequest
	1, // 5: pb.ExamService.CreateExam:output_type -> pb.CreateExamResponse
	3, // 6: pb.ExamService.UpdateExam:output_type -> pb.UpdateExamResponse
	5, // 7: pb.ExamService.FindExams:output_type -> pb.FindExamsResponse
	8, // 8: pb.ExamService.DeleteExam:output_type -> pb.DeleteExamResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_exam_service_proto_init() }
func file_exam_service_proto_init() {
	if File_exam_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_exam_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateExamRequest); i {
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
		file_exam_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateExamResponse); i {
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
		file_exam_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateExamRequest); i {
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
		file_exam_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateExamResponse); i {
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
		file_exam_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindExamsRequest); i {
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
		file_exam_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FindExamsResponse); i {
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
		file_exam_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Exam); i {
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
		file_exam_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteExamRequest); i {
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
		file_exam_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteExamResponse); i {
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
			RawDescriptor: file_exam_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_exam_service_proto_goTypes,
		DependencyIndexes: file_exam_service_proto_depIdxs,
		MessageInfos:      file_exam_service_proto_msgTypes,
	}.Build()
	File_exam_service_proto = out.File
	file_exam_service_proto_rawDesc = nil
	file_exam_service_proto_goTypes = nil
	file_exam_service_proto_depIdxs = nil
}
