// Code generated by protoc-gen-go. DO NOT EDIT.
// source: githubstats.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type GenerateReportRequest struct {
	Username             string   `protobuf:"bytes,1,opt,name=Username,proto3" json:"Username,omitempty"`
	Token                string   `protobuf:"bytes,2,opt,name=Token,proto3" json:"Token,omitempty"`
	RepositoryName       []string `protobuf:"bytes,3,rep,name=RepositoryName,proto3" json:"RepositoryName,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenerateReportRequest) Reset()         { *m = GenerateReportRequest{} }
func (m *GenerateReportRequest) String() string { return proto.CompactTextString(m) }
func (*GenerateReportRequest) ProtoMessage()    {}
func (*GenerateReportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4c99c921aa3a466, []int{0}
}

func (m *GenerateReportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateReportRequest.Unmarshal(m, b)
}
func (m *GenerateReportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateReportRequest.Marshal(b, m, deterministic)
}
func (m *GenerateReportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateReportRequest.Merge(m, src)
}
func (m *GenerateReportRequest) XXX_Size() int {
	return xxx_messageInfo_GenerateReportRequest.Size(m)
}
func (m *GenerateReportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateReportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateReportRequest proto.InternalMessageInfo

func (m *GenerateReportRequest) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *GenerateReportRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *GenerateReportRequest) GetRepositoryName() []string {
	if m != nil {
		return m.RepositoryName
	}
	return nil
}

type GenerateReportResponse struct {
	Repository           []*Repository `protobuf:"bytes,1,rep,name=Repository,proto3" json:"Repository,omitempty"`
	Err                  string        `protobuf:"bytes,2,opt,name=Err,proto3" json:"Err,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *GenerateReportResponse) Reset()         { *m = GenerateReportResponse{} }
func (m *GenerateReportResponse) String() string { return proto.CompactTextString(m) }
func (*GenerateReportResponse) ProtoMessage()    {}
func (*GenerateReportResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4c99c921aa3a466, []int{1}
}

func (m *GenerateReportResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenerateReportResponse.Unmarshal(m, b)
}
func (m *GenerateReportResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenerateReportResponse.Marshal(b, m, deterministic)
}
func (m *GenerateReportResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenerateReportResponse.Merge(m, src)
}
func (m *GenerateReportResponse) XXX_Size() int {
	return xxx_messageInfo_GenerateReportResponse.Size(m)
}
func (m *GenerateReportResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenerateReportResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenerateReportResponse proto.InternalMessageInfo

func (m *GenerateReportResponse) GetRepository() []*Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *GenerateReportResponse) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type Repository struct {
	RepositoryFullName   string   `protobuf:"bytes,1,opt,name=RepositoryFullName,proto3" json:"RepositoryFullName,omitempty"`
	RepositoryName       string   `protobuf:"bytes,2,opt,name=RepositoryName,proto3" json:"RepositoryName,omitempty"`
	CloneURL             string   `protobuf:"bytes,3,opt,name=CloneURL,proto3" json:"CloneURL,omitempty"`
	LatestCommitDate     string   `protobuf:"bytes,4,opt,name=LatestCommitDate,proto3" json:"LatestCommitDate,omitempty"`
	LatestCommitAuthor   string   `protobuf:"bytes,5,opt,name=LatestCommitAuthor,proto3" json:"LatestCommitAuthor,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Repository) Reset()         { *m = Repository{} }
func (m *Repository) String() string { return proto.CompactTextString(m) }
func (*Repository) ProtoMessage()    {}
func (*Repository) Descriptor() ([]byte, []int) {
	return fileDescriptor_b4c99c921aa3a466, []int{2}
}

func (m *Repository) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Repository.Unmarshal(m, b)
}
func (m *Repository) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Repository.Marshal(b, m, deterministic)
}
func (m *Repository) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Repository.Merge(m, src)
}
func (m *Repository) XXX_Size() int {
	return xxx_messageInfo_Repository.Size(m)
}
func (m *Repository) XXX_DiscardUnknown() {
	xxx_messageInfo_Repository.DiscardUnknown(m)
}

var xxx_messageInfo_Repository proto.InternalMessageInfo

func (m *Repository) GetRepositoryFullName() string {
	if m != nil {
		return m.RepositoryFullName
	}
	return ""
}

func (m *Repository) GetRepositoryName() string {
	if m != nil {
		return m.RepositoryName
	}
	return ""
}

func (m *Repository) GetCloneURL() string {
	if m != nil {
		return m.CloneURL
	}
	return ""
}

func (m *Repository) GetLatestCommitDate() string {
	if m != nil {
		return m.LatestCommitDate
	}
	return ""
}

func (m *Repository) GetLatestCommitAuthor() string {
	if m != nil {
		return m.LatestCommitAuthor
	}
	return ""
}

func init() {
	proto.RegisterType((*GenerateReportRequest)(nil), "pb.GenerateReportRequest")
	proto.RegisterType((*GenerateReportResponse)(nil), "pb.GenerateReportResponse")
	proto.RegisterType((*Repository)(nil), "pb.repository")
}

func init() { proto.RegisterFile("githubstats.proto", fileDescriptor_b4c99c921aa3a466) }

var fileDescriptor_b4c99c921aa3a466 = []byte{
	// 291 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0x41, 0x4f, 0xb3, 0x40,
	0x10, 0x0d, 0xe5, 0xeb, 0x17, 0x3b, 0x4d, 0x48, 0x9d, 0xa8, 0x59, 0x39, 0x11, 0x0e, 0x86, 0x78,
	0xe0, 0x50, 0x7f, 0x81, 0xa9, 0xca, 0xa5, 0xe9, 0x61, 0xb5, 0x1e, 0xbc, 0x41, 0x32, 0xb1, 0x44,
	0x60, 0xe9, 0xee, 0x70, 0xf0, 0xbf, 0xfa, 0x63, 0xcc, 0x82, 0x96, 0x4a, 0xb9, 0xed, 0xbc, 0xf7,
	0x76, 0x76, 0xe6, 0xbd, 0x85, 0xf3, 0xf7, 0x9c, 0x77, 0x4d, 0x66, 0x38, 0x65, 0x13, 0xd7, 0x5a,
	0xb1, 0xc2, 0x49, 0x9d, 0x85, 0x7b, 0xb8, 0x4c, 0xa8, 0x22, 0x9d, 0x32, 0x49, 0xaa, 0x95, 0x66,
	0x49, 0xfb, 0x86, 0x0c, 0xa3, 0x0f, 0x67, 0x5b, 0x43, 0xba, 0x4a, 0x4b, 0x12, 0x4e, 0xe0, 0x44,
	0x33, 0x79, 0xa8, 0xf1, 0x02, 0xa6, 0x2f, 0xea, 0x83, 0x2a, 0x31, 0x69, 0x89, 0xae, 0xc0, 0x1b,
	0xf0, 0x6c, 0x0b, 0x93, 0xb3, 0xd2, 0x9f, 0x1b, 0x7b, 0xcf, 0x0d, 0xdc, 0x68, 0x26, 0x07, 0x68,
	0xf8, 0x06, 0x57, 0xc3, 0x27, 0x4d, 0xad, 0x2a, 0x43, 0x18, 0x03, 0xf4, 0x5a, 0xe1, 0x04, 0x6e,
	0x34, 0x5f, 0x7a, 0x71, 0x9d, 0xc5, 0xfa, 0x80, 0xca, 0x23, 0x05, 0x2e, 0xc0, 0x7d, 0xd4, 0xfa,
	0x67, 0x0a, 0x7b, 0x0c, 0xbf, 0x1c, 0x80, 0x5e, 0x8c, 0x31, 0x60, 0x2f, 0x7f, 0x6a, 0x8a, 0x62,
	0xd3, 0xaf, 0x33, 0xc2, 0x8c, 0xac, 0xd0, 0xf5, 0x1e, 0xa0, 0xd6, 0x9c, 0x55, 0xa1, 0x2a, 0xda,
	0xca, 0xb5, 0x70, 0x3b, 0x73, 0x7e, 0x6b, 0xbc, 0x85, 0xc5, 0x3a, 0x65, 0x32, 0xbc, 0x52, 0x65,
	0x99, 0xf3, 0x43, 0xca, 0x24, 0xfe, 0xb5, 0x9a, 0x13, 0xdc, 0xce, 0x77, 0x8c, 0xdd, 0x37, 0xbc,
	0x53, 0x5a, 0x4c, 0xbb, 0xf9, 0x4e, 0x99, 0xe5, 0x2b, 0xcc, 0x93, 0x36, 0xc6, 0x67, 0x1b, 0x23,
	0x26, 0xe0, 0xfd, 0x75, 0x12, 0xaf, 0xad, 0x5b, 0xa3, 0x81, 0xfa, 0xfe, 0x18, 0xd5, 0x19, 0x9f,
	0xfd, 0x6f, 0x3f, 0xc4, 0xdd, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8c, 0x17, 0xbb, 0x5f, 0x25,
	0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// GithubStatsClient is the client API for GithubStats service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type GithubStatsClient interface {
	// Generates a github statistics report
	GenerateReport(ctx context.Context, in *GenerateReportRequest, opts ...grpc.CallOption) (*GenerateReportResponse, error)
}

type githubStatsClient struct {
	cc *grpc.ClientConn
}

func NewGithubStatsClient(cc *grpc.ClientConn) GithubStatsClient {
	return &githubStatsClient{cc}
}

func (c *githubStatsClient) GenerateReport(ctx context.Context, in *GenerateReportRequest, opts ...grpc.CallOption) (*GenerateReportResponse, error) {
	out := new(GenerateReportResponse)
	err := c.cc.Invoke(ctx, "/pb.GithubStats/GenerateReport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GithubStatsServer is the server API for GithubStats service.
type GithubStatsServer interface {
	// Generates a github statistics report
	GenerateReport(context.Context, *GenerateReportRequest) (*GenerateReportResponse, error)
}

func RegisterGithubStatsServer(s *grpc.Server, srv GithubStatsServer) {
	s.RegisterService(&_GithubStats_serviceDesc, srv)
}

func _GithubStats_GenerateReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GenerateReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GithubStatsServer).GenerateReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.GithubStats/GenerateReport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GithubStatsServer).GenerateReport(ctx, req.(*GenerateReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _GithubStats_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.GithubStats",
	HandlerType: (*GithubStatsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GenerateReport",
			Handler:    _GithubStats_GenerateReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "githubstats.proto",
}
