// Code generated by protoc-gen-go. DO NOT EDIT.
// source: consignment/consignment.proto

package consignment

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "golang.org/x/net/context"
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

// 货轮承运的一批货物
type Consignment struct {
	Id                   string       `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty" bson:"-"`
	Description          string       `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty" bson:"description,omitempty"`
	Weight               int32        `protobuf:"varint,3,opt,name=weight,proto3" json:"weight,omitempty" bson:"weight,omitempty"`
	Containers           []*Container `protobuf:"bytes,4,rep,name=containers,proto3" json:"containers,omitempty" bson:"containers,omitempty"`
	VesselId             string       `protobuf:"bytes,5,opt,name=vessel_id,json=vesselId,proto3" json:"vessel_id,omitempty" bson:"vessel_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-" bson:"-"`
	XXX_unrecognized     []byte       `json:"-" bson:"-"`
	XXX_sizecache        int32        `json:"-" bson:"-"`
}

func (m *Consignment) Reset()         { *m = Consignment{} }
func (m *Consignment) String() string { return proto.CompactTextString(m) }
func (*Consignment) ProtoMessage()    {}
func (*Consignment) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c6f11f9923110da, []int{0}
}

func (m *Consignment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Consignment.Unmarshal(m, b)
}
func (m *Consignment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Consignment.Marshal(b, m, deterministic)
}
func (m *Consignment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Consignment.Merge(m, src)
}
func (m *Consignment) XXX_Size() int {
	return xxx_messageInfo_Consignment.Size(m)
}
func (m *Consignment) XXX_DiscardUnknown() {
	xxx_messageInfo_Consignment.DiscardUnknown(m)
}

var xxx_messageInfo_Consignment proto.InternalMessageInfo

func (m *Consignment) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Consignment) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Consignment) GetWeight() int32 {
	if m != nil {
		return m.Weight
	}
	return 0
}

func (m *Consignment) GetContainers() []*Container {
	if m != nil {
		return m.Containers
	}
	return nil
}

func (m *Consignment) GetVesselId() string {
	if m != nil {
		return m.VesselId
	}
	return ""
}

// 单个集装箱
type Container struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	CustomerId           string   `protobuf:"bytes,2,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Origin               string   `protobuf:"bytes,3,opt,name=origin,proto3" json:"origin,omitempty"`
	UserId               string   `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Container) Reset()         { *m = Container{} }
func (m *Container) String() string { return proto.CompactTextString(m) }
func (*Container) ProtoMessage()    {}
func (*Container) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c6f11f9923110da, []int{1}
}

func (m *Container) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Container.Unmarshal(m, b)
}
func (m *Container) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Container.Marshal(b, m, deterministic)
}
func (m *Container) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Container.Merge(m, src)
}
func (m *Container) XXX_Size() int {
	return xxx_messageInfo_Container.Size(m)
}
func (m *Container) XXX_DiscardUnknown() {
	xxx_messageInfo_Container.DiscardUnknown(m)
}

var xxx_messageInfo_Container proto.InternalMessageInfo

func (m *Container) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Container) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *Container) GetOrigin() string {
	if m != nil {
		return m.Origin
	}
	return ""
}

func (m *Container) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

// 托运结果
type Response struct {
	Created              bool           `protobuf:"varint,1,opt,name=created,proto3" json:"created,omitempty"`
	Consignment          *Consignment   `protobuf:"bytes,2,opt,name=consignment,proto3" json:"consignment,omitempty"`
	Consignments         []*Consignment `protobuf:"bytes,3,rep,name=consignments,proto3" json:"consignments,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c6f11f9923110da, []int{2}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCreated() bool {
	if m != nil {
		return m.Created
	}
	return false
}

func (m *Response) GetConsignment() *Consignment {
	if m != nil {
		return m.Consignment
	}
	return nil
}

func (m *Response) GetConsignments() []*Consignment {
	if m != nil {
		return m.Consignments
	}
	return nil
}

// 查看货物信息的请求
// 客户端想要从服务端请求数据，必须有请求格式，哪怕为空
type GetRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c6f11f9923110da, []int{3}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Consignment)(nil), "consignment.Consignment")
	proto.RegisterType((*Container)(nil), "consignment.Container")
	proto.RegisterType((*Response)(nil), "consignment.Response")
	proto.RegisterType((*GetRequest)(nil), "consignment.GetRequest")
}

func init() { proto.RegisterFile("consignment/consignment.proto", fileDescriptor_2c6f11f9923110da) }

var fileDescriptor_2c6f11f9923110da = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x52, 0xdd, 0x4e, 0xc2, 0x30,
	0x14, 0x76, 0xfc, 0xef, 0x8c, 0x48, 0x6c, 0x22, 0x34, 0x1a, 0xe3, 0xb2, 0x2b, 0xae, 0x30, 0xc1,
	0xc4, 0x0b, 0xe3, 0xdd, 0x12, 0x09, 0xb7, 0xe5, 0x01, 0x0c, 0x6e, 0x27, 0xe3, 0x24, 0xd2, 0xce,
	0xb6, 0xe0, 0xdb, 0xf0, 0x04, 0x3e, 0xa4, 0xa1, 0x63, 0x52, 0x34, 0xdc, 0xf5, 0xfb, 0xeb, 0xf9,
	0x7a, 0x36, 0xb8, 0xcb, 0x94, 0x34, 0x54, 0xc8, 0x35, 0x4a, 0xfb, 0xe0, 0x9d, 0x27, 0xa5, 0x56,
	0x56, 0xb1, 0xc8, 0xa3, 0x92, 0xef, 0x00, 0xa2, 0xf4, 0x88, 0xd9, 0x25, 0x34, 0x28, 0xe7, 0x41,
	0x1c, 0x8c, 0x43, 0xd1, 0xa0, 0x9c, 0xc5, 0x10, 0xe5, 0x68, 0x32, 0x4d, 0xa5, 0x25, 0x25, 0x79,
	0xc3, 0x09, 0x3e, 0xc5, 0x86, 0xd0, 0xf9, 0x42, 0x2a, 0x56, 0x96, 0x37, 0xe3, 0x60, 0xdc, 0x16,
	0x07, 0xc4, 0x9e, 0x00, 0x32, 0x25, 0xed, 0x92, 0x24, 0x6a, 0xc3, 0x5b, 0x71, 0x73, 0x1c, 0x4d,
	0x87, 0x13, 0xbf, 0x4e, 0x5a, 0xcb, 0xc2, 0x73, 0xb2, 0x5b, 0x08, 0xb7, 0x68, 0x0c, 0x7e, 0xbc,
	0x51, 0xce, 0xdb, 0x6e, 0x5e, 0xaf, 0x22, 0xe6, 0x79, 0xb2, 0x86, 0xf0, 0x37, 0xf5, 0xaf, 0xeb,
	0x3d, 0x44, 0xd9, 0xc6, 0x58, 0xb5, 0x46, 0xbd, 0xcf, 0x56, 0x5d, 0xa1, 0xa6, 0xe6, 0xf9, 0xbe,
	0xaa, 0xd2, 0x54, 0x90, 0x74, 0x55, 0x43, 0x71, 0x40, 0x6c, 0x04, 0xdd, 0x8d, 0xa9, 0x42, 0xad,
	0x4a, 0xd8, 0xc3, 0x79, 0x9e, 0xec, 0x02, 0xe8, 0x09, 0x34, 0xa5, 0x92, 0x06, 0x19, 0x87, 0x6e,
	0xa6, 0x71, 0x69, 0xb1, 0x9a, 0xd9, 0x13, 0x35, 0x64, 0xcf, 0xe0, 0xef, 0xd4, 0x0d, 0x8e, 0xa6,
	0xfc, 0xef, 0x5b, 0xeb, 0xb3, 0xf0, 0xcd, 0xec, 0x05, 0xfa, 0x1e, 0x34, 0xbc, 0xe9, 0x16, 0x75,
	0x3e, 0x7c, 0xe2, 0x4e, 0xfa, 0x00, 0x33, 0xb4, 0x02, 0x3f, 0x37, 0x68, 0xec, 0x74, 0x17, 0xc0,
	0x60, 0xb1, 0xa2, 0xb2, 0x24, 0x59, 0x2c, 0x50, 0x6f, 0x29, 0x43, 0xf6, 0x0a, 0x57, 0xa9, 0xab,
	0xe9, 0x7f, 0xe5, 0xb3, 0xd7, 0xdf, 0x5c, 0x9f, 0x28, 0xf5, 0xdb, 0x93, 0x0b, 0x96, 0xc2, 0x60,
	0x86, 0xd6, 0xb3, 0x1a, 0x36, 0x3a, 0xf1, 0x1e, 0x7b, 0x9c, 0xbd, 0xe4, 0xbd, 0xe3, 0xfe, 0xc0,
	0xc7, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x12, 0x9b, 0x52, 0x59, 0xa2, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for ShippingService service

type ShippingServiceClient interface {
	// 托运一批货物
	CreateConsignment(ctx context.Context, in *Consignment, opts ...client.CallOption) (*Response, error)
	// 查看托运货物的信息
	GetConsignments(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*Response, error)
}

type shippingServiceClient struct {
	c           client.Client
	serviceName string
}

func NewShippingServiceClient(serviceName string, c client.Client) ShippingServiceClient {
	if c == nil {
		c = client.NewClient()
	}
	if len(serviceName) == 0 {
		serviceName = "consignment"
	}
	return &shippingServiceClient{
		c:           c,
		serviceName: serviceName,
	}
}

func (c *shippingServiceClient) CreateConsignment(ctx context.Context, in *Consignment, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "ShippingService.CreateConsignment", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *shippingServiceClient) GetConsignments(ctx context.Context, in *GetRequest, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.serviceName, "ShippingService.GetConsignments", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ShippingService service

type ShippingServiceHandler interface {
	// 托运一批货物
	CreateConsignment(context.Context, *Consignment, *Response) error
	// 查看托运货物的信息
	GetConsignments(context.Context, *GetRequest, *Response) error
}

func RegisterShippingServiceHandler(s server.Server, hdlr ShippingServiceHandler, opts ...server.HandlerOption) {
	s.Handle(s.NewHandler(&ShippingService{hdlr}, opts...))
}

type ShippingService struct {
	ShippingServiceHandler
}

func (h *ShippingService) CreateConsignment(ctx context.Context, in *Consignment, out *Response) error {
	return h.ShippingServiceHandler.CreateConsignment(ctx, in, out)
}

func (h *ShippingService) GetConsignments(ctx context.Context, in *GetRequest, out *Response) error {
	return h.ShippingServiceHandler.GetConsignments(ctx, in, out)
}
